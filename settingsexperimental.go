// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	bolt "github.com/coreos/bbolt"
)

type QCExperimentalSettings struct {
	UseMaxFPSLimit          bool
	UseMaxFPSLimitMinimized bool
	UseFPSSmoothing         bool
	MaxFPSLimit             int
	MaxFPSLimitMinimized    int
}

func (s *QCExperimentalSettings) get(ls *LauncherStore) error {
	ls.checkDataFile(false)
	if err := ls.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSettings))
		decerr := s.decode(b.Get([]byte(keyQCExperimentalSettings)))
		if decerr != nil {
			logger.Errorw(fmt.Sprintf("%s: error decoding QC experimental settings from datastore during get operation", GetCaller()),
				"error", decerr)
			return decerr
		}
		return nil
	}); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting QC experimental settings from datastore", GetCaller()), "error", err)
		return err
	}
	return nil
}

func (s *QCExperimentalSettings) save(ls *LauncherStore) error {
	return ls.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketSettings))
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating settings bucket in datastore during save operation",
				GetCaller()), "error", err)
			return err
		}
		encoded, err := s.encode()
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error encoding QC experimental settings during datastore save operation", GetCaller()),
				"error", err)
			return err
		}
		return b.Put([]byte(keyQCExperimentalSettings), encoded)
	})
}

func (s *QCExperimentalSettings) decode(data []byte) error {
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&s); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decoding QC experimental settings data", GetCaller()), "error", err)
		return err
	}
	return nil
}

func (s *QCExperimentalSettings) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(s); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error encoding QC experimental settings data", GetCaller()), "error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *QCExperimentalSettings) validate() error {
	if s == nil {
		return errors.New("QC Experimental setting info was not entered")
	}
	return nil
}

func (s *QCExperimentalSettings) resetBoolForZeroValues() {
	// values of zero for the following have absolutely no effect, so reset the state prior to saving
	if s.UseMaxFPSLimit && s.MaxFPSLimit == 0 {
		s.UseMaxFPSLimit = false
	}
	if s.UseMaxFPSLimitMinimized && s.MaxFPSLimitMinimized == 0 {
		s.UseMaxFPSLimitMinimized = false
	}
}
