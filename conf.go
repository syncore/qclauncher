// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"
	"path/filepath"
)

const (
	LogFile  = "qclauncher.log"
	DataFile = "data.qcl"
	bDefBase = "buildinfo.cdp.bethesda.net"
	sDefBase = "services.bethesda.net"
	version  = 1.01
)

var (
	ConfLocal          bool
	ConfDebug          bool
	ConfLocalAddr      string
	ConfXAppVer        string
	ConfXLibVer        string
	ConfUpdateInterval int64
	ConfSkipUpdates    bool
	ConfEnforceHash    bool
	ConfBaseSvc        string
	ConfBaseBi         string
	ConfSrcFp          string
)

func Setup() {
	setLogger()
	setBaseAddr()
	setSrcFp()
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

func getLogFilePath() string {
	return filepath.Join(getExecutingPath(), LogFile)
}
