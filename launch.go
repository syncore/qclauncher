// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

const (
	QCExe                     = "QuakeChampions.exe"
	defArgs                   = `--startup --set /Config/GAME_CONFIG/bethesdaGameCode "%GAMECODE%" --set /Config/GAME_CONFIG/bethesdaLoginEnabled 1 --set /Config/Bethesda/Language "%LANGUAGE%" --set /Config/GAME_CONFIG/bethesdaEndpointUrl "https://services.bethesda.net/agora_beam/"`
	qcEntitlmentID            = 48329
	qcProjectID               = 11
	qcDefaultBranchIdentifier = "Default"
	gameCodeTempl             = "%GAMECODE%"
	langTempl                 = "%LANGUAGE%"
	defLang                   = "en"
)

func Launch() error {
	running, _, _, _, namepids, err := IsProcessRunning(QCExe)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error enumerating running processes to see if QC is already running", GetCaller()),
			"error", err)
		ShowErrorMsg("Error", "Unable to enumerate processes to determine if Quake Champions is already running", nil)
		return err
	}
	if running {
		if willClose := ShowQCRunningMsg(namepids[QCExe]); !willClose {
			return &alreadyRunningError{emsg: "Quake Champions is already running, cannot start."}
		}
	}
	cfg, err := GetConfiguration()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: fatal error: unable to load data file during launch", GetCaller()), "error", err)
		ShowFatalErrorMsg("Error", fmt.Sprintf("Could not read your %s file. Cannot start.", DataFile), nil)
		return err
	}
	lc := newLauncherClient(defTimeout)
	lc.checkServerStatus()
	if err = CheckUpdate(ConfEnforceHash, UpdateQC); IsErrHashMismatch(err) {
		return &hashMismatchError{emsg: err.Error()}
	}
	err = lc.authenticate(cfg)
	if err != nil {
		return err
	}
	entitlementInfo, err := lc.getEntitlementInfo()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: getEntitlementInfo error", GetCaller()), "error", err, "data", entitlementInfo)
		return err
	}
	logger.Debugw("Entitlement info", "entitlementInfo", entitlementInfo)
	projectID, branchID, _, err := getProjectBranchBuildIdentifiers(entitlementInfo)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: getBuildBranchIdentifiers error", GetCaller()), "error", err)
		return err
	}
	branchInfo, err := lc.getBranchInfo(projectID, branchID)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: getBranchInfo error", GetCaller()), "error", err, "data", branchInfo)
		return err
	}
	logger.Debugw("Branch info", "branchInfo", branchInfo)
	launchArgs, err := lc.getLaunchArgs(projectID)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: getLaunchArgs error", GetCaller()), "error", err, "data", launchArgs)
		return err
	}
	exArgs := launchArgs.extractLaunchArgs(branchInfo.LaunchinfoList[0], cfg.Core.Language)
	logger.Debugw("Launch args", "launchArgs", launchArgs)
	logger.Debugw("Extracted launch args", "exArgs", exArgs)
	gameCode, err := lc.getGameCode(projectID)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: getGameCode error", GetCaller()), "error", err, "data", gameCode)
		return err
	}
	logger.Debugw("Game code", "gameCode.GameCode", gameCode.Gamecode)
	baseArgs := strings.Replace(exArgs, gameCodeTempl, gameCode.Gamecode, -1)
	return runQC(cfg, baseArgs)
}

func Exit(code int) {
	Lock.Unlock()
	os.Exit(code)
}

func (r *LaunchArgsResponse) extractLaunchArgs(liKey int, language string) string {
	fallback := strings.Replace(defArgs, langTempl, defLang, -1)
	if r == nil {
		logger.Info("extractLaunchArgs: launch args was nil, using fallback arguments")
		return fallback
	}
	if (LaunchInfo{}) == r.LaunchinfoSet {
		logger.Info("extractLaunchArgs: launch info set had default struct values, using fallback arguments")
		return fallback
	}
	if (LaunchInfoItem{}) == r.LaunchinfoSet.Default {
		logger.Info("extractLaunchArgs: default launch info item had default struct values, using fallback arguments")
		return fallback
	}
	likstr := strconv.Itoa(liKey)
	val := reflect.ValueOf(r.LaunchinfoSet)
	found := false
	for i := 0; i < val.Type().NumField(); i++ {
		v := val.Type().Field(i).Tag.Get("json")
		if v != likstr {
			continue
		}
		found = true
	}
	if found {
		v := r.LaunchinfoSet.Default.LaunchArgs
		v = strings.Replace(v, "\\", "", -1)
		v = strings.Replace(v, langTempl, language, -1)
		return v
	}
	logger.Infof("extractLaunchArgs: launch info key (%d) passed in had no match, using fallback arguments", liKey)
	return fallback
}

func formatUnexpectedResponse(event string) error {
	// Formatting for general errors that occur during launcher client actions (displayed in msg box)
	return fmt.Errorf("Received an unexpected response when %s", event)
}

func buildArgs(cfg *Configuration, baseArgs string) string {
	largs := []string{baseArgs}
	if ConfAppendCustomArgs != "" {
		largs = append(largs, ConfAppendCustomArgs)
	}
	// Keep support for cmd-line (shortcut) Max FPS argument for backwards compatibility w/ previous version
	// If the "maxfps" cmd-line argument is specified, it takes precedence over the value configured in the UI
	if ConfMaxFPS != 0 {
		largs = append(largs, fmt.Sprintf("--set /Config/CONFIG/maxFpsValue %d", ConfMaxFPS))
	} else if cfg.Experimental.UseMaxFPSLimit {
		largs = append(largs, fmt.Sprintf("--set /Config/CONFIG/maxFpsValue %d", cfg.Experimental.MaxFPSLimit))
	}
	if cfg.Experimental.UseMaxFPSLimitMinimized {
		largs = append(largs, fmt.Sprintf("--set /Config/CONFIG/maxFpsValueMinimized %d", cfg.Experimental.MaxFPSLimitMinimized))
	}
	if cfg.Experimental.UseFPSSmoothing {
		largs = append(largs, fmt.Sprintf("--set /Config/CONFIG/enableFpsSmooth 1"))
	}
	return strings.Join(largs, " ")
}

func runQC(cfg *Configuration, baseArgs string) error {
	qc := exec.Command(cfg.Core.FilePath)
	qc.Dir = filepath.Dir(cfg.Core.FilePath)
	a := buildArgs(cfg, baseArgs)
	logger.Debugf("Final arguments: %s", a)
	// Handle arg quote-escaping manually (see golang issue #15566)
	qc.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CmdLine:       fmt.Sprintf(` %s`, a),
		CreationFlags: 0,
	}
	logger.Debug("Launching....")
	if err := qc.Start(); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error starting QC", GetCaller()), "error", err)
		return err
	}
	handlePostLaunch(cfg)
	return nil
}

func getProjectBranchBuildIdentifiers(r *EntitlementInfoResponse) (projectID int, buildID int, branchID int, err error) {
	for _, v := range r.Branches {
		if v.Project == qcProjectID && strings.EqualFold(v.Name, qcDefaultBranchIdentifier) {
			return v.Project, v.ID, v.Build, nil
		}
	}
	return 0, 0, 0, fmt.Errorf("QC build/branch identifiers were not found")
}

func handlePostLaunch(cfg *Configuration) {
	if cfg.Launcher.ExitOnLaunch {
		exitFromUI()
		return
	} else if cfg.Launcher.MinimizeOnLaunch {
		if cfg.Launcher.MinimizeToTray && qclauncherMainWindow.TrayIcon.Visible() {
			// Already minimized to tray and being launched from tray context menu item;  do nothing
			return
		}
		qclauncherMainWindow.minimize(true)
	}
}
