// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"bytes"
	"fmt"
	"image/png"
	"os/exec"
	"strconv"

	"github.com/lxn/walk"
	"github.com/lxn/win"
	"github.com/syncore/qclauncher/resources"
)

var (
	UILaunchErrorMsg = fmt.Sprintf("An error occurred while executing the launch process. See %s for more information. For support visit: https://github.com/syncore/qclauncher/issues",
		LogFile)
	title = fmt.Sprintf("QCLauncher %.2f by syncore", version)
)

func ShowErrorMsg(title, message string, owner walk.Form) {
	walk.MsgBox(owner, title, message, walk.MsgBoxIconError)
}

func ShowFatalErrorMsg(title, message string, owner walk.Form) {
	walk.MsgBox(owner, title, message, walk.MsgBoxIconError)
	Exit(1)
}

func ShowWarningMsg(title, message string, owner walk.Form) {
	walk.MsgBox(owner, title, message, walk.MsgBoxIconWarning)
}

func ShowInfoMsg(title, message string, owner walk.Form) {
	walk.MsgBox(owner, title, message, walk.MsgBoxIconInformation)
}

func ShowQCRunningMsg(pid int) bool {
	result := walk.MsgBox(nil, "Already Running",
		"Quake Champions is already running. Should QCLauncher exit Quake Champions for you?", walk.MsgBoxYesNo)
	if result == win.IDYES {
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
		if err := cmd.Run(); err != nil {
			ShowFatalErrorMsg("Error", "Unable to exit Quake Champions. Please exit QC and restart QCLauncher.", nil)
		}
		return true
	} else if result == win.IDNO || result == 0 {
		return false
	}
	return false
}

func getAppIcon() (icon *walk.Icon) {
	for i := 0; i < 128; i++ {
		if icon, err := walk.NewIconFromResourceId(i); err == nil {
			return icon
		}
	}
	logger.Error(fmt.Sprintf("%s: Could not find icon from resource id", GetCaller()))
	return nil
}

func loadAppLogo(assetPath string) (*walk.Bitmap, error) {
	i, err := resources.Asset(assetPath)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error loading logo image asset", GetCaller()), "error", err)
		return nil, err
	}
	img, err := png.Decode(bytes.NewBuffer(i))
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decoding logo image asset", GetCaller()), "error", err)
		return nil, err
	}
	bm, err := walk.NewBitmapFromImage(img)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error converting logo image to usable bitmap", GetCaller()), "error", err)
		return nil, err
	}
	return bm, nil
}

func exitFromUI() {
	Lock.Unlock()
	walk.App().Exit(0)
}
