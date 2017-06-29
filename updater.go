// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type updateType int

const (
	updateAll updateType = iota
	updateQC
	updateLauncher
)

func Update(enforceHashIntegrity bool) bool {
	continueLaunch := true
	updateData, err := getLastUpdateTime()
	if err != nil {
		logUpdateError(err, updateQC, time.Now().Unix())
		return continueLaunch
	}
	// In the current beta state of the game, updates will be frequent, so verify QC against the
	// latest version from Bethesda on every launch.
	checkForQCUpdate(enforceHashIntegrity)

	if isUpdateDue(updateData.LastLauncherUpdateTime) {
		continueLaunch = checkForLauncherUpdate()
	}
	return continueLaunch
}

func checkForQCUpdate(enforceHashIntegrity bool) {
	now := time.Now().Unix()
	t := now
	qcHash, err := getHash()
	if err != nil {
		logUpdateError(err, updateQC, now)
		return
	}
	lc := newLauncherClient(10)
	qcUpdateInfo, err := lc.getQCUpdateInfo()
	if err != nil {
		logUpdateError(err, updateQC, now)
		return
	}
	if !strings.EqualFold(qcHash, qcUpdateInfo.Hash) {
		t = 0 // try next time
		// Only allow launching with the latest version of QC (default) unless enforcement is specifically disabled
		if !enforceHashIntegrity {
			ShowWarningMsg("Warning",
				"Your QC did not match the newest version from Bethesda. Please run the Bethesda Launcher to update Quake Champions! Launching anyway, but it will probably be unsuccessful.")
		} else {
			logger.Errorf("QC hash mismatch, locally got: %s, latest Bethesda version has: %s", strings.ToUpper(qcHash), qcUpdateInfo.Hash)
			ShowErrorMsg("Error",
				"Your QC did not match the newest version from Bethesda. Please run the Bethesda Launcher to update Quake Champions! Exiting.")
			err = updateLastCheckTime(updateQC, t)
			if err != nil {
				logger.Errorw("checkForQCUpdate: Error updating time info", "error", err)
			}
			os.Exit(1)
		}
	}
	err = updateLastCheckTime(updateQC, t)
	if err != nil {
		logger.Errorw("checkForQCUpdate: Error updating time info", "error", err)
	}
}

func checkForLauncherUpdate() bool {
	now := time.Now().Unix()
	continueLaunch := true
	lc := newLauncherClient(10)
	linfo, err := lc.getLauncherUpdateInfo()
	if err != nil {
		logUpdateError(err, updateLauncher, now)
		return continueLaunch
	}
	if version < linfo.LatestVersion {
		err = updateLastCheckTime(updateLauncher, now)
		if err != nil {
			logger.Errorw("checkForLauncherUpdate: Error updating time info", "error", err)
		}
		continueLaunch = showUpdateMsg(&LauncherUpdateInfo{LatestVersion: linfo.LatestVersion, Date: linfo.Date, URL: linfo.URL})
	}
	return continueLaunch
}

func logUpdateError(originalErr error, ut updateType, unixTime int64) {
	logger.Errorw("Update time info error", "error", originalErr)
	err := updateLastCheckTime(ut, unixTime)
	if err != nil {
		logger.Errorw("Error updating time info", "error", err)
	}
}

func isUpdateDue(unixTime int64) bool {
	return unixTime == 0 || int64(time.Since(time.Unix(unixTime, 0)).Seconds()) > ConfUpdateInterval
}

func getHash() (string, error) {
	qcOpts, err := getQCOptions()
	if err != nil {
		logger.Errorw("getHash: Error getting QC options", "error", err)
		return "", err
	}
	qc, err := os.Open(qcOpts.QCFilePath)
	if err != nil {
		logger.Errorw("getHash: Error opening QC executable", "error", err)
		return "", err
	}
	defer qc.Close()
	h := sha256.New224()
	if _, err := io.Copy(h, qc); err != nil {
		logger.Errorw("getHash: Error calculating QC executable hash", "error", err)
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
