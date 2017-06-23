// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"bytes"
	"fmt"
	"image/png"
	"os"

	"github.com/lxn/walk"
	"github.com/lxn/win"
	"github.com/skratchdot/open-golang/open"
	"github.com/syncore/qclauncher/resources"
)

func ShowErrorMsg(title, message string) {
	walk.MsgBox(nil, title, message, walk.MsgBoxIconError)
}

func ShowFatalErrorMsg(title, message string) {
	walk.MsgBox(nil, title, message, walk.MsgBoxIconError)
	os.Exit(1)
}

func ShowWarningMsg(title, message string) {
	walk.MsgBox(nil, title, message, walk.MsgBoxIconWarning)
}

func showUpdateMsg(updateInfo *LauncherUpdateInfo) bool {
	msg := fmt.Sprintf("An update is available for QCLauncher!\nYour version: %.2f\nLatest version: %.2f\nDate: %s\nClick \"Yes\" to exit and go to the download site or \"No\" to continue.",
		version, updateInfo.LatestVersion, updateInfo.Date.Format("Mon Jan 2 15:04:05 MST 2006"))
	result := walk.MsgBox(nil, "QCLauncher Update", msg, walk.MsgBoxYesNo)
	if result == win.IDYES {
		err := open.Start(updateInfo.URL)
		if err != nil {
			logger.Errorw("showUpdateMsg: error launching web browser to load update page", "error", err)
			ShowErrorMsg("Error", "Unable to open web browser to load update page")
		}
		return false
	} else if result == win.IDNO || result == 0 {
		return true
	}
	return true
}

func getAppIcon() (icon *walk.Icon) {
	for i := 0; i < 128; i++ {
		if icon, err := walk.NewIconFromResourceId(i); err == nil {
			return icon
		}
	}
	logger.Error("getAppIcon: Could not find icon from resource id")
	return nil
}

func loadAppLogo(assetPath string) (*walk.Bitmap, error) {
	i, err := resources.Asset(assetPath)
	if err != nil {
		logger.Errorw("loadAppLogo: error loading logo image asset", "error", err)
		return nil, err
	}
	img, err := png.Decode(bytes.NewBuffer(i))
	if err != nil {
		logger.Errorw("loadLogo: error decoding logo image asset", "error", err)
		return nil, err
	}
	bm, err := walk.NewBitmapFromImage(img)
	if err != nil {
		logger.Errorw("loadAppLogo: error converting logo image to usable bitmap", "error", err)
		return nil, err
	}
	return bm, nil
}
