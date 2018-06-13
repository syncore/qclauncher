// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"
	"path/filepath"
)

const (
	LogFile            = "qclauncher.log"
	DataFile           = "data.qcl"
	LockFile           = "qcl.lock"
	ShowMainWindowFlag = "show"
	XAppDefVer         = "1.36.2"
	XLibDefVer         = "1.36.2"
	bDefBase           = "buildinfo.cdp.bethesda.net"
	sDefBase           = "services.bethesda.net"
	defTimeout         = 10
	version            = 1.04
)

var (
	ConfLocal            bool
	ConfDebug            bool
	ConfAppendCustomArgs string
	ConfLocalAddr        string
	ConfXAppVer          string
	ConfXLibVer          string
	ConfUpdateInterval   int64
	ConfSkipUpdates      bool
	ConfEnforceHash      bool
	ConfMaxFPS           int
	ConfBaseSvc          string
	ConfBaseBi           string
	ConfSrcFp            string
	ConfShowMainWindow   bool
	Lock                 *Single
)

func Setup() {
	setLogger()
	setLock()
	setBaseAddr()
	setSrcFp()
	setVersionInfo()
}

func GetDataFilePath() string {
	return filepath.Join(getExecutingPath(), DataFile)
}

func setBaseAddr() {
	if ConfLocal {
		addr := fmt.Sprintf("http://%s", ConfLocalAddr)
		ConfBaseBi, ConfBaseSvc = addr, addr
	} else {
		ConfBaseBi, ConfBaseSvc = fmt.Sprintf("https://%s", bDefBase), fmt.Sprintf("https://%s", sDefBase)
	}
}

func setLock() {
	if Lock == nil {
		Lock = NewSingle(LockFile)
	}
}

func setSrcFp() {
	ConfSrcFp = genFp()
}

func setVersionInfo() {
	if ConfXAppVer != XAppDefVer && ConfXLibVer != XLibDefVer {
		return
	}
	lc := newLauncherClient(7)
	uinfo, err := lc.getQCUpdateInfo()
	if err != nil {
		return
	}
	cachedQCUpdateInfo = uinfo
	if uinfo.BVer == "" {
		return
	}
	if ConfXAppVer == XAppDefVer {
		ConfXAppVer = uinfo.BVer
	}
	if ConfXLibVer == XLibDefVer {
		ConfXLibVer = uinfo.BVer
	}
}

func getLogFilePath() string {
	return filepath.Join(getExecutingPath(), LogFile)
}
