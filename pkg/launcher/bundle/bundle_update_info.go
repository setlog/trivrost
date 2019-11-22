package bundle

import (
	"github.com/setlog/trivrost/pkg/launcher/config"
	log "github.com/sirupsen/logrus"
)

// BundleUpdateInfo contains information on what files need updating on the user's machine for the bundle specified by the embedded BundleConfig.
type BundleUpdateInfo struct {
	config.BundleConfig
	IsSystemBundle bool
	PresentState   config.FileInfoMap
	RemoteState    config.FileInfoMap
	WantedState    config.FileInfoMap
}

func (bui *BundleUpdateInfo) LogChanges() {
	for filePath, wantedFileInfo := range bui.WantedState {
		presentFileInfo, ok := bui.PresentState[filePath]
		if ok {
			if wantedFileInfo.SHA256 == "" {
				log.Infof("\"%s\": Delete: %s", filePath, presentFileInfo.SHA256)
			} else {
				log.Infof("\"%s\": %s -> %s", filePath, presentFileInfo.SHA256, wantedFileInfo.SHA256)
			}
		} else {
			log.Infof("\"%s\": Create: %s", filePath, wantedFileInfo.SHA256)
		}
	}
}
