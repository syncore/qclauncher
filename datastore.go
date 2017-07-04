// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

// TODO: refactor/reduce the verbosity of this datastore file; it's crazy even by go standards & can be simplified

package qclauncher

import (
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"
)

var (
	qcCredBucket          = "qccb"
	qcOptBucket           = "qcob"
	updBucket             = "updb"
	qcUserKey             = "userqc"
	qcPassKey             = "passqc"
	rndencKey             = "rndenc"
	qcAuthTokenKey        = "authqc"
	qcFilePathKey         = "fpathqc"
	qcLangKey             = "langqc"
	lastUpdateQCKey       = "lastqcupd"
	lastUpdateLauncherKey = "lastlchupd"
	dfVerKey              = "dbver"
	dataFileIncompatible  = fmt.Sprintf("Your %s file is incompatible with this version of the launcher. Please delete %s and restart QCLauncher.",
		DataFile, DataFile)
	errUnableToOpenDatafile = fmt.Sprintf("Unable to open %s data file.", DataFile)
	tmpToken                string
	tmpKey                  *[]byte
)

const dataFileVersion int64 = 2

func updateAuthToken(isPreSaveVerification bool, authToken string) error {
	if isPreSaveVerification {
		// Data file won't exist on first-run credential verification; which is the entry point into
		// the data store, so save token & key in temp vars so they will be applied when the final
		// settings are saved.
		tmpToken = authToken
		tmpKey = genKey()
		return nil
	}
	return updateExistingAuthToken(authToken)
}

func updateExistingAuthToken(authToken string) error {
	checkDataFile(false)
	var k []byte
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("updateExistingAuthToken: error opening data file", "error", err)
		return err
	}
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(qcCredBucket))
		k = b.Get([]byte(rndencKey))
		return nil
	}); err != nil {
		logger.Errorw("updateExistingAuthToken: error getting key from datastore", "error", err)
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b, dberr := tx.CreateBucketIfNotExists([]byte(qcCredBucket))
		if dberr != nil {
			logger.Errorw("updateExistingAuthToken: error creating credential bucket", "error", dberr)
			return dberr
		}
		b = tx.Bucket([]byte(qcCredBucket))
		encToken, err := encrypt(authToken, &k)
		if err != nil {
			logger.Errorw("updateExistingAuthToken: error encrypting token credential", "error", err)
			return err
		}
		dberr = b.Put([]byte(qcAuthTokenKey), []byte(encToken))
		if dberr != nil {
			logger.Errorw("updateExistingAuthToken: error saving token to datastore", "error", dberr)
			return dberr
		}
		return nil
	})
}

func updateLastCheckTime(ut updateType, unixTime int64) error {
	checkDataFile(false)
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("updateLastCheckTime: error opening data file", "error", err)
		return err
	}
	defer db.Close()
	lastQCUpdateTime, lastLauncherUpdateTime := make([]byte, 8), make([]byte, 8)
	binary.LittleEndian.PutUint64(lastQCUpdateTime, uint64(unixTime))
	binary.LittleEndian.PutUint64(lastLauncherUpdateTime, uint64(unixTime))

	return db.Update(func(tx *bolt.Tx) error {
		b, dberr := tx.CreateBucketIfNotExists([]byte(updBucket))
		if dberr != nil {
			logger.Errorw("updateLastCheckTime: error creating update time info bucket", "error", dberr)
			return dberr
		}
		var uperr, lqerr, llerr error
		switch ut {
		case updateQC:
			uperr = b.Put([]byte(lastUpdateQCKey), lastQCUpdateTime)
		case updateLauncher:
			uperr = b.Put([]byte(lastUpdateLauncherKey), lastLauncherUpdateTime)
		case updateAll:
			lqerr = b.Put([]byte(lastUpdateQCKey), lastQCUpdateTime)
			llerr = b.Put([]byte(lastUpdateLauncherKey), lastLauncherUpdateTime)
		default:
			logger.Errorw("updateLastCheckTime: got unknown update type", "updateType", ut)
		}
		if uperr != nil {
			logger.Errorw("updateLastCheckTime: error saving update time info to datastore", "error", uperr)
			return uperr
		}
		if lqerr != nil {
			logger.Errorw("updateLastCheckTime: error saving last qc update (all) time info to datastore", "error", lqerr)
			return lqerr
		}
		if llerr != nil {
			logger.Errorw("updateLastCheckTime: error saving last launcher update (all) time info to datastore", "error", llerr)
			return llerr
		}
		return nil
	})
}

func getAuthToken() (string, error) {
	checkDataFile(false)
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("getAuthToken: error opening data file", "error", err)
		return "", err
	}
	defer db.Close()
	var k, qcEncAuthToken []byte
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(qcCredBucket))
		k = b.Get([]byte(rndencKey))
		qcEncAuthToken = b.Get([]byte(qcAuthTokenKey))
		return nil
	}); err != nil {
		logger.Errorw("getAuthToken: error getting token from datastore", "error", err)
		return "", err
	}
	qcDecAuthToken, err := decrypt(string(qcEncAuthToken), &k)
	if err != nil {
		logger.Errorw("getAuthToken: error decrypting token credential", "error", err)
		return "", err
	}
	return qcDecAuthToken, nil
}

func getUserCredentials() (*QCUserCredentials, error) {
	checkDataFile(false)
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("getUserCredentials: error opening data file", "error", err)
		return nil, err
	}
	defer db.Close()
	var k, qcEncUser, qcEncPass, qcEncAuthToken []byte
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(qcCredBucket))
		k = b.Get([]byte(rndencKey))
		qcEncUser = b.Get([]byte(qcUserKey))
		qcEncPass = b.Get([]byte(qcPassKey))
		qcEncAuthToken = b.Get([]byte(qcAuthTokenKey))
		return nil
	}); err != nil {
		logger.Errorw("getUserCredentials: error getting credentials from datastore", "error", err)
		return nil, err
	}
	qcDecUser, err := decrypt(string(qcEncUser), &k)
	if err != nil {
		logger.Errorw("getUserCredentials: error decrypting username credential", "error", err)
		return nil, err
	}
	qcDecPass, err := decrypt(string(qcEncPass), &k)
	if err != nil {
		logger.Errorw("getUserCredentials: error decrypting password credential", "error", err)
		return nil, err
	}
	qcDecToken, err := decrypt(string(qcEncAuthToken), &k)
	if err != nil {
		logger.Errorw("getUserCredentials: error decrypting token credential", "error", err)
		return nil, err
	}
	return &QCUserCredentials{
		Username: qcDecUser,
		Password: qcDecPass,
		Token:    qcDecToken,
	}, nil
}

func getLastUpdateTime() (*UpdateTime, error) {
	checkDataFile(false)
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("getLastUpdateTime: error opening data file", "error", err)
		return &UpdateTime{}, err
	}
	defer db.Close()
	var lqut, llut []byte
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(updBucket))
		lqut = b.Get([]byte(lastUpdateQCKey))
		llut = b.Get([]byte(lastUpdateLauncherKey))
		return nil
	}); err != nil {
		logger.Errorw("getLastUpdateTime: error getting update time info from datastore", "error", err)
		return &UpdateTime{}, err
	}
	qcLast, lchLast := int64(binary.LittleEndian.Uint64(lqut)), int64(binary.LittleEndian.Uint64(llut))
	return &UpdateTime{LastQCUpdateTime: qcLast, LastLauncherUpdateTime: lchLast}, nil
}

func getQCOptions() (*QCOptions, error) {
	checkDataFile(false)
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("getQCOptions: error opening data file", "error", err)
		return nil, err
	}
	defer db.Close()
	var qcfp, qclang []byte
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(qcOptBucket))
		qcfp = b.Get([]byte(qcFilePathKey))
		qclang = b.Get([]byte(qcLangKey))
		return nil
	}); err != nil {
		logger.Errorw("getQCOptions: error getting QC options from datastore", "error", err)
		return nil, err
	}
	return &QCOptions{QCFilePath: string(qcfp), QCLanguage: string(qclang)}, nil
}

func checkDataFile(skipVersionCheck bool) {
	if !FileExists(GetDataFilePath()) {
		logger.Error("checkDataFile: data file does not exist")
		ShowFatalErrorMsg("Error", fmt.Sprintf("Your %s file is missing. Try re-running QCLauncher.", DataFile))
	}
	if skipVersionCheck {
		return
	}
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("getLastUpdateTime: error opening data file", "error", err)
		db.Close()
		ShowFatalErrorMsg("Error", errUnableToOpenDatafile)
	}
	defer db.Close()
	var v []byte
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(updBucket))
		v = b.Get([]byte(dfVerKey))
		return nil
	}); err != nil {
		logger.Errorw("checkDataFile: error getting data file version from datastore", "error", err)
		db.Close()
		ShowFatalErrorMsg("Error", fmt.Sprintf("Unable to determine %s data file version. Please delete %s and re-launch.",
			DataFile, DataFile))
	}
	savedVer := int64(binary.LittleEndian.Uint64(v))
	if savedVer != dataFileVersion {
		logger.Errorf("checkDataFile: data file version incompatibility found. Detected: %d, needed: %d", savedVer,
			dataFileVersion)
		db.Close()
		ShowFatalErrorMsg("Error", dataFileIncompatible)
	}
}

func saveQCOptions(qco *QCOptions) error {
	checkDataFile(true)
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("saveQCOptions: error opening data file", "error", err)
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b, dberr := tx.CreateBucketIfNotExists([]byte(qcOptBucket))
		if dberr != nil {
			logger.Errorw("saveQCOptions: error creating QC options bucket", "error", dberr)
			return dberr
		}
		dberr = b.Put([]byte(qcFilePathKey), []byte(qco.QCFilePath))
		if dberr != nil {
			logger.Errorw("saveQCOptions: error saving QC filepath option to datastore", "error", dberr)
			return dberr
		}
		dberr = b.Put([]byte(qcLangKey), []byte(qco.QCLanguage))
		if dberr != nil {
			logger.Errorw("saveQCOptions: error saving QC language option to datastore", "error", dberr)
			return dberr
		}
		return nil
	})
}

func saveUserCredentials(quc *QCUserCredentials) error {
	checkDataFile(true)
	encUser, err := encrypt(quc.Username, tmpKey)
	if err != nil {
		logger.Errorw("saveUserCredentials: error encrypting username credential", "error", err)
		return err
	}
	encPass, err := encrypt(quc.Password, tmpKey)
	if err != nil {
		logger.Errorw("saveUserCredentials: error encrypting password credential", "error", err)
		return err
	}
	encToken, err := encrypt(tmpToken, tmpKey)
	if err != nil {
		logger.Errorw("saveUserCredentials: error encrypting token credential", "error", err)
		return err
	}
	db, err := bolt.Open(GetDataFilePath(), 0600, nil)
	if err != nil {
		logger.Errorw("saveUserCredentials: error opening data file", "error", err)
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b, dberr := tx.CreateBucketIfNotExists([]byte(qcCredBucket))
		if dberr != nil {
			logger.Errorw("saveUserCredentials: error creating credentials bucket", "error", dberr)
			return dberr
		}
		dberr = b.Put([]byte(rndencKey), *tmpKey)
		if dberr != nil {
			logger.Errorw("saveUserCredentials: error saving credential key to datastore", "error", dberr)
			return dberr
		}
		dberr = b.Put([]byte(qcUserKey), []byte(encUser))
		if dberr != nil {
			logger.Errorw("saveUserCredentials: error saving username credential to datastore", "error", dberr)
			return dberr
		}
		dberr = b.Put([]byte(qcPassKey), []byte(encPass))
		if dberr != nil {
			logger.Errorw("saveUserCredentials: error saving password credential to datastore", "error", dberr)
			return dberr
		}
		dberr = b.Put([]byte(qcAuthTokenKey), []byte(encToken))
		if dberr != nil {
			logger.Errorw("saveUserCredentials: error saving auth token to datastore", "error", dberr)
			return dberr
		}
		b, dberr = tx.CreateBucketIfNotExists([]byte(updBucket))
		if dberr != nil {
			logger.Errorw("saveUserCredentials: error creating update time info bucket", "error", dberr)
			return dberr
		}
		dfv := make([]byte, 8)
		binary.LittleEndian.PutUint64(dfv, uint64(dataFileVersion))
		dberr = b.Put([]byte(dfVerKey), dfv)
		if dberr != nil {
			logger.Errorw("saveUserCredentials: error saving data file version to datastore", "error", dberr)
			return dberr
		}
		return nil
	})
}
