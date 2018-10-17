// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"encoding/binary"
	"fmt"

	bolt "github.com/coreos/bbolt"
)

type LauncherStore struct {
	*bolt.DB
}

type Storable interface {
	get(*LauncherStore) error
	save(*LauncherStore) error
}

const (
	dataFileVersion           int64 = 4
	bucketSettings                  = "sb"
	bucketLastUpdate                = "lub"
	keyQCCoreSettings               = "core"
	keyQCExperimentalSettings       = "exp"
	keyLauncherSettings             = "lch"
	keyTokenAuth                    = "atkn"
	keyTokenKey                     = "rndenc"
	keyLastUpdateQC                 = "luqc"
	keyLastUpdateLauncher           = "lulc"
	keyDfVer                        = "dfver"
)

var (
	dataFileIncompatible = fmt.Sprintf(
		"Your %s file is from an old version of QCLauncher and is not supported by this version. You must restart QCLauncher to reset your settings.",
		DataFile)
	tmpToken string
	tmpKey   *[]byte
	tmpFp    string
)

func Save(s Storable) error {
	ls, err := newLauncherDataStore()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error initializing datastore during save operation", GetCaller()), "error", err)
		return err
	}
	defer ls.Close()
	return s.save(ls)
}

func Get(s Storable) error {
	ls, err := newLauncherDataStore()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error initializing datastore during get operation", GetCaller()), "error", err)
		return err
	}
	defer ls.Close()
	return s.get(ls)
}

func newLauncherDataStore() (*LauncherStore, error) {
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error opening data file", GetCaller()), "error", err)
		if db != nil {
			db.Close()
		}
		ShowFatalErrorMsg("Error", fmt.Sprintf("Unable to open file: %s. If %s exists, then delete it and try restarting QCLauncher.",
			DataFile, DataFile), nil)
		return nil, err
	}
	return &LauncherStore{db}, db.Update(func(tx *bolt.Tx) error {
		_, dberr := tx.CreateBucketIfNotExists([]byte(bucketSettings))
		if dberr != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating settings bucket", GetCaller()), "error", dberr)
			if dberr = DeleteFile(GetDataFilePath()); dberr != nil {
				logger.Errorw(fmt.Sprintf("%s: error deleting datafile during failed settings bucket creation", GetCaller()),
					"error", dberr)
			}
			return dberr
		}
		_, dberr = tx.CreateBucketIfNotExists([]byte(bucketLastUpdate))
		if dberr != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating update time info bucket", GetCaller()), "error", dberr)
			dberr = DeleteFile(GetDataFilePath())
			if dberr != nil {
				logger.Errorw(fmt.Sprintf("%s: error deleting datafile during update time bucket creation", GetCaller()),
					"error", dberr)
			}
			return dberr
		}
		return nil
	})
}

func (ls *LauncherStore) checkDataFile(skipVersionCheck bool) {
	if !FileExists(GetDataFilePath()) {
		logger.Error(fmt.Sprintf("%s: data file does not exist", GetCaller()))
		ShowFatalErrorMsg("Error", fmt.Sprintf("Your %s file is missing. Try re-running QCLauncher.", DataFile), nil)
		return
	}
	if skipVersionCheck {
		return
	}
	var v []byte
	if err := ls.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketLastUpdate))
		v = b.Get([]byte(keyDfVer))
		return nil
	}); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting data file version from datastore", GetCaller()), "error", err)
		if ls != nil && ls.DB != nil {
			ls.Close()
		}
		DeleteConfiguration(true)
		ShowFatalErrorMsg("Error", "Could not determine data file version. Please restart QCLauncher to reset your settings.", nil)
		return
	}
	if len(v) == 0 {
		logger.Error(fmt.Sprintf("%s: error getting data file version from datastore (may be an old version of %s)", GetCaller(),
			DataFile))
		if ls != nil && ls.DB != nil {
			ls.Close()
		}
		DeleteConfiguration(true)
		ShowFatalErrorMsg("Error", dataFileIncompatible, nil)
		return
	}
	savedVer := int64(binary.LittleEndian.Uint64(v))
	if savedVer != dataFileVersion {
		logger.Error(fmt.Sprintf("%s: data file version incompatibility found. Detected: %d, needed: %d", GetCaller(),
			savedVer, dataFileVersion))
		if ls != nil && ls.DB != nil {
			ls.Close()
		}
		DeleteConfiguration(true)
		ShowFatalErrorMsg("Error", dataFileIncompatible, nil)
		return
	}
}
