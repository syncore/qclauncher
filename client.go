// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	defArgs        = `--startup --set /Config/GAME_CONFIG/bethesdaGameCode "%GAMECODE%" --set /Config/GAME_CONFIG/bethesdaLoginEnabled 1 --set /Config/Bethesda/Language "%LANGUAGE%" --set /Config/GAME_CONFIG/bethesdaEndpointUrl "https://services.bethesda.net/agora_beam/"`
	qcEntitlmentID = 48329
	gameCodeTempl  = "%GAMECODE%"
	langTempl      = "%LANGUAGE%"
	defLang        = "en"
)

type launcherClient struct {
	Hc *http.Client
}

type authFailedError struct {
	emsg string
}

func (e *authFailedError) Error() string {
	return e.emsg
}

func Launch() error {
	lc := newLauncherClient(defTimeout)
	lc.checkServerStatus()
	err := lc.authenticate()
	if err != nil {
		return err
	}
	qcOpts, err := getQCOptions()
	if err != nil {
		logger.Errorw("Launch: getQCOptions error", "error", err, "data", qcOpts)
		return err
	}
	buildInfo, err := lc.getBuildInfo()
	if err != nil {
		logger.Errorw("Launch: getBuildInfo error", "error", err, "data", buildInfo)
		return err
	}
	logger.Debugw("Build info", "buildInfo", buildInfo)
	branchInfo, err := lc.getBranchInfo(buildInfo.Projects[0].ID, buildInfo.Branches[0].ID)
	if err != nil {
		logger.Errorw("Launch: getBranchInfo error", "error", err, "data", branchInfo)
		return err
	}
	logger.Debugw("Branch info", "branchInfo", branchInfo)
	launchArgs, err := lc.getLaunchArgs(buildInfo.Projects[0].ID)
	if err != nil {
		logger.Errorw("Launch: getLaunchArgs error", "error", err, "data", launchArgs)
		return err
	}
	exArgs := launchArgs.extractLaunchArgs(branchInfo.LaunchinfoList[0], qcOpts.QCLanguage)
	logger.Debugw("Launch args", "launchArgs", launchArgs)
	logger.Debugw("Extracted launch args", "exArgs", exArgs)
	gameCode, err := lc.getGameCode(buildInfo.Projects[0].ID)
	if err != nil {
		logger.Errorw("Launch: getGameCode error", "error", err, "data", gameCode)
		return err
	}
	logger.Debugw("Game code", "gameCode.GameCode", gameCode.Gamecode)
	finalArgs := strings.Replace(exArgs, gameCodeTempl, gameCode.Gamecode, -1)
	return runQC(qcOpts.QCFilePath, finalArgs)
}

func newLauncherClient(timeout int) *launcherClient {
	return &launcherClient{
		Hc: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (lc *launcherClient) getBuildInfo() (*BuildInfoResponse, error) {
	req := &buildInfoRequest{}
	err := req.build(getBuildInfoEndpoint())
	if err != nil {
		logger.Errorw("getBuildInfo: error building build info request", "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if buildInfoResponse, ok := res.(BuildInfoResponse); ok {
		return &buildInfoResponse, nil
	} else if err != nil {
		logger.Errorw("getBuildInfo: error parsing build info response", "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw("getBuildInfo: got an unexpected build info response", "error", err, "data", res)
		return nil, formatUnexpectedResponse("receiving build information")
	}
}

func (lc *launcherClient) getBranchInfo(projectID, branchID int) (*BranchInfoResponse, error) {
	req := &branchInfoRequest{}
	err := req.build(getBranchInfoEndpoint(projectID, branchID))
	if err != nil {
		logger.Errorw("getBranchInfo: error building branch info request", "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if branchInfoResponse, ok := res.(BranchInfoResponse); ok {
		return &branchInfoResponse, nil
	} else if err != nil {
		logger.Errorw("getBranchInfo: error parsing branch info response", "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw("getBranchInfo: got an unexpected branch info response", "error", err, "data", res)
		return nil, formatUnexpectedResponse("receiving branch information")
	}
}

func (lc *launcherClient) getLaunchArgs(projectID int) (*LaunchArgsResponse, error) {
	req := &launchArgsRequest{}
	err := req.build(getLaunchArgsEndpoint(projectID))
	if err != nil {
		logger.Errorw("getLaunchArgs: error building launch args request", "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if launchArgsResponse, ok := res.(LaunchArgsResponse); ok {
		return &launchArgsResponse, nil
	} else if err != nil {
		logger.Errorw("getLaunchArgs: error parsing launch args response", "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw("getLaunchArgs: got an unexpected launch args response", "error", err, "data", res)
		return nil, formatUnexpectedResponse("retrieving launch args")
	}
}

func (lc *launcherClient) getGameCode(projectID int) (*GameCodeResponse, error) {
	req := &gameCodeRequest{}
	err := req.build(getGameCodeEndpoint(projectID))
	if err != nil {
		logger.Errorw("getGameCode: error building game code request", "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if gameCodeResponse, ok := res.(GameCodeResponse); ok {
		return &gameCodeResponse, nil
	} else if err != nil {
		logger.Errorw("getGameCode: error parsing game code response", "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw("getGameCode: got an unexpected game code response", "error", err, "data", res)
		return nil, formatUnexpectedResponse("receiving game code")
	}
}

func (lc *launcherClient) getServerStatus() (*ServerStatusResponse, error) {
	req := &serverStatusRequest{}
	err := req.build(getServerStatusEndpoint())
	if err != nil {
		logger.Errorw("getServerStatus: error building server status request", "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if serverStatusResponse, ok := res.(ServerStatusResponse); ok {
		return &serverStatusResponse, nil
	} else if err != nil {
		logger.Errorw("getServerStatus: error parsing server status response", "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw("getServerStatus: got an unexpected server status response", "error", err, "data", res)
		return nil, formatUnexpectedResponse("checking QC server status")
	}
}

func (lc *launcherClient) getQCUpdateInfo() (*UpdateQCResponse, error) {
	req := &updateQCRequest{}
	err := req.build(updateQCEndpoint)
	if err != nil {
		logger.Errorw("getQCUpdateInfo: error building QC update request", "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if updateQCResponse, ok := res.(UpdateQCResponse); ok {
		return &updateQCResponse, nil
	} else if err != nil {
		logger.Errorw("getQCUpdateInfo: error parsing QC update response", "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw("getQCUpdateInfo: got an unexpected QC update response", "error", err, "data", res)
		return nil, formatUnexpectedResponse("checking for QC updates")
	}
}

func (lc *launcherClient) getLauncherUpdateInfo() (*UpdateLauncherResponse, error) {
	req := &updateLauncherRequest{}
	err := req.build(updateLauncherEndpoint)
	if err != nil {
		logger.Errorw("getLauncherUpdateInfo: error building launcher update request", "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if updateLauncherResponse, ok := res.(UpdateLauncherResponse); ok {
		return &updateLauncherResponse, nil
	} else if err != nil {
		logger.Errorw("getLauncherUpdateInfo: error parsing launcher update response", "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw("getLauncherUpdateInfo: got an unexpected launcher update response", "error", err, "data", res)
		return nil, formatUnexpectedResponse("checking for launcher updates")
	}
}

func (lc *launcherClient) checkServerStatus() {
	status, err := lc.getServerStatus()
	if err != nil {
		logger.Errorw("checkServerStatus: error checking server status", "error", err)
		return
	}
	if strings.EqualFold(status.Platform.Response.Quake, "DOWN") {
		ShowWarningMsg("Warning", "The QC servers are currently offline. Launch will continue but you will be unable to play.")
	}
}

func (lc *launcherClient) authenticate() error {
	var err error
	creds, cerr := getUserCredentials()
	if cerr != nil {
		logger.Errorw("authenticate: error occurred when getting user credentials", "error", err, "data", creds)
		return cerr
	}
	// Verify
	if creds.Token != "" {
		vreq := &verifyRequest{}
		err = vreq.build(getVerifyEndpoint())
		if err != nil {
			logger.Errorw("authenticate: error building verify request", "error", err)
			return err
		}
		vres, err := lc.send(vreq)
		if _, ok := vres.(AuthResponse); !ok {
			logger.Errorw("authenticate: unexpected verify response type", "error", err, "data", vres)
			return errors.New("Received an unexpected response during authentication")
		}
		if err != nil {
			logger.Errorw("authenticate: error receiving verify response", "error", err)
			return err
		}
	} else {
		// Auth
		areq := &authRequest{}
		err := areq.build(getAuthEndpoint())
		if err != nil {
			logger.Errorw("authenticate: error building auth request", "error", err)
			return err
		}
		ares, err := lc.send(areq)
		if _, ok := ares.(AuthResponse); !ok {
			logger.Errorw("authenticate: unexpected auth response type", "error", err, "data", ares)
			return formatUnexpectedResponse("performing authentication")
		}
		if err != nil {
			logger.Errorw("authenticate: error receiving auth response", "error", err)
			return err
		}
	}
	return nil
}

func (lc *launcherClient) verifyCredentials(user, password string) error {
	areq := &preSaveVerifyRequest{Username: user, Password: password}
	err := areq.build(getAuthEndpoint())
	if err != nil {
		logger.Errorw("verifyCredentials: error building pre-save auth request", "error", err)
		return err
	}
	ares, err := lc.send(areq)
	if _, ok := err.(*authFailedError); ok {
		emsg := "Login failed. This needs to be the same as your Bethesda Launcher login information. Please try again."
		logger.Errorf("verifyCredentials: %s", emsg)
		return errors.New(emsg)
	}
	if _, ok := ares.(AuthResponse); !ok {
		logger.Errorw("verifyCredentials: unexpected auth response type", "error", err, "data", ares)
		return formatUnexpectedResponse("verifying login information")
	}
	if err != nil {
		logger.Errorw("verifyCredentials: error receiving auth response", "error", err)
		return err
	}
	return nil
}

func (lc *launcherClient) send(req localRequest) (interface{}, error) {
	p := req.getParams()
	var br io.Reader
	if req.needsContent() {
		j, err := json.Marshal(req)
		br = strings.NewReader(string(j))
		if err != nil {
			logger.Errorw("send: error marshaling JSON", "error", err, "data", req)
			return nil, err
		}
	}
	hr, err := http.NewRequest(req.action(), p.endpointAddr, br)
	if err != nil {
		logger.Errorw("send: error creating request", "error", err, "data", hr)
		return nil, err
	}
	hr.Header = p.header
	res, err := lc.Hc.Do(hr)
	if err != nil {
		logger.Errorw("send: error sending request", "error", err, "data", hr)
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorw("send: error reading response body", "error", err, "data", string(b))
		return nil, err
	}
	logger.Debugw("send: response body", "expectedResponse", req.expectedResponse(), "body", string(b))
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			logger.Errorf("send: got unauthorized response when accessing resource (%s) requiring authentication", p.endpointAddr)
			return nil, &authFailedError{emsg: "User authentication failed"}
		}
		logger.Errorw("send: got non-OK status code", "error", err, "statusCode", res.StatusCode)
		return nil, fmt.Errorf("send: Non-OK status code received: %d", res.StatusCode)
	}
	rd := &remoteResponseData{ResponseType: req.expectedResponse()}
	err = json.Unmarshal(b, &rd.Data)
	if err != nil {
		logger.Errorw("send: error unmarshaling resposne body into model", "error", err, "data", string(b))
		return nil, err
	}
	response, err := parseRemoteResponseData(rd)
	if err != nil {
		logger.Error("send: error parsing remote response data", "error", err, "data", string(rd.Data))
		return nil, err
	}
	return response, nil
}

func (r *LaunchArgsResponse) extractLaunchArgs(liKey int, language string) string {
	fallback := strings.Replace(defArgs, langTempl, defLang, -1)
	if r == nil {
		logger.Info("extractLaunchArgs: launch args was nil, using fallback arguments")
		return fallback
	}
	if (LaunchInfo{}) == r.LaunchinfoSet {
		logger.Info("extractLaunchArgs: launch info set had default struct values, using fallback arguments")
		return fallback
	}
	if (LaunchInfoItem{}) == r.LaunchinfoSet.Default {
		logger.Info("extractLaunchArgs: default launch info item had default struct values, using fallback arguments")
		return fallback
	}
	likstr := strconv.Itoa(liKey)
	val := reflect.ValueOf(r.LaunchinfoSet)
	found := false
	for i := 0; i < val.Type().NumField(); i++ {
		v := val.Type().Field(i).Tag.Get("json")
		if v != likstr {
			continue
		}
		found = true
	}
	if found {
		v := r.LaunchinfoSet.Default.LaunchArgs
		v = strings.Replace(v, "\\", "", -1)
		v = strings.Replace(v, langTempl, language, -1)
		return v
	}
	logger.Infof("extractLaunchArgs: launch info key (%d) passed in had no match, using fallback arguments", liKey)
	return fallback
}

func formatUnexpectedResponse(event string) error {
	// Formatting for generai errors that occur during launcher client actions (displayed in msg box)
	return fmt.Errorf("Received an unexpected response when %s.", event)
}

func runQC(qcPath, qcArgs string) error {
	qc := exec.Command(qcPath)
	qc.Dir = filepath.Dir(qcPath)
	allArgs := []string{qcArgs}
	if ConfAppendCustomArgs != "" {
		allArgs = append(allArgs, ConfAppendCustomArgs)
	}
	if ConfMaxFPS != 0 {
		allArgs = append(allArgs, fmt.Sprintf("--set /Config/CONFIG/maxFpsValue %d", ConfMaxFPS))
	}
	a := strings.Join(allArgs, " ")
	logger.Debugf("Final arguments: %s", a)
	// Handle arg quote-escaping manually (see golang issue #15566)
	qc.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CmdLine:       fmt.Sprintf(` %s`, a),
		CreationFlags: 0,
	}
	logger.Debug("Launching....")
	if err := qc.Start(); err != nil {
		logger.Errorw("runQC: error starting QC", "error", err)
		return err
	}
	return nil
}
