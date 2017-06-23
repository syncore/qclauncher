// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gtank/cryptopasta"
	ps "github.com/keybase/go-ps"
)

func IsProcessRunning(processNames ...string) (bool, string, error) {
	runningProcesses, err := ps.Processes()
	if err != nil {
		return false, "", fmt.Errorf("Error enumerating processes: %s", err)
	}
	var results []string
	found := false
	for _, rp := range runningProcesses {
		for _, pn := range processNames {
			if !strings.EqualFold(pn, rp.Executable()) {
				continue
			}
			found = true
			results = append(results, pn)
			break
		}
	}
	return found, strings.Join(results, ", "), nil
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getExecutingPath() string {
	executingPath, err := os.Executable()
	if err != nil {
		panic("Unable to determine execution path")
	}
	return filepath.Dir(executingPath)
}

func dirExists(dir string) bool {
	f, err := os.Stat(dir)
	return err == nil && f.IsDir()
}

func createEmptyFile(fullpath string, allowOverwrite bool) error {
	if FileExists(fullpath) && !allowOverwrite {
		return fmt.Errorf("File %s exists and allowOverwrite is false", fullpath)
	}
	f, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("Couldn't create empty file '%s': %s\n", fullpath, err)
	}
	defer f.Close()
	return nil
}

func randStr(strlen int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, strlen)
	const chars = "AbcDEfGHiJKlmnOpQRsTuVWxyZ012345679/+"
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func genFp() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	h := sha1.New()
	h.Write([]byte(randStr(r.Intn(64)*4 + 64)))
	s := fmt.Sprintf("%x", h.Sum(nil))
	return strings.ToUpper(s)
}

func genKey() *[]byte {
	k := cryptopasta.NewEncryptionKey()
	ek := make([]byte, 32)
	for i, b := range *k {
		ek[i] = b
	}
	return &ek
}

func encrypt(unencrypted string, key *[]byte) (encrypted string, err error) {
	keyLength := len(*key)
	if keyLength != 32 {
		return "", fmt.Errorf("Key has invalid size of %d bytes; expected size of 32 bytes", keyLength)
	}
	var k [32]byte
	for i, b := range *key {
		k[i] = b
	}
	enc, err := cryptopasta.Encrypt([]byte(unencrypted), &k)
	return string(enc), err
}

func decrypt(encrypted string, key *[]byte) (decrypted string, err error) {
	keyLength := len(*key)
	if keyLength != 32 {
		return "", fmt.Errorf("Key has invalid size of %d bytes; expected size of 32 bytes", keyLength)
	}
	var k [32]byte
	for i, b := range *key {
		k[i] = b
	}
	dec, err := cryptopasta.Decrypt([]byte(encrypted), &k)
	return string(dec), err
}
