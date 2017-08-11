// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"github.com/skratchdot/open-golang/open"
)

type UpdateTime struct {
	LastQCUpdateTime       int64
	LastLauncherUpdateTime int64
}

type LauncherUpdateInfo struct {
	LatestVersion float32
	Date          time.Time
	URL           string
}

type UpdateType int

const (
	UpdateAll UpdateType = iota
	UpdateQC
	UpdateLauncher
)

var cachedQCUpdateInfo *UpdateQCResponse

func CheckUpdate(enforceHashIntegrity bool, ut UpdateType) error {
	ls, err := newLauncherDataStore()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error initializing data store", GetCaller()), "error", err)
		return nil // swallow & continue
	}
	updateData, err := ls.getLastUpdateTime()
	if err != nil {
		logUpdateError(err, ut, time.Now().Unix())
		return nil // swallow & continue
	}
	l := newLauncherClient(defTimeout)
	if ut == UpdateAll || (ut == UpdateLauncher && isUpdateDue(updateData.LastLauncherUpdateTime)) {
		if continueRunning := l.checkForLauncherUpdate(); !continueRunning {
			exitFromUI()
			return nil
		}
	}
	var qcuperr error
	if ut == UpdateQC || ut == UpdateAll {
		// Verify QC against the latest version (default) from Bethesda on every launch unless disabled
		qcuperr = l.checkForQCUpdate(enforceHashIntegrity)
		if IsErrHashMismatch(qcuperr) {
			qcuperr = &hashMismatchError{emsg: qcuperr.Error()} // for presentation to top-level caller
		} else {
			qcuperr = nil // swallow
		}
	}
	return qcuperr
}

func promptForLauncherUpdate(updateInfo *LauncherUpdateInfo) bool {
	msg := fmt.Sprintf("An update is available for QCLauncher!\nYour version: %.2f\nLatest version: %.2f\nDate: %s\nClick \"Yes\" to exit and go to the download site or \"No\" to continue.",
		version, updateInfo.LatestVersion, updateInfo.Date.Format("Mon Jan 2 15:04:05 MST 2006"))
	result := walk.MsgBox(nil, "QCLauncher Update", msg, walk.MsgBoxYesNo)
	if result == win.IDYES {
		ShowInfoMsg("Update", "QCLauncher will exit so that you can download and install the new version.", nil)
		err := open.Start(updateInfo.URL)
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error launching web browser to load update page", GetCaller()), "error", err)
			ShowErrorMsg("Error", "Unable to open web browser to load update page", nil)
			return true
		}
		return false
	}
	return true
}

func (lc *launcherClient) checkForQCUpdate(enforceHashIntegrity bool) error {
	var err error
	var qcUpdateInfo *UpdateQCResponse
	now := time.Now().Unix()
	t := now
	if cachedQCUpdateInfo == nil {
		logger.Debug("cachedQCUpdateInfo is nil, getting fresh update info")
		qcUpdateInfo, err = lc.getQCUpdateInfo()
		if err != nil {
			logUpdateError(err, UpdateQC, now)
			return err
		}
	} else {
		qcUpdateInfo = cachedQCUpdateInfo
		logger.Debugf("using already-present cachedQCUpdateInfo: %+v", cachedQCUpdateInfo)
	}
	h := []FileHash{}
	for _, fh := range qcUpdateInfo.Hashes {
		h = append(h, FileHash{File: strings.Replace(fh.File, "/", "\\", -1), Hash: fh.Hash})
	}
	cherr := compareHashes(h)
	if _, ok := cherr.(*hashMismatchError); ok {
		t = 0 // try next time
		// Only allow launching with the latest version of QC (default) unless enforcement is specifically disabled
		if !enforceHashIntegrity {
			ShowWarningMsg("Warning",
				"Your QC files did not match the newest versions from Bethesda. Please run the Bethesda Launcher to update Quake Champions! Launching anyway, but it will probably be unsuccessful.", nil)
		} else {
			logger.Error(cherr)
			if uerr := updateLastCheckTime(UpdateQC, t); uerr != nil {
				logger.Errorw(fmt.Sprintf("%s: error updating time info", GetCaller()), "error", uerr)
			}
			return &hashMismatchError{emsg: cherr.Error()}
		}
	}
	if err != nil {
		logUpdateError(err, UpdateQC, now)
		return err
	}
	err = updateLastCheckTime(UpdateQC, t)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error updating time info", GetCaller()), "error", err)
	}
	return nil
}

func (lc *launcherClient) checkForLauncherUpdate() bool {
	now := time.Now().Unix()
	ignore := true
	linfo, err := lc.getLauncherUpdateInfo()
	if err != nil {
		logUpdateError(err, UpdateLauncher, now)
		return ignore
	}
	if version < linfo.LatestVersion {
		err = updateLastCheckTime(UpdateLauncher, now)
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error updating time info", GetCaller()), "error", err)
		}
		ignore = promptForLauncherUpdate(&LauncherUpdateInfo{LatestVersion: linfo.LatestVersion, Date: linfo.Date, URL: linfo.URL})
	}
	return ignore
}

func logUpdateError(originalErr error, ut UpdateType, unixTime int64) {
	logger.Errorw(fmt.Sprintf("%s: Update time info error", GetCaller()), "error", originalErr)
	err := updateLastCheckTime(ut, unixTime)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error updating time info", GetCaller()), "error", err)
	}
}

func isUpdateDue(unixTime int64) bool {
	return unixTime == 0 || int64(time.Since(time.Unix(unixTime, 0)).Seconds()) > ConfUpdateInterval
}

func compareHashes(hashes []FileHash) error {
	matches := 0
	cfg, err := GetConfiguration()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting configuration", GetCaller()), "error", err)
		return err
	}
	p := strings.ToLower(cfg.Core.FilePath)
	for _, fh := range hashes {
		r := strings.Replace(p, "client\\bin\\pc\\quakechampions.exe", fh.File, -1)
		f, err := os.Open(r)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: error opening %s to compare hash: %s", GetCaller(), r, err))
			return err
		}
		defer f.Close()
		h := sha256.New224()
		if _, err := io.Copy(h, f); err != nil {
			logger.Error(fmt.Sprintf("%s: error calculating hash for %s: %s", GetCaller(), r, err))
			return err
		}
		localCalc := fmt.Sprintf("%x", h.Sum(nil))
		if strings.EqualFold(localCalc, fh.Hash) {
			matches++
			logger.Debugf("Got hash match for %s", r)
		} else {
			logger.Errorw(fmt.Sprintf("%s: File hash mismatch, %s local version hash: %s, latest Bethesda version hash: %s",
				GetCaller(), r, strings.ToUpper(localCalc), fh.Hash))
			return &hashMismatchError{
				emsg: "One or more of your QC files did not match the newest version from Bethesda. Please run the Bethesda Launcher to update Quake Champions!"}
		}
	}
	if matches != len(hashes) {
		return &hashMismatchError{emsg: "All file hashes did not match latest Bethesda versions"}
	}
	return nil
}

func (ls *LauncherStore) getLastUpdateTime() (*UpdateTime, error) {
	ls.checkDataFile(false)
	defer ls.Close()
	var lqut, llut []byte
	if err := ls.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketLastUpdate))
		lqut = b.Get([]byte(keyLastUpdateQC))
		llut = b.Get([]byte(keyLastUpdateLauncher))
		return nil
	}); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting update time info from datastore", GetCaller()), "error", err)
		return &UpdateTime{}, err
	}
	qcLast, lchLast := int64(binary.LittleEndian.Uint64(lqut)), int64(binary.LittleEndian.Uint64(llut))
	return &UpdateTime{LastQCUpdateTime: qcLast, LastLauncherUpdateTime: lchLast}, nil
}

func updateLastCheckTime(ut UpdateType, unixTime int64) error {
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error opening data file", GetCaller()), "error", err)
		return err
	}
	defer db.Close()
	lastQCUpdateTime, lastLauncherUpdateTime := make([]byte, 8), make([]byte, 8)
	binary.LittleEndian.PutUint64(lastQCUpdateTime, uint64(unixTime))
	binary.LittleEndian.PutUint64(lastLauncherUpdateTime, uint64(unixTime))

	return db.Update(func(tx *bolt.Tx) error {
		b, dberr := tx.CreateBucketIfNotExists([]byte(bucketLastUpdate))
		if dberr != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating update time info bucket", GetCaller()), "error", dberr)
			return dberr
		}
		var uperr, lqerr, llerr error
		switch ut {
		case UpdateQC:
			uperr = b.Put([]byte(keyLastUpdateQC), lastQCUpdateTime)
		case UpdateLauncher:
			uperr = b.Put([]byte(keyLastUpdateLauncher), lastLauncherUpdateTime)
		case UpdateAll:
			lqerr = b.Put([]byte(keyLastUpdateQC), lastQCUpdateTime)
			// do not reset the existing launcher last update time when saving new settings
			if existingll := b.Get([]byte(keyLastUpdateLauncher)); len(existingll) != 0 {
				llerr = b.Put([]byte(keyLastUpdateLauncher), existingll)
			} else {
				llerr = b.Put([]byte(keyLastUpdateLauncher), lastLauncherUpdateTime)
			}
		default:
			logger.Errorw(fmt.Sprintf("%s: got unknown update type", GetCaller()), "updateType", ut)
		}
		if uperr != nil {
			logger.Errorw(fmt.Sprintf("%s: error saving update time info to datastore", GetCaller()), "error", uperr)
			return uperr
		}
		if lqerr != nil {
			logger.Errorw(fmt.Sprintf("%s: error saving last qc update (all) time info to datastore", GetCaller()), "error", lqerr)
			return lqerr
		}
		if llerr != nil {
			logger.Errorw(fmt.Sprintf("%s: error saving last launcher update (all) time info to datastore", GetCaller()), "error", llerr)
			return llerr
		}
		return nil
	})
}
