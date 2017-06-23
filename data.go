// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import "time"

type UserSettings struct {
	Username          string
	Password          string
	FilePath          string
	PasswordConfirm   string
	SetAsNonSteamGame bool
}

type QCUserCredentials struct {
	Username string
	Password string
	Token    string
}

type QCOptions struct {
	QCFilePath string
}

type UpdateTime struct {
	LastQCUpdateTime       int64
	LastLauncherUpdateTime int64
}

type LauncherUpdateInfo struct {
	LatestVersion float32
	Date          time.Time
	URL           string
}
