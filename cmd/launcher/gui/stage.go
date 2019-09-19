package gui

import log "github.com/sirupsen/logrus"

// Stage is an abstraction of the various steps which the launcher goes through when it runs.
type Stage int

const (
	StageAcquireLock Stage = iota
	StageGetDeploymentConfig
	StageDetermineLocalLauncherVersion
	StageRetrieveRemoteLauncherVersion
	StageSelfUpdate
	StageDetermineLocalBundleVersions
	StageRetrieveRemoteBundleVersions
	StageAwaitApplicationsTerminated
	StageDownloadBundleUpdates
	StageLaunchApplication
)

var (
	textAcquireLock                   = "Waiting for other launcher instance to finish..."
	textGetDeploymentConfig           = "Retrieving application configuration..."
	textDetermineLocalLauncherVersion = "Determining launcher version..."
	textRetrieveRemoteLauncherVersion = "Checking for launcher updates..."
	textSelfUpdate                    = "Updating launcher..."
	textDetermineLocalBundleVersions  = "Determining application version..."
	textRetrieveRemoteBundleVersions  = "Checking for application updates..."
	textAwaitApplicationsTerminated   = "Please close all instances of the application to apply the update."
	textDownloadBundleUpdates         = "Retrieving application update..."
	textLaunchApplication             = "Launching application..."
)

func SetStageText(s Stage, text string) {
	switch s {
	case StageAcquireLock:
		textAcquireLock = text
	case StageGetDeploymentConfig:
		textGetDeploymentConfig = text
	case StageDetermineLocalLauncherVersion:
		textDetermineLocalLauncherVersion = text
	case StageRetrieveRemoteLauncherVersion:
		textRetrieveRemoteLauncherVersion = text
	case StageSelfUpdate:
		textSelfUpdate = text
	case StageDetermineLocalBundleVersions:
		textDetermineLocalBundleVersions = text
	case StageRetrieveRemoteBundleVersions:
		textRetrieveRemoteBundleVersions = text
	case StageAwaitApplicationsTerminated:
		textAwaitApplicationsTerminated = text
	case StageDownloadBundleUpdates:
		textDownloadBundleUpdates = text
	case StageLaunchApplication:
		textLaunchApplication = text
	}
}

func (s Stage) getProgressInterval() (lowerEnd, upperEnd int) {
	switch s {
	case StageAcquireLock:
		return 0, 0
	case StageGetDeploymentConfig:
		return 1, 1
	case StageDetermineLocalLauncherVersion:
		return 2, 2
	case StageRetrieveRemoteLauncherVersion:
		return 3, 3
	case StageSelfUpdate:
		return 4, 10
	case StageDetermineLocalBundleVersions:
		return 11, 17
	case StageRetrieveRemoteBundleVersions:
		return 18, 18
	case StageAwaitApplicationsTerminated:
		return 19, 19
	case StageDownloadBundleUpdates:
		return 20, 99
	case StageLaunchApplication:
		return 100, 100
	}
	log.Warnf("No progress interval for stage %v.\n", s)
	return 0, 100
}

func (s Stage) getText() string {
	switch s {
	case StageAcquireLock:
		return textAcquireLock
	case StageGetDeploymentConfig:
		return textGetDeploymentConfig
	case StageDetermineLocalLauncherVersion:
		return textDetermineLocalLauncherVersion
	case StageRetrieveRemoteLauncherVersion:
		return textRetrieveRemoteLauncherVersion
	case StageSelfUpdate:
		return textSelfUpdate
	case StageDetermineLocalBundleVersions:
		return textDetermineLocalBundleVersions
	case StageRetrieveRemoteBundleVersions:
		return textRetrieveRemoteBundleVersions
	case StageAwaitApplicationsTerminated:
		return textAwaitApplicationsTerminated
	case StageDownloadBundleUpdates:
		return textDownloadBundleUpdates
	case StageLaunchApplication:
		return textLaunchApplication
	}
	log.Warnf("No status text for stage %v.\n", s)
	return "Working..."
}

func (s Stage) IsDownloadStage() bool {
	return s == StageGetDeploymentConfig ||
		s == StageRetrieveRemoteLauncherVersion ||
		s == StageSelfUpdate ||
		s == StageRetrieveRemoteBundleVersions ||
		s == StageDownloadBundleUpdates
}
