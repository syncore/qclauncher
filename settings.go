// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"errors"
	"fmt"
	"os"
)

type Configuration struct {
	Core         *QCCoreSettings
	Experimental *QCExperimentalSettings
	Launcher     *LauncherSettings
	Auth         *TokenAuth
}

var isCollectingSettings = false

func GetConfiguration() (*Configuration, error) {
	coreQCSettings := &QCCoreSettings{}
	err := Get(coreQCSettings)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error retrieving core QC configuration settings", GetCaller()), "error", err)
		return nil, err
	}
	expQCSettings := &QCExperimentalSettings{}
	err = Get(expQCSettings)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error retrieving experimental QC configuration settings", GetCaller()), "error", err)
		return nil, err
	}
	launcherSettings := &LauncherSettings{}
	err = Get(launcherSettings)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error retrieving launcher configuration settings", GetCaller()), "error", err)
		return nil, err
	}
	authToken := &TokenAuth{}
	err = Get(authToken)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error retrieving auth token", GetCaller()), "error", err)
		return nil, err
	}
	return &Configuration{
		Core:         coreQCSettings,
		Experimental: expQCSettings,
		Launcher:     launcherSettings,
		Auth:         authToken,
	}, nil
}

func GetEmptyConfiguration() *Configuration {
	return &Configuration{
		Core:         &QCCoreSettings{},
		Experimental: &QCExperimentalSettings{},
		Launcher:     &LauncherSettings{},
	}
}

func DeleteConfiguration(removeLock bool) {
	err := DeleteFile(GetDataFilePath())
	if err != nil && !os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("%s: error deleting %s: %s", GetCaller(), DataFile, err))
		Lock.Unlock()
		ShowFatalErrorMsg("Error",
			fmt.Sprintf("Unable to delete existing %s file during reset. Please manually delete it and then re-run QCLauncher.",
				DataFile), nil)
		return
	}
	if removeLock {
		Lock.Unlock()
	}
}

func configureSettings() {
	if isCollectingSettings {
		return
	}
	var err error
	var cfg *Configuration
	if FileExists(GetDataFilePath()) {
		cfg, err = GetConfiguration()
		if err != nil {
			ShowErrorMsg("Error", "An error occurred when retrieving your settings. Resetting.", nil)
			DeleteConfiguration(false)
			cfg = GetEmptyConfiguration()
		}
	} else {
		cfg = GetEmptyConfiguration()
	}
	isCollectingSettings = true
	settingsWindow := newSettingsWindow(cfg)
	settingsWindow.Run()
}

func saveConfiguration(cfg *Configuration) error {
	if err := validateOnSave(cfg); err != nil {
		return err
	}
	checkLog := fmt.Sprintf("please check the %s file for more details.", LogFile)
	if err := Save(cfg.Core); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error saving QC core settings", GetCaller()), "error", err)
		return fmt.Errorf("Unable to save QC core settings, %s", checkLog)
	}
	cfg.Experimental.resetBoolForZeroValues() // handle 0 values (no effect) and others
	if err := Save(cfg.Experimental); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error saving QC experimental settings", GetCaller()), "error", err)
		return fmt.Errorf("Unable to save QC experimental settings, %s", checkLog)
	}
	// Steam launch should be a one-time event
	launchSteam := cfg.Launcher.SetAsNonSteamGame
	cfg.Launcher.SetAsNonSteamGame = false
	if err := Save(cfg.Launcher); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error saving launcher settings", GetCaller()), "error", err)
		return fmt.Errorf("Unable to save launcher settings, %s", checkLog)
	}
	if err := updateLastCheckTime(UpdateAll, 0); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error setting last update check time", GetCaller()), "error", err)
		return err
	}
	if launchSteam {
		ShowInfoMsg("Launching Steam",
			"Now opening Steam to add Quake Champions. When Steam loads, browse to and select your qclauncher.exe file.",
			qclauncherSettingsWindow)
		if err := addNonSteamGame(getSteamInstallPath()); err != nil {
			ShowErrorMsg("Error", "Unable to launch Steam", qclauncherSettingsWindow)
		}
	}
	return nil
}

func validateOnSave(cfg *Configuration) error {
	if cfg == nil {
		return errors.New("No settings were entered")
	}
	if err := cfg.Core.validate(); err != nil {
		return err
	}
	if err := cfg.Experimental.validate(); err != nil {
		return err
	}
	if err := cfg.Launcher.validate(); err != nil {
		return err
	}
	return nil
}
