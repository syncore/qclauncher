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
	"strings"
	"time"
)

type launcherClient struct {
	*http.Client
}

func newLauncherClient(timeout int) *launcherClient {
	return &launcherClient{
		&http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

func (lc *launcherClient) getEntitlementInfo() (*EntitlementInfoResponse, error) {
	req := &entitlementInfoRequest{}
	if err := req.build(getEntitlementInfoEndpoint()); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building entitlement info request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if entitlementInfoResponse, ok := res.(EntitlementInfoResponse); ok {
		return &entitlementInfoResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing entitlement info response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected entitlement info response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("receiving entitlement information")
	}
}

func (lc *launcherClient) getBuildInfo() (*BuildInfoResponse, error) {
	req := &buildInfoRequest{}
	if err := req.build(getBuildInfoEndpoint()); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building build info request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if buildInfoResponse, ok := res.(BuildInfoResponse); ok {
		return &buildInfoResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing build info response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected build info response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("receiving build information")
	}
}

func (lc *launcherClient) checkEntitlement() (*EntitlementCheckAPIResponse, error) {
	req := &entitlementCheckAPIRequest{}
	if err := req.build(entitlementCheckAPIEndpoint); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building entitlement check API request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if entitlementCheckAPIResponse, ok := res.(EntitlementCheckAPIResponse); ok {
		return &entitlementCheckAPIResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing entitlement check API response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected entitlement check API response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("receiving entitlement check API information")
	}
}

func (lc *launcherClient) getBranchInfo(projectID, branchID int) (*BranchInfoResponse, error) {
	req := &branchInfoRequest{}
	if err := req.build(getBranchInfoEndpoint(projectID, branchID)); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building branch info request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if branchInfoResponse, ok := res.(BranchInfoResponse); ok {
		return &branchInfoResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing branch info response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected branch info response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("receiving branch information")
	}
}

func (lc *launcherClient) getLaunchArgs(projectID int) (*LaunchArgsResponse, error) {
	req := &launchArgsRequest{}
	if err := req.build(getLaunchArgsEndpoint(projectID)); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building launch args request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if launchArgsResponse, ok := res.(LaunchArgsResponse); ok {
		return &launchArgsResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing launch args response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected launch args response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("retrieving launch args")
	}
}

func (lc *launcherClient) getGameCode(projectID int) (*GameCodeResponse, error) {
	req := &gameCodeRequest{}
	if err := req.build(getGameCodeEndpoint(projectID)); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building game code request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if gameCodeResponse, ok := res.(GameCodeResponse); ok {
		return &gameCodeResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing game code response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected game code response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("receiving game code")
	}
}

func (lc *launcherClient) getServerStatus() (*ServerStatusResponse, error) {
	req := &serverStatusRequest{}
	if err := req.build(getServerStatusEndpoint()); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building server status request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if serverStatusResponse, ok := res.(ServerStatusResponse); ok {
		return &serverStatusResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing server status response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected server status response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("checking QC server status")
	}
}

func (lc *launcherClient) getQCUpdateInfo() (*UpdateQCResponse, error) {
	req := &updateQCRequest{}
	if err := req.build(updateQCEndpoint); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building QC update request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if updateQCResponse, ok := res.(UpdateQCResponse); ok {
		return &updateQCResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing QC update response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected QC update response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("checking for QC updates")
	}
}

func (lc *launcherClient) getLauncherUpdateInfo() (*UpdateLauncherResponse, error) {
	req := &updateLauncherRequest{}
	if err := req.build(updateLauncherEndpoint); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building launcher update request", GetCaller()), "error", err, "data", req)
		return nil, err
	}
	res, err := lc.send(req)
	if updateLauncherResponse, ok := res.(UpdateLauncherResponse); ok {
		return &updateLauncherResponse, nil
	} else if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing launcher response", GetCaller()), "error", err, "data", res)
		return nil, err
	} else {
		logger.Errorw(fmt.Sprintf("%s: got an unexpected launcher update response", GetCaller()), "error", err, "data", res)
		return nil, formatUnexpectedResponse("checking for launcher updates")
	}
}

func (lc *launcherClient) checkServerStatus() {
	status, err := lc.getServerStatus()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error checking server status", GetCaller()), "error", err)
		return
	}
	if strings.EqualFold(status.Platform.Response.Quake, "DOWN") {
		ShowWarningMsg("Warning", "The QC servers are currently offline. Launch will continue but you will be unable to play.", nil)
	}
}

func (lc *launcherClient) authenticate(cfg *Configuration) error {
	// Verify
	if cfg.Auth.Token != "" {
		vreq := &verifyRequest{}
		if err := vreq.build(getVerifyEndpoint()); err != nil {
			logger.Errorw(fmt.Sprintf("%s: error building verify request", GetCaller()), "error", err)
			return err
		}
		vres, err := lc.send(vreq)
		if IsErrAuthFailed(err) {
			logger.Error(fmt.Sprintf("%s: stale authentication token. clearing token for next attempt.", GetCaller()))
			if cerr := clearAuthToken(); cerr != nil {
				logger.Errorw(fmt.Sprintf("%s: unable to clear stale authentication token, data file will need to be reset", GetCaller()),
					"error", cerr)
				DeleteConfiguration(true)
				return cerr
			}
			return &authFailedError{emsg: "Bethesda server authentication error. Please try launching again."}
		}
		if _, ok := vres.(AuthResponse); !ok {
			logger.Errorw(fmt.Sprintf("%s: unexpected verify response type", GetCaller()), "error", err, "data", vres)
			return errors.New("Received an unexpected response during authentication")
		}
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error receiving verify response", GetCaller()), "error", err, "data", vres)
			return err
		}
	} else {
		// Auth
		areq := &authRequest{}
		if err := areq.build(getAuthEndpoint()); err != nil {
			logger.Errorw(fmt.Sprintf("%s: error building auth request", GetCaller()), "error", err)
			return err
		}
		ares, err := lc.send(areq)
		if _, ok := ares.(AuthResponse); !ok {
			logger.Errorw(fmt.Sprintf("%s: unexpected auth response type", GetCaller()), "error", err, "data", ares)
			return formatUnexpectedResponse("performing authentication")
		}
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error receiving auth response", GetCaller()), "error", err, "data", ares)
			return err
		}
	}
	return nil
}

func (lc *launcherClient) verifyCredentials(user, password string) error {
	areq := &preSaveVerifyRequest{Username: user, Password: password}
	if err := areq.build(getAuthEndpoint()); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error building pre-save auth request", GetCaller()), "error", err)
		return err
	}
	ares, err := lc.send(areq)
	if IsErrAuthFailed(err) {
		emsg := "Login failed. This needs to be the same as your Bethesda Launcher login information. Please try again."
		logger.Error(fmt.Sprintf("%s: %s", GetCaller(), emsg))
		return errors.New(emsg)
	}
	if _, ok := ares.(AuthResponse); !ok {
		logger.Errorw(fmt.Sprintf("%s: unexpected auth response type", GetCaller()), "error", err, "data", ares)
		return formatUnexpectedResponse("verifying login information")
	}
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error receiving auth response", GetCaller()), "error", err)
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
			logger.Errorw(fmt.Sprintf("%s: error marshaling JSON", GetCaller()), "error", err, "data", req)
			return nil, err
		}
	}
	hr, err := http.NewRequest(req.action(), p.endpointAddr, br)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error creating request", GetCaller()), "error", err, "data", hr)
		return nil, err
	}
	hr.Header = p.header
	res, err := lc.Do(hr)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error sending request", GetCaller()), "error", err, "data", hr)
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error reading response body", GetCaller()), "error", err, "data", string(b))
		return nil, err
	}
	logger.Debugw("send: response body", "expectedResponse", req.expectedResponse(), "body", string(b))
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			logger.Error(fmt.Sprintf("%s: got unauthorized response when accessing resource (%s) requiring authentication",
				GetCaller(), p.endpointAddr))
			return nil, &authFailedError{emsg: "User authentication failed"}
		}
		logger.Errorw(fmt.Sprintf("%s: got non-OK status code", GetCaller()), "error", err, "statusCode", res.StatusCode)
		return nil, fmt.Errorf("send: Non-OK status code received: %d", res.StatusCode)
	}
	rd := &remoteResponseData{ResponseType: req.expectedResponse()}
	err = json.Unmarshal(b, &rd.Data)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error unmarshaling resposne body into model", GetCaller()), "error", err, "data", string(b))
		return nil, err
	}
	response, err := parseRemoteResponseData(rd)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error parsing remote response data", GetCaller()), "error", err, "data", string(rd.Data))
		return nil, err
	}
	return response, nil
}
