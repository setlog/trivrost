package fetching

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/setlog/trivrost/pkg/misc"
	"github.com/setlog/trivrost/pkg/signatures"
	"git.sr.ht/~tslocum/preallocate"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"
)

const defaultTimeout = time.Second * 30

const MaxConcurrentDownloads = 5

// Downloader has helper functions for common use cases of Download, such as writing a resource to a file while downloading it,
// downloading multiple resources in parallel and verifying the hashsum or signature of downloading resources.
type Downloader struct {
	handler DownloadProgressHandler
	client  *http.Client
	ctx     context.Context
}

func NewDownloader(ctx context.Context, handler DownloadProgressHandler) *Downloader {
	return &Downloader{handler: handler, client: MakeClient(), ctx: ctx}
}

func (downloader *Downloader) DownloadSignedResource(fromURL string, keys []*rsa.PublicKey) ([]byte, error) {
	data, err := downloader.DownloadSignedResources([]string{fromURL}, keys)
	if err != nil {
		return nil, err
	}
	return data[fromURL], nil
}

func (downloader *Downloader) DownloadSignedResources(urls []string, keys []*rsa.PublicKey) (map[string][]byte, error) {
	fileMapWithSignatures := make(config.FileInfoMap)
	for _, url := range urls {
		fileMapWithSignatures[url] = &config.FileInfo{}
		if strings.HasPrefix(url, "file://") {
			log.Printf("Skipping signature validation for resource \"%s\" because it uses the \"file://\"-scheme.", url)
		} else {
			fileMapWithSignatures[url+".signature"] = &config.FileInfo{}
		}
	}
	fileData, err := downloader.DownloadToRAM(fileMapWithSignatures)
	if err != nil {
		return nil, err
	}
	for _, url := range urls {
		if !strings.HasPrefix(url, "file://") && !signatures.IsSignatureValid(fileData[url], fileData[url+".signature"], keys) {
			log.Panicf("Invalid signature of resource %s.", url)
		}
	}
	validatedResources := make(map[string][]byte)
	for _, url := range urls {
		validatedResources[url] = fileData[url]
	}
	return validatedResources, nil
}

func (downloader *Downloader) DownloadBytes(fromURL string) (data []byte) {
	success := false
	var err error
	for !success {
		dl := downloader.newDownload(fromURL)
		data, err = ioutil.ReadAll(dl)
		if err != nil {
			log.Printf("Download of \"%s\" failed: %v", fromURL, err)
		}
		if downloader.ctx.Err() != nil {
			panic(downloader.ctx.Err())
		}
		success = err == nil
	}
	return
}

func (downloader *Downloader) newDownload(resourceUrl string) *Download {
	return &Download{
		url:      resourceUrl,
		client:   downloader.client,
		ctx:      downloader.ctx,
		handler:  downloader.handler,
		workerId: 0,
	}
}

func (downloader *Downloader) MustDownloadToTempDirectory(baseUrl string, fileMap config.FileInfoMap, localDirPath string) (tempDirectoryPath string) {
	reachedEndOfFunction := false // https://stackoverflow.com/a/34851179/10513183
	tempDirectoryPath = system.MustMakeTempDirectory(localDirPath)
	defer func() {
		if !reachedEndOfFunction {
			system.TryRemoveDirectory(tempDirectoryPath)
		}
	}()
	downloader.MustDownloadToDirectory(baseUrl, fileMap, tempDirectoryPath)
	reachedEndOfFunction = true
	return tempDirectoryPath
}

func (downloader *Downloader) MustDownloadToDirectory(baseUrl string, fileMap config.FileInfoMap, localDirPath string) {
	err := os.MkdirAll(localDirPath, 0700)
	if err != nil {
		panic(system.NewFileSystemError(fmt.Sprintf("Could not create directory \"%s\"", localDirPath), err))
	}
	err = downloader.DownloadToDirectory(baseUrl, fileMap, localDirPath)
	if err != nil {
		panic(err)
	}
}

func (downloader *Downloader) DownloadToDirectory(baseUrl string, fileMap config.FileInfoMap, localDirPath string) error {
	fileMap = fileMap.OmitEntriesWithMissingSha()
	urlToPathMap := make(map[string]string)
	urlToInfoMap := make(config.FileInfoMap)
	for relativeFilePath, fileInfo := range fileMap {
		url := misc.MustJoinURL(baseUrl, filepath.ToSlash(relativeFilePath))
		urlToPathMap[url] = relativeFilePath
		urlToInfoMap[url] = fileInfo
	}
	return downloader.DownloadResources(stringStringMapKeys(urlToPathMap), func(dl *Download) error {
		return updateFile(dl, urlToInfoMap[dl.url], filepath.Join(localDirPath, urlToPathMap[dl.url]))
	})
}

func (downloader *Downloader) DownloadFile(url string, fileInfo *config.FileInfo, filePath string) error {
	return downloader.DownloadResource(url, func(dl *Download) error {
		return updateFile(dl, fileInfo, filePath)
	})
}

func (downloader *Downloader) DownloadToRAM(fileMap config.FileInfoMap) (fileData map[string][]byte, dlErr error) {
	m := &sync.Mutex{}
	fileData = make(map[string][]byte)
	return fileData, downloader.DownloadResources(fileMap.FilePaths(), func(dl *Download) error {
		wantedFileInfo := fileMap[dl.url]
		hash := sha256.New()
		data, err := ioutil.ReadAll(io.TeeReader(dl, hash))
		if err != nil {
			if downloader.ctx.Err() != nil {
				return downloader.ctx.Err()
			}
			return fmt.Errorf("ioutil.ReadAll failed: %v", err)
		}
		dlFileSha := hex.EncodeToString(hash.Sum(nil))
		if wantedFileInfo.SHA256 != "" && !strings.EqualFold(wantedFileInfo.SHA256, dlFileSha) {
			return fmt.Errorf("Sha of downloaded file \"%s\" does not match expected value \"%s\". Was \"%s\"",
				dl.url, wantedFileInfo.SHA256, dlFileSha)
		}
		m.Lock()
		defer m.Unlock()
		fileData[dl.url] = data
		return nil
	})
}

func (downloader *Downloader) DownloadResource(url string, processDownload func(dl *Download) error) error {
	ctx, cancelFunc := context.WithCancel(downloader.ctx)
	defer cancelFunc()
	return downloader.runDownloadWorkers(ctx, cancelFunc, []string{url}, processDownload)
}

// DownloadResources makes a goroutined call of processDownload with a *Download ready for Read()s for every url
// in urls, never allowing more than MaxConcurrentDownloads to be processed at the same time and returning a non-nil
// error if and only if any of the calls to processDownload do return a non-nil error or the context is cancelled.
func (downloader *Downloader) DownloadResources(urls []string, processDownload func(dl *Download) error) error {
	ctx, cancelFunc := context.WithCancel(downloader.ctx)
	defer cancelFunc()
	return downloader.runDownloadWorkers(ctx, cancelFunc, urls, processDownload)
}

func (downloader *Downloader) runDownloadWorkers(ctx context.Context, cancelFunc context.CancelFunc, urls []string, processDownload func(dl *Download) error) error {
	availableWorkerIds := createWorkerIdChannel(MaxConcurrentDownloads)
	errChan := make(chan error, 1)
	workerErrChan := make(chan error, len(urls))
	allWorkersDoneCond := &sync.Cond{L: &sync.Mutex{}}
	go reportWorkersResult(ctx, cancelFunc, workerErrChan, len(urls), errChan)
	for _, url := range urls {
		if ctx.Err() != nil {
			break
		}
		select {
		case <-ctx.Done():
			break
		case workerId := <-availableWorkerIds:
			dl := NewDownloadForConcurrentUse(ctx, url, downloader.client, downloader.handler, workerId)
			go downloadWorker(dl, availableWorkerIds, allWorkersDoneCond, workerErrChan, processDownload)
		}
	}
	allWorkersDoneCond.L.Lock()
	for len(availableWorkerIds) < MaxConcurrentDownloads {
		allWorkersDoneCond.Wait()
	}
	allWorkersDoneCond.L.Unlock()
	return <-errChan
}

func reportWorkersResult(ctx context.Context, cancelFunc context.CancelFunc, workerErrChan chan error, expectedResultCount int, errChan chan error) {
	resultCount := 0
	for resultCount < expectedResultCount {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		case err := <-workerErrChan:
			if err != nil {
				cancelFunc()
				errChan <- err
				return
			}
			resultCount++
		}
	}
	errChan <- nil
}

func downloadWorker(dl *Download, workerIds chan int, allWorkersDoneCond *sync.Cond, workerErrChan chan error, processDownload func(dl *Download) error) {
	var workerErr error
	defer func() {
		panicObject := recover()
		if panicObject != nil {
			workerErr = fmt.Errorf("worker %d panicked: %v", dl.workerId, panicObject)
		}
		allWorkersDoneCond.L.Lock()
		workerErrChan <- workerErr
		workerIds <- dl.workerId
		allWorkersDoneCond.L.Unlock()
		allWorkersDoneCond.Signal()
		if panicObject != nil {
			defer misc.LogPanic()
			panic(panicObject)
		}
	}()
	workerErr = processDownload(dl)
}

func updateFile(dl *Download, expectedFileInfo *config.FileInfo, localFilePath string) (returnError error) {
	system.MustMakeDir(filepath.Dir(localFilePath))
	file, err := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return system.NewFileSystemError(fmt.Sprintf("Could not open file \"%s\" for writing", localFilePath), err)
	}
	defer system.CleanUpFileOperation(file, &returnError)
	if err = preallocate.File(file, expectedFileInfo.Size); err != nil { // Important: Screws up royally on files opened with the os.O_APPEND-flag.
		log.Printf("Could not preallocate file \"%s\" with %d bytes: %v", localFilePath, expectedFileInfo.Size, err)
	}
	n, dlFileSha, err := ioHashingCopy(dl.ctx, file, dl)
	if err != nil {
		return err
	}
	if !strings.EqualFold(expectedFileInfo.SHA256, dlFileSha) {
		return fmt.Errorf("Sha of downloaded file \"%s\" does not match expected value \"%s\" for file \"%s\". Was \"%s\"",
			dl.url, expectedFileInfo.SHA256, localFilePath, dlFileSha)
	}
	if n < expectedFileInfo.Size { // Needed to prevent trailing null bytes.
		return fmt.Errorf("Wrote less bytes than expected after preallocating file \"%s\". Written: %d; Expected: %d", localFilePath, n, expectedFileInfo.Size)
	}
	return nil
}

func ioHashingCopy(ctx context.Context, dst io.Writer, src io.Reader) (int64, string, error) {
	hash := sha256.New()
	n, err := io.Copy(dst, io.TeeReader(src, hash))
	if err != nil {
		if ctx.Err() != nil {
			return n, hex.EncodeToString(hash.Sum(nil)), ctx.Err()
		}
		return n, hex.EncodeToString(hash.Sum(nil)), fmt.Errorf("io.Copy failed: %v", err)
	}
	return n, hex.EncodeToString(hash.Sum(nil)), nil
}
