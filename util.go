// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gtank/cryptopasta"
	ps "github.com/keybase/go-ps"
)

// Single provides a mechanism to ensure that only one instance of a program is running
// https://github.com/WeltN24/single
type Single struct {
	name   string
	file   *os.File
	Locked bool
}

func IsProcessRunning(processNames ...string) (bool, string, int, []string, map[string]int, error) {
	var pfilenames []string
	namepid := make(map[string]int)
	runningProcesses, err := ps.Processes()
	if err != nil {
		return false, "", 0, pfilenames, namepid, fmt.Errorf("Error enumerating processes: %s", err)
	}
	found := false
	for _, rp := range runningProcesses {
		for _, pn := range processNames {
			if !strings.EqualFold(pn, rp.Executable()) {
				continue
			}
			found = true
			pfilenames = append(pfilenames, pn)
			if _, ok := namepid[pn]; !ok {
				namepid[pn] = rp.Pid()
			}
			break
		}
	}
	return found, strings.Join(pfilenames, ", "), len(pfilenames), pfilenames, namepid, nil
}

func GetCaller() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	fnName := f.Name()
	i := strings.LastIndex(fnName, ".")
	if i != -1 {
		return fnName[i+1:]
	}
	return fnName
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func DeleteFile(filepath string) error {
	if !FileExists(filepath) {
		return os.ErrNotExist
	}
	err := os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("Error deleting file: %s", err)
	}
	return nil
}

func NewSingle(name string) *Single {
	return &Single{name: name, Locked: false}
}

func (s *Single) Wait() {
	locked := true
	for locked {
		time.Sleep(time.Millisecond)
		err := s.Lock()
		locked = err != nil
		if err == nil {
			_ = s.Unlock()
		}
	}
}

func (s *Single) Filename(useFullPath bool) string {
	if useFullPath {
		return filepath.Join(getExecutingPath(), s.name)
	}
	return fmt.Sprintf(s.name)
}

func (s *Single) Lock() error {
	if err := os.Remove(s.Filename(true)); err != nil && !os.IsNotExist(err) {
		return &alreadyRunningError{emsg: fmt.Sprintf("QCLauncher v%.2f is already running.", version)}
	}
	file, err := os.OpenFile(s.Filename(true), os.O_EXCL|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	s.file = file
	return nil
}

func (s *Single) Unlock() error {
	if err := s.file.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.Filename(true)); err != nil {
		return err
	}
	return nil
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
