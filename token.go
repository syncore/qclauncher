// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"bytes"
	"encoding/gob"
	"fmt"

	bolt "github.com/coreos/bbolt"
)

type TokenAuth struct {
	Token string
}

type TokenKey struct {
	Key []byte
}

func (t *TokenAuth) get(ls *LauncherStore) error {
	ls.checkDataFile(false)
	if err := ls.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSettings))
		decerr := t.decode(b.Get([]byte(keyTokenAuth)), b.Get([]byte(keyTokenKey)))
		if decerr != nil {
			logger.Errorw(fmt.Sprintf("%s: error decoding auth token from datastore during get operation", GetCaller()),
				"error", decerr)
			return decerr
		}
		return nil
	}); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting auth token from datastore", GetCaller()), "error", err)
		return err
	}
	return nil
}

func (t *TokenAuth) save(ls *LauncherStore) error {
	return ls.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketSettings))
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating settings bucket in datastore during save operation",
				GetCaller()), "error", err)
			return err
		}
		key := b.Get([]byte(keyTokenKey))
		encoded, err := t.encode(key)
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error encoding auth token during datastore save operation", GetCaller()),
				"error", err)
			return err
		}
		return b.Put([]byte(keyTokenAuth), encoded)
	})
}

func (t *TokenAuth) decode(data, key []byte) error {
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&t); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decoding auth token data", GetCaller()), "error", err)
		return err
	}
	qcDecToken, err := decrypt(t.Token, &key)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decrypting auth token credential", GetCaller()), "error", err)
		return err
	}
	t.Token = qcDecToken
	return nil
}

func (t *TokenAuth) encode(existingKey []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	encToken, err := encrypt(t.Token, &existingKey)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error encrypting auth token credential", GetCaller()), "error", err)
		return nil, err
	}
	t.Token = encToken
	if err = gob.NewEncoder(buf).Encode(t); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error encoding auth token data", GetCaller()), "error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TokenKey) get(ls *LauncherStore) error {
	if err := ls.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSettings))
		decerr := t.decode(b.Get([]byte(keyTokenKey)))
		if decerr != nil {
			logger.Errorw(fmt.Sprintf("%s: error decoding credential key from datastore during get operation", GetCaller()),
				"error", decerr)
			return decerr
		}
		return nil
	}); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting credential key from datastore", GetCaller()), "error", err)
		return err
	}
	return nil
}

func (t *TokenKey) save(ls *LauncherStore) error {
	return ls.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketSettings))
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error creating settings bucket in datastore during save operation",
				GetCaller()), "error", err)
			return err
		}
		encoded, err := t.encode()
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error encoding credential key during datastore save operation", GetCaller()),
				"error", err)
			return err
		}
		return b.Put([]byte(keyTokenKey), encoded)
	})
}

func (t *TokenKey) decode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(&t); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error decoding credential key data", GetCaller()), "error", err)
		return err
	}
	return nil
}

func (t *TokenKey) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(t)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error encoding credential key data", GetCaller()), "error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func updateAuthToken(isPreSaveVerification bool, token string) error {
	if isPreSaveVerification {
		// Data file won't exist on first-run credential verification; which is the entry point into
		// the data store, so save token & key in temp vars so they will be applied when the core
		// settings are saved.
		tmpToken = token
		tmpKey = genKey()
		return nil
	}
	return Save(&TokenAuth{Token: token})
}

func clearAuthToken() error {
	return Save(&TokenAuth{Token: ""})
}
