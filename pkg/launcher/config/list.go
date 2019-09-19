package config

import (
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func NewFileInfoMap() FileInfoMap {
	return make(FileInfoMap)
}

func (fm FileInfoMap) Join(other FileInfoMap) {
	for k, v := range other {
		fm[k] = v
	}
}

func (fm FileInfoMap) MustGetOnly() (key string, fileInfo *FileInfo) {
	if len(fm) != 1 {
		panic(spew.Sprintf("Not exactly one entry in FileInfoMap. Got %d: %+v.", len(fm), fm))
	}
	for k, v := range fm {
		return k, v
	}
	panic("Invalid state: Failed to return the only key-value pair in the FileInfoMap")
}

// Returns a FileInfoMap which describes the changes from have to want.
// In the returned FileInfoMap, a key which maps to a *FileInfo with its SHA256-field being non-empty indicates
// that the file described by the key-string requires an update or is new. A key which maps to a *FileInfo
// with an empty SHA256-field indicates that the file described by the key-string should be removed.
func MakeDiffFileInfoMap(have FileInfoMap, want FileInfoMap) FileInfoMap {
	fm := make(FileInfoMap)
	for presentKey, presentFileInfo := range have {
		if _, ok := want[presentKey]; !ok {
			log.Debugf("Have %s but want no hash for \"%s\". Delete.", presentFileInfo.SHA256[:8], presentKey)
			newFileInfo := *presentFileInfo
			newFileInfo.SHA256 = ""
			fm[presentKey] = &newFileInfo
		}
	}
	for wantedKey, wantedFileInfo := range want {
		if presentFileInfo, ok := have[wantedKey]; ok {
			if wantedFileInfo.SHA256 != presentFileInfo.SHA256 {
				log.Debugf("Have %s but want %s for \"%s\". Update.", presentFileInfo.SHA256[:8], wantedFileInfo.SHA256[:8], wantedKey)
				newFileInfo := *wantedFileInfo
				fm[wantedKey] = &newFileInfo
			}
		} else {
			log.Debugf("Have no hash but want %s for \"%s\". Update.", wantedFileInfo.SHA256[:8], wantedKey)
			newFileInfo := *wantedFileInfo
			fm[wantedKey] = &newFileInfo
		}
	}
	return fm
}

func (fm FileInfoMap) HasChanges() bool {
	return len(fm) > 0
}

func (fm FileInfoMap) HasUpdates() bool {
	return fm.UpdateFileCount() > 0
}

func (fm FileInfoMap) DeleteFileCount() uint64 {
	var total uint64
	for _, v := range fm {
		if v.SHA256 == "" {
			total++
		}
	}
	return total
}

func (fm FileInfoMap) UpdateFileCount() uint64 {
	var total uint64
	for _, v := range fm {
		if v.SHA256 != "" {
			total++
		}
	}
	return total
}

func (fm FileInfoMap) UpdateByteCount() uint64 {
	var total uint64
	for _, v := range fm {
		if v.SHA256 != "" {
			total += uint64(v.Size)
		}
	}
	return total
}
