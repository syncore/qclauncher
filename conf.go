// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"
	"path/filepath"
)

const (
	LogFile    = "qclauncher.log"
	DataFile   = "data.qcl"
	XAppDefVer = "1.20.5"
	XLibDefVer = "1.20.5"
	bDefBase   = "buildinfo.cdp.bethesda.net"
	sDefBase   = "services.bethesda.net"
	defTimeout = 10
	version    = 1.02
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
)

func Setup() {
	setLogger()
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
