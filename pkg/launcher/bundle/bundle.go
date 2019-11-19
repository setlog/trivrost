package bundle

import (
	"path/filepath"

	"github.com/setlog/trivrost/pkg/system"
)

func (u *Updater) haveSystemBundleWithName(localDirectory string) bool {
	return system.FolderExists(filepath.Join(u.systemBundlesFolderPath, localDirectory))
}

func (u *Updater) wantBundleWithName(localDirectory string) bool {
	for _, bundleUpdateInfo := range u.bundleUpdateInfos {
		if localDirectory == bundleUpdateInfo.LocalDirectory {
			return true
		}
	}
	return false
}

func (u *Updater) HasChangesToSystemBundles(considerMandatoryChangesOnly bool) bool {
	for _, bundleUpdateInfo := range u.bundleUpdateInfos {
		if (!considerMandatoryChangesOnly || bundleUpdateInfo.IsUpdateMandatory) && bundleUpdateInfo.IsSystemBundle && bundleUpdateInfo.WantedState.HasChanges() {
			return true
		}
	}
	return false
}

func (u *Updater) HasChangesToUserBundles() bool {
	for _, bundleUpdateInfo := range u.bundleUpdateInfos {
		if !bundleUpdateInfo.IsSystemBundle && bundleUpdateInfo.WantedState.HasChanges() {
			return true
		}
	}
	return false
}

func countUpdatesBytes(bundleUpdateConfigs []*BundleUpdateInfo) uint64 {
	var total uint64
	for _, bundleUpdateConfig := range bundleUpdateConfigs {
		total += bundleUpdateConfig.WantedState.UpdateByteCount()
	}
	return total
}
