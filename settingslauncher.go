// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

type LauncherSettings struct {
	AutoStartQC       bool
	ExitOnLaunch      bool
	MinimizeOnLaunch  bool
	MinimizeToTray    bool
	SetAsNonSteamGame bool
}

func (s *LauncherSettings) get(ls *LauncherStore) error {
	ls.checkDataFile(false)
	if err := ls.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSettings))
		decerr := s.decode(b.Get([]byte(keyLauncherSettings)))
		if decerr != nil {
			logger.Errorw(fmt.Sprintf("%s: error decoding launcher settings from datastore during get operation", GetCaller()),
				"error", decerr)
		}
		return nil
	}); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting launcher settings from datastore", GetCaller()), "error", err)
		return err
	}
	return nil
}

func (s *LauncherSettings) save(ls *LauncherStore) error {
	return ls.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketSettings))
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating settings bucket in datastore during save operation",
				GetCaller()), "error", err)
			return err
		}
		encoded, err := s.encode()
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error encoding launcher settings during datastore save operation", GetCaller()),
				"error", err)
			return err
		}
		err = b.Put([]byte(keyLauncherSettings), encoded)
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error saving encoded launcher settings to datastore", GetCaller()), "error", err)
			return err
		}
		b, err = tx.CreateBucketIfNotExists([]byte(bucketLastUpdate))
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating update bucket for saving data file version to datastore", GetCaller()),
				"error", err)
			return err
		}
		dfv := make([]byte, 8)
		binary.LittleEndian.PutUint64(dfv, uint64(dataFileVersion))
		err = b.Put([]byte(keyDfVer), dfv)
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error saving data file version to datastore", GetCaller()), "error", err)
			return err
		}
		return nil
	})
}

func (s *LauncherSettings) decode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&s)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decoding launcher settings data", GetCaller()), "error", err)
		return err
	}
	return nil
}

func (s *LauncherSettings) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error encoding launcher settings data", GetCaller()), "error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *LauncherSettings) validate() error {
	if s == nil {
		return errors.New("Launcher setting info was not entered")
	}
	return nil
}
