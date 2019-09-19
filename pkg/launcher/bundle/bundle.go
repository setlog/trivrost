package bundle

import (
	"path/filepath"

	"github.com/setlog/trivrost/pkg/misc"
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

func (u *Updater) assertUpdatePossible() {
	if u.hasChangesToSystemBundles() && u.hasChangesToUserBundles() {
		panic(misc.NewNestedError("There is an update which needs to make changes to protected system folders.\n"+
			"Your system administrator should already be aware of this, and apply the update shortly.", nil))
	}
}

func (u *Updater) hasChangesToSystemBundles() bool {
	for _, bundleUpdateInfo := range u.bundleUpdateInfos {
		if bundleUpdateInfo.IsSystemBundle && bundleUpdateInfo.WantedState.HasChanges() {
			return true
		}
	}
	return false
}

func (u *Updater) hasChangesToUserBundles() bool {
	for _, bundleUpdateInfo := range u.bundleUpdateInfos {
		if !bundleUpdateInfo.IsSystemBundle && bundleUpdateInfo.WantedState.HasChanges() {
			return true
		}
	}
	return false
}

func (u *Updater) isAtLeastOneChangeRequired() bool {
	return u.hasChangesToSystemBundles() || u.hasChangesToUserBundles()
}

func countUpdatesBytes(bundleUpdateConfigs []*BundleUpdateInfo) uint64 {
	var total uint64
	for _, bundleUpdateConfig := range bundleUpdateConfigs {
		total += bundleUpdateConfig.WantedState.UpdateByteCount()
	}
	return total
}
