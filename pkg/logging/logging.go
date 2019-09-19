package logging

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/setlog/trivrost/pkg/system"

	golog "log"

	log "github.com/sirupsen/logrus"
)

var folderPath string
var topic string
var filePath string
var wasInitializeCalled bool
var initTime time.Time

const maxIndex = 9999

var digitCount = len(strconv.Itoa(maxIndex))

func Initialize(logFolderPath string, logTopic string, logIndex int, logInstance int) int {
	if wasInitializeCalled {
		panic("Initialize() was already called.")
	}
	wasInitializeCalled = true

	folderPath = logFolderPath
	topic = logTopic

	log.SetReportCaller(true)
	log.SetFormatter(&LogFormatter{})

	logFileName, nextLogIndex := getLogFileName(logIndex, logInstance)
	filePath = filepath.Join(folderPath, logFileName)

	vTermError := enableVirtualTerminalProcessing()
	if vTermError != nil {
		log.Warnf("Enable virtual terminal processing: %v", vTermError)
	}

	configureLogrusOutput(openLogFile(), vTermError == nil)
	configureDefaultLogOutput()

	initTime = time.Now()
	log.WithFields(log.Fields{"Binary": system.GetBinaryPath(), "Program": system.GetProgramPath(), "Args": os.Args,
		"OS": runtime.GOOS, "Arch": system.GetOSArch(), "Log": filePath, "Local Time": initTime, "UTC Time": initTime.UTC()}).Info("")

	return nextLogIndex
}

func configureLogrusOutput(logFile io.Writer, isVirtualTerminalProcessingEnabled bool) {
	if logFile == nil {
		if isVirtualTerminalProcessingEnabled {
			log.SetOutput(os.Stderr)
		} else {
			log.SetOutput(&withoutStyleWriter{target: os.Stderr})
		}
	} else {
		if isVirtualTerminalProcessingEnabled {
			log.SetOutput(io.MultiWriter(&withoutStyleWriter{target: logFile}, os.Stderr))
		} else {
			log.SetOutput(&withoutStyleWriter{target: io.MultiWriter(logFile, os.Stderr)})
		}
	}
}

func configureDefaultLogOutput() {
	golog.SetFlags(golog.Llongfile) // logRelay parses this information.
	golog.SetPrefix("")
	golog.SetOutput(&logRelay{})
}

func GetLogFilePath() string {
	return filePath
}

func DeleteOldLogFiles() {
	maxLogFileAge := time.Hour * 24 * 20
	infos, err := ioutil.ReadDir(folderPath)
	now := time.Now()
	if err != nil {
		log.Errorf("Could not read contents of log directory \"%s\": %v", folderPath, err)
	} else {
		for _, info := range infos {
			if isLogFileName(info.Name()) && info.ModTime().Add(maxLogFileAge).Before(now) {
				deleteLogFilePath := filepath.Join(folderPath, info.Name())
				err = os.Remove(deleteLogFilePath)
				if err != nil {
					log.Errorf("Could not remove old log file \"%s\": %v", deleteLogFilePath, err)
				} else {
					log.Infof("Removed old log file \"%s\".", deleteLogFilePath)
				}
			}
		}
	}
}

func openLogFile() *os.File {
	err := os.MkdirAll(folderPath, 0700)
	if err != nil {
		log.WithFields(log.Fields{"folderPath": folderPath, "err": err}).
			Warn("Could not create nested directory for log files.")
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0600)
	if err != nil {
		log.Errorf("Failed to log on file system: %v", err)
	} else {
		log.RegisterExitHandler(func() {
			file.Close()
		})
	}
	return file
}

func getLogFileName(logIndex int, logInstance int) (string, int) {
	now := time.Now().UTC()
	descriptor := topic + "." + now.Format("2006-01-02_15-04-05")
	if logIndex != -1 {
		return concatenateLogFileNameArtifacts(logIndex, logInstance, descriptor), logIndex
	}
	infos, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "folderPath": folderPath}).Warning("Could not read directory contents.")
		return concatenateLogFileNameArtifacts(0, logInstance, descriptor), 0
	}
	useIndex := getNextIndex(infos)
	return concatenateLogFileNameArtifacts(useIndex, logInstance, descriptor), useIndex
}

func getNextIndex(infos []os.FileInfo) int {
	latestIndex := getLatestIndex(infos)
	if latestIndex < maxIndex {
		return (latestIndex + 1) % (maxIndex + 1)
	}
	firstAvailableIndex := getFirstAvailableIndex(infos)
	if firstAvailableIndex > maxIndex {
		log.Warnf("No free log index left.")
		return maxIndex
	}
	return firstAvailableIndex
}

func getLatestIndex(infos []os.FileInfo) int {
	for i := len(infos) - 1; i >= 0; i-- {
		name := infos[i].Name()
		if isLogFileName(name) {
			infoIndex, err := strconv.Atoi(name[:digitCount])
			if err == nil {
				return infoIndex
			}
		}
	}
	return -1
}

func getFirstAvailableIndex(infos []os.FileInfo) int {
	firstAvailableIndex := 0
	for _, info := range infos {
		name := info.Name()
		if isLogFileName(name) {
			infoIndex, err := strconv.Atoi(name[:digitCount])
			if err == nil {
				if infoIndex > firstAvailableIndex {
					return firstAvailableIndex
				}
				firstAvailableIndex = infoIndex + 1
			}
		}
	}
	return firstAvailableIndex
}

func isLogFileName(fileName string) bool {
	var index int
	var letter rune
	var rest string
	_, err := fmt.Sscanf(fileName, "%0"+strconv.Itoa(digitCount)+"d%c.%s", &index, &letter, &rest)
	if err != nil || strings.Index(fileName, topic) != 6 || !strings.HasSuffix(fileName, ".log") {
		return false
	}
	return true
}

func concatenateLogFileNameArtifacts(logIndex int, logInstance int, descriptor string) string {
	logInstanceRune := 'a' + logInstance
	if logInstanceRune > 'z' {
		logInstanceRune = 'z'
	}
	return fmt.Sprintf("%0"+strconv.Itoa(digitCount)+"d%c.%s.log", logIndex, logInstanceRune, descriptor)
}
