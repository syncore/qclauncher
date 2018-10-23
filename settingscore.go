// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"

	bolt "github.com/coreos/bbolt"
)

type QCCoreSettings struct {
	Username string
	Password string
	FilePath string
	Language string
	FP       string
}

func (s *QCCoreSettings) get(ls *LauncherStore) error {
	ls.checkDataFile(false)
	if err := ls.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSettings))
		decerr := s.decode(b.Get([]byte(keyQCCoreSettings)), b.Get([]byte(keyTokenKey)))
		if decerr != nil {
			logger.Errorw(fmt.Sprintf("%s: error decoding QC core settings from datastore during get operation", GetCaller()),
				"error", decerr)
			return decerr
		}
		return nil
	}); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting QC core settings from datastore", GetCaller()), "error", err)
		return nil
	}
	return nil
}

func (s *QCCoreSettings) save(ls *LauncherStore) error {
	if err := ls.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketSettings))
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating settings bucket in datastore during save operation",
				GetCaller()), "error", err)
			return err
		}
		s.FP = tmpFp
		encoded, err := s.encode()
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error encoding QC core settings during datastore save operation", GetCaller()),
				"error", err)
			return err
		}
		if err = b.Put([]byte(keyQCCoreSettings), encoded); err != nil {
			logger.Errorw(fmt.Sprintf("%s: error saving encoded QC core settings to datastore", GetCaller()), "error", err)
		}
		if err = b.Put([]byte(keyTokenKey), *tmpKey); err != nil {
			logger.Errorw(fmt.Sprintf("%s: error saving credential key to datastore", GetCaller()), "error", err)
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	t := &TokenAuth{Token: tmpToken}
	return t.save(ls)
}

func (s *QCCoreSettings) decode(data, key []byte) error {
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&s); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decoding QC core settings data", GetCaller()), "error", err)
		return err
	}
	qcDecUser, err := decrypt(s.Username, &key)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decrypting username credential", GetCaller()), "error", err)
		return err
	}
	qcDecPass, err := decrypt(s.Password, &key)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decrypting password credential", GetCaller()), "error", err)
		return err
	}
	s.Username = qcDecUser
	s.Password = qcDecPass
	return nil
}

func (s *QCCoreSettings) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	encUser, err := encrypt(s.Username, tmpKey)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error encrypting username credential", GetCaller()), "error", err)
		return nil, err
	}
	encPass, err := encrypt(s.Password, tmpKey)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error encrypting password credential", GetCaller()), "error", err)
		return nil, err
	}
	s.Username = encUser
	s.Password = encPass
	if err = gob.NewEncoder(buf).Encode(s); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error encoding QC core settings data", GetCaller()), "error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *QCCoreSettings) validate() error {
	if s == nil {
		return errors.New("QC Core setting info was not entered")
	}
	if s.Username == "" {
		return errors.New("QC username must be specified")
	}
	if s.Password == "" {
		return errors.New("QC password must be specified")
	}
	if s.FilePath == "" {
		return errors.New("QC EXE location must be specified")
	}
	if !strings.Contains(strings.ToUpper(s.FilePath), strings.ToUpper(QCExe)) {
		return errors.New("Invalid QC EXE was specified")
	}
	if s.Language == "" {
		return errors.New("QC language must be specified")
	}
	if isFPOverride() {
		s.FP = ConfXSrcFp
	}
	fp, err := validateAccount(s.Username, s.Password, s.FP)
	if fp == "" {
		return errors.New("Unable to get required hardware fingerprint from Bethesda Launcher. Please try again.")
	}
	return err
}

func validateAccount(username, password, fp string) (string, error) {
	var fpErr error
	if !FileExists(GetDataFilePath()) {
		if fp == "" {
			fp, fpErr = getBNLFingerprint()
			if fpErr != nil {
				return "", fpErr
			}
			return fp, newLauncherClient(defTimeout).verifyCredentials(username, password)
		} else {
			tmpFp = fp
		}
	} else {
		cfg, err := GetConfiguration()
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error getting configuration during pre-save account validation",
				GetCaller()), "error", err)
			fp, fpErr = getBNLFingerprint()
			if fpErr != nil {
				return "", fpErr
			}
			return fp, newLauncherClient(defTimeout).verifyCredentials(username, password)
		}
		if isFPOverride() {
			fp = ConfXSrcFp
		} else if cfg.Core.FP != "" {
			fp = cfg.Core.FP
		} else {
			fp, fpErr = getBNLFingerprint()
			if fpErr != nil {
				return "", fpErr
			}
		}
		if cfg.Core.Username == username && cfg.Core.Password == password {
			token := &TokenAuth{}
			if err := Get(token); err != nil {
				return fp, newLauncherClient(defTimeout).verifyCredentials(username, password)
			}
			tmpKey = genKey()
			tmpToken = token.Token
			tmpFp = fp
			return fp, nil
		}
	}
	return fp, newLauncherClient(defTimeout).verifyCredentials(username, password)
}
