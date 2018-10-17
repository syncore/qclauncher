package qclauncher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Microsoft/go-winio"
	"github.com/syncore/qclauncher/resources"
)

const pipePrefix = "\\\\.\\pipe\\"
const pipeName = "blffqcl"
const fpAttempts = 4

var extractArgs = []string{fmt.Sprintf("-p=%s", pipeName), fmt.Sprintf("-r=%d", fpAttempts)}

//var extractArgs = []string{fmt.Sprintf("-p=%s", pipeName), "-t", "-f=testfp.json", fmt.Sprintf("-r=%d", fpAttempts)}

type fpChanResult struct {
	fp  []byte
	err error
}

type bnlFingerprint struct {
	FP *string `json:"fp"`
}

func getBNLFingerprint() (string, error) {
	fpChan := make(chan fpChanResult)
	go readFromPipe(fpChan)
	if err := extractFp(); err != nil {
		return "", err
	}
	result := <-fpChan
	rfp, rerr := result.fp, result.err
	if rerr != nil {
		return "", rerr
	}
	bnl := &bnlFingerprint{}
	if err := json.Unmarshal(rfp, bnl); err != nil {
		return "", err
	}
	if bnl.FP == nil {
		err := "Received null FP value indicating that FP retrieval process could not find FP"
		logger.Errorf("%s: %s", GetCaller(), err)
		return "", fmt.Errorf("%s", err)
	}
	tmpFp = *bnl.FP // set until it can be written to storage
	return *bnl.FP, nil
}

func readFromPipe(c chan fpChanResult) {
	lp, err := winio.ListenPipe(pipePrefix+pipeName, &winio.PipeConfig{
		MessageMode:      true,
		InputBufferSize:  64,
		OutputBufferSize: 64,
	})
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error creating named pipe", GetCaller()), "error", err)
		c <- fpChanResult{fp: []byte{}, err: err}
		return
	}
	defer lp.Close()
	pipe, err := lp.Accept()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error accepting named pipe connection", GetCaller()), "error", err)
		c <- fpChanResult{fp: []byte{}, err: err}
		return
	}
	defer pipe.Close()
	result, err := ioutil.ReadAll(bufio.NewReader(pipe))
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error reading from named pipe", GetCaller()), "error", err)
		c <- fpChanResult{fp: []byte{}, err: err}
		return
	}
	logger.Debugf("readFromPipe: Read %d bytes. FP (raw): %+x output: %s", len(result), result, string(result))
	logger.Debugw("readFromPipe", "result", string(result))
	c <- fpChanResult{fp: result, err: nil}
}

func extractFp() error {
	a, err := resources.Asset("../../resources/bin/blff.exe") // https://github.com/syncore/blff
	//a, err := resources.Asset("../../resources/bin/blff_console.exe")
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error reading FP extraction tool asset", GetCaller()), "error", err)
		return err
	}
	outname := filepath.Join(getExecutingPath(), "ExtractBNLauncherFP.exe")
	err = ioutil.WriteFile(outname, a, 0644)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error extracting FP extraction tool", GetCaller()), "error", err)
		return err
	}
	blff := exec.Command(outname, extractArgs...)
	logger.Debugf("extractFp: Executing blff (%s) and awaiting completion...", outname)
	err = blff.Run()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error occurred while running FP extraction tool", GetCaller()), "error", err)
		return err
	}
	err = os.Remove(outname)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error occurred while cleaning up FP extraction tool", GetCaller()), "error", err)
		return err
	}
	return nil
}