// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

// +build windows,amd64

package main

//go:generate goversioninfo -manifest "../../resources/qclauncher.manifest" -icon "../../resources/qclauncher.ico" -o "qclauncher_amd64.syso" -64 ../../resources/versioninfo.json
//go:generate go-bindata -pkg "resources" -o ../../resources/logo.go ../../resources/img

import (
	"flag"
	"fmt"
	"os"

	"github.com/syncore/qclauncher"
)

func init() {
	flag.BoolVar(&qclauncher.ConfLocal, "local", false, "Run in a local test environment")
	flag.BoolVar(&qclauncher.ConfDebug, "debug", false, "Log debug messages in addition to errors")
	flag.StringVar(&qclauncher.ConfLocalAddr, "localaddr", "localhost:30002", "Local endpoint host:port for test environment")
	flag.StringVar(&qclauncher.ConfXAppVer, "xappver", qclauncher.XAppDefVer, "Manually specify app version for request header")
	flag.StringVar(&qclauncher.ConfXLibVer, "xlibver", qclauncher.XLibDefVer, "Manually specify lib version for request header")
	flag.StringVar(&qclauncher.ConfAppendCustomArgs, "customargs", "", "Append the specified args to the launch args")
	flag.Int64Var(&qclauncher.ConfUpdateInterval, "updateinterval", 86400, "Time in seconds between checking for launcher updates") // 24 hours (86400)
	flag.BoolVar(&qclauncher.ConfSkipUpdates, "skipupdates", false, "Skip checking for QC and launcher updates")
	flag.BoolVar(&qclauncher.ConfEnforceHash, "enforcehash", true, "Enforce QC game hash checking (disabling is not recommended)")
	flag.IntVar(&qclauncher.ConfMaxFPS, "maxfps", 0, "Max value to limit FPS to (experimental)")
}

func main() {
	flag.Parse()
	qclauncher.Setup()
	execMain()
}

func execMain() {
	mainLogger := qclauncher.NewLogger()
	lmsg := fmt.Sprintf("An error occurred while executing the launch process. See %s for more information.", qclauncher.LogFile)
	running, procs, err := qclauncher.IsProcessRunning("QuakeChampions.exe")
	if err != nil {
		mainLogger.Errorw("main: Error checking running processes.", "error", err)
		qclauncher.ShowErrorMsg("Error", fmt.Sprintf("An error occurred while checking the running processes. See %s for more information.", qclauncher.LogFile))
		return
	}
	if running {
		qclauncher.ShowErrorMsg("Already running", fmt.Sprintf("The following are currently running: %s. Cannot continue. Exiting.", procs))
		return
	}
	if !qclauncher.FileExists(qclauncher.GetDataFilePath()) {
		qclauncher.OpenSettings()
		return
	}
	if qclauncher.ConfSkipUpdates {
		if err := qclauncher.Launch(); err != nil {
			mainLogger.Errorw("main: Error occurred while executing the launch process.", "error", err)
			qclauncher.ShowErrorMsg("Error", lmsg)
			os.Exit(1)
		}
		return
	}
	if qclauncher.Update(qclauncher.ConfEnforceHash) {
		if err := qclauncher.Launch(); err != nil {
			mainLogger.Errorw("main: Error occurred while executing the launch process.", "error", err)
			qclauncher.ShowErrorMsg("Error", lmsg)
			os.Exit(1)
		}
	}
}
