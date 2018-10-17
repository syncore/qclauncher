// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

// +build windows,amd64

package main

//go:generate goversioninfo -manifest "../../resources/qclauncher.manifest" -icon "../../resources/qclauncher.ico" -o "qclauncher_amd64.syso" -64 ../../resources/versioninfo.json
//go:generate go-bindata -pkg "resources" -o ../../resources/res.go ../../resources/img ../../resources/bin

import (
	"flag"
	"fmt"

	"github.com/syncore/qclauncher"
)

func init() {
	flag.BoolVar(&qclauncher.ConfLocal, "local", false, "Run in a local test environment")
	flag.BoolVar(&qclauncher.ConfDebug, "debug", false, "Log debug messages in addition to errors")
	flag.StringVar(&qclauncher.ConfLocalAddr, "localaddr", "localhost:30002", "Local endpoint host:port for test environment")
	flag.StringVar(&qclauncher.ConfXAppVer, "xappver", qclauncher.XAppDefVer, "Manually specify app version for request header")
	flag.StringVar(&qclauncher.ConfXLibVer, "xlibver", qclauncher.XLibDefVer, "Manually specify lib version for request header")
	flag.StringVar(&qclauncher.ConfXSrcFp, "fp", qclauncher.XSrcFpDef, "Manually specify Bethesda hardware fingerprint for request header")
	flag.StringVar(&qclauncher.ConfAppendCustomArgs, "customargs", "", "Append the specified args to the launch args")
	flag.Int64Var(&qclauncher.ConfUpdateInterval, "updateinterval", 86400, "Time in seconds between checking for launcher updates") // 24 hours (86400)
	flag.BoolVar(&qclauncher.ConfSkipUpdates, "skipupdates", false, "Skip checking for QC and launcher updates")
	flag.BoolVar(&qclauncher.ConfEnforceHash, "enforcehash", true, "Enforce QC game hash checking (disabling is not recommended)")
	flag.IntVar(&qclauncher.ConfMaxFPS, "maxfps", 0, "Max value to limit FPS to (experimental)")
	flag.BoolVar(&qclauncher.ConfShowMainWindow, "show", false, "Restore the QCLauncher main UI window")
	flag.BoolVar(&qclauncher.ConfUseEntitlementAPI, "entitlement", false, "Use Bethesda.net entitlement API")
}

func main() {
	flag.Parse()
	qclauncher.Setup()
	execMain()
}

func execMain() {
	err := qclauncher.Lock.Lock()
	if qclauncher.IsErrAlreadyRunning(err) {
		qclauncher.ShowFatalErrorMsg("Error", fmt.Sprintf("%s If this is an error, delete the %s file and try again.",
			err.Error(), qclauncher.Lock.Filename(false)), nil)
		return
	}
	defer qclauncher.Lock.Unlock()
	mainlogger := qclauncher.NewLogger()
	running, _, _, _, namepids, err := qclauncher.IsProcessRunning(qclauncher.QCExe)
	if err != nil {
		mainlogger.Errorw(fmt.Sprintf("%s: error checking running processes", qclauncher.GetCaller()), "error", err)
		qclauncher.ShowErrorMsg("Error", "Unable to enumerate processes to determine if Quake Champions is already running", nil)
		return
	}
	if running {
		if willClose := qclauncher.ShowQCRunningMsg(namepids[qclauncher.QCExe]); !willClose {
			qclauncher.ShowErrorMsg("Already running", "Please close Quake Champions and then re-run QCLauncher.", nil)
			return
		}
	}
	if !qclauncher.FileExists(qclauncher.GetDataFilePath()) {
		qclauncher.LoadUI(qclauncher.GetEmptyConfiguration())
		return
	}
	cfg, err := qclauncher.GetConfiguration()
	if err != nil {
		qclauncher.ShowErrorMsg("Error", "An error occurred when retrieving your settings. Resetting.", nil)
		qclauncher.DeleteConfiguration(false)
		qclauncher.LoadUI(qclauncher.GetEmptyConfiguration())
		return
	}
	if !qclauncher.ConfSkipUpdates {
		// param of type UpdateLauncher to this call throws no error
		_ = qclauncher.CheckUpdate(qclauncher.ConfEnforceHash, qclauncher.UpdateLauncher)
	}
	if qclauncher.ConfUseEntitlementAPI {
		qclauncher.UseEntitlementAPI = false
	} else {
		qclauncher.SetEntitlementAPI()
	}
	if qclauncher.ConfShowMainWindow {
		qclauncher.LoadUI(cfg)
		return
	}
	if cfg.Launcher.AutoStartQC {
		if err := qclauncher.Launch(); err != nil {
			mainlogger.Errorw(fmt.Sprintf("%s: %s", qclauncher.GetCaller(), "error occurred while executing the launch process."),
				"error", err)
			if qclauncher.IsErrAlreadyRunning(err) || qclauncher.IsErrHashMismatch(err) || qclauncher.IsErrAuthFailed(err) {
				qclauncher.ShowErrorMsg("Error", err.Error(), nil)
			} else {
				qclauncher.ShowErrorMsg("Error", qclauncher.UILaunchErrorMsg, nil)
			}
			qclauncher.Exit(1)
		}
		return
	}
	qclauncher.LoadUI(cfg)
}
