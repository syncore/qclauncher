// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	uuid "github.com/satori/go.uuid"
)

const (
	actionGET  = "GET"
	actionPOST = "POST"
)

type localRequest interface {
	build(addr string) error
	action() string
	getParams() *requestParams
	expectedResponse() remoteResponseType
	needsContent() bool
}

type requestParams struct {
	header       map[string][]string
	endpointAddr string
}

type authRequest struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	Password  string `json:"password"`
	params    *requestParams
}

func (r *authRequest) build(addr string) error {
	creds, err := getUserCredentials()
	if err != nil {
		logger.Errorw("authRequest.build: error getting user credentials", "error", err)
		return err
	}
	r.Username = creds.Username
	r.Password = creds.Password
	r.SessionID = uuid.NewV4().String()
	header := &requestHeaderAuth{}
	err = header.build()
	if err != nil {
		logger.Errorw("authRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{
		header: header.values, endpointAddr: addr,
	}
	return nil
}

func (r *authRequest) getParams() *requestParams {
	return r.params
}

func (r *authRequest) expectedResponse() remoteResponseType {
	return rrAuth
}

func (r *authRequest) action() string {
	return actionPOST
}

func (r *authRequest) needsContent() bool {
	return true
}

type verifyRequest struct {
	SessionID string `json:"session_id"`
	params    *requestParams
}

func (r *verifyRequest) build(addr string) error {
	r.SessionID = uuid.NewV4().String()
	header := &requestHeaderVerify{}
	err := header.build()
	if err != nil {
		logger.Errorw("verifyRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *verifyRequest) getParams() *requestParams {
	return r.params
}

func (r *verifyRequest) expectedResponse() remoteResponseType {
	return rrAuth
}

func (r *verifyRequest) action() string {
	return actionPOST
}

func (r *verifyRequest) needsContent() bool {
	return true
}

type buildInfoRequest struct {
	params *requestParams
}

func (r *buildInfoRequest) build(addr string) error {
	header := &requestHeaderBuildInfo{}
	err := header.build()
	if err != nil {
		logger.Errorw("buildInfoRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *buildInfoRequest) getParams() *requestParams {
	return r.params
}

func (r *buildInfoRequest) expectedResponse() remoteResponseType {
	return rrBuildInfo
}

func (r *buildInfoRequest) action() string {
	return actionGET
}

func (r *buildInfoRequest) needsContent() bool {
	return false
}

type preSaveVerifyRequest struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	Password  string `json:"password"`
	params    *requestParams
}

func (r *preSaveVerifyRequest) build(addr string) error {
	// Login credentials are passed in the struct prior to allowing save
	r.SessionID = uuid.NewV4().String()
	header := &requestHeaderAuth{}
	err := header.build()
	if err != nil {
		logger.Errorw("preSaveVerifyRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{
		header: header.values, endpointAddr: addr,
	}
	return nil
}

func (r *preSaveVerifyRequest) getParams() *requestParams {
	return r.params
}

func (r *preSaveVerifyRequest) expectedResponse() remoteResponseType {
	return rrPreSave
}

func (r *preSaveVerifyRequest) action() string {
	return actionPOST
}

func (r *preSaveVerifyRequest) needsContent() bool {
	return true
}

type branchInfoRequest struct {
	params *requestParams
}

func (r *branchInfoRequest) build(addr string) error {
	header := &requestHeaderBranchInfo{}
	err := header.build()
	if err != nil {
		logger.Errorw("branchInfoRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *branchInfoRequest) getParams() *requestParams {
	return r.params
}

func (r *branchInfoRequest) expectedResponse() remoteResponseType {
	return rrBranchInfo
}

func (r *branchInfoRequest) action() string {
	return actionGET
}

func (r *branchInfoRequest) needsContent() bool {
	return false
}

type launchArgsRequest struct {
	params *requestParams
}

func (r *launchArgsRequest) build(addr string) error {
	header := &requestHeaderLaunchArgs{}
	err := header.build()
	if err != nil {
		logger.Errorw("launchArgsRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *launchArgsRequest) getParams() *requestParams {
	return r.params
}

func (r *launchArgsRequest) expectedResponse() remoteResponseType {
	return rrLaunchArgs
}

func (r *launchArgsRequest) action() string {
	return actionGET
}

func (r *launchArgsRequest) needsContent() bool {
	return false
}

type gameCodeRequest struct {
	params *requestParams
}

func (r *gameCodeRequest) build(addr string) error {
	header := &requestHeaderGameCode{}
	err := header.build()
	if err != nil {
		logger.Errorw("gameCodeRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *gameCodeRequest) getParams() *requestParams {
	return r.params
}

func (r *gameCodeRequest) expectedResponse() remoteResponseType {
	return rrGameCode
}

func (r *gameCodeRequest) action() string {
	return actionGET
}

func (r *gameCodeRequest) needsContent() bool {
	return false
}

type serverStatusRequest struct {
	params *requestParams
}

func (r *serverStatusRequest) build(addr string) error {
	header := &requestHeaderServerStatus{}
	err := header.build()
	if err != nil {
		logger.Errorw("serverStatusRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *serverStatusRequest) getParams() *requestParams {
	return r.params
}

func (r *serverStatusRequest) expectedResponse() remoteResponseType {
	return rrServerStatus
}

func (r *serverStatusRequest) action() string {
	return actionGET
}

func (r *serverStatusRequest) needsContent() bool {
	return false
}

type updateQCRequest struct {
	params *requestParams
}

func (r *updateQCRequest) build(addr string) error {
	header := &requestHeaderUpdateQC{}
	err := header.build()
	if err != nil {
		logger.Errorw("updateQCRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *updateQCRequest) getParams() *requestParams {
	return r.params
}

func (r *updateQCRequest) expectedResponse() remoteResponseType {
	return rrUpdateQC
}

func (r *updateQCRequest) action() string {
	return actionGET
}

func (r *updateQCRequest) needsContent() bool {
	return false
}

type updateLauncherRequest struct {
	params *requestParams
}

func (r *updateLauncherRequest) build(addr string) error {
	header := &requestHeaderUpdateLauncher{}
	err := header.build()
	if err != nil {
		logger.Errorw("updateLauncherRequest.build: error getting headers", "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *updateLauncherRequest) getParams() *requestParams {
	return r.params
}

func (r *updateLauncherRequest) expectedResponse() remoteResponseType {
	return rrUpdateLauncher
}

func (r *updateLauncherRequest) action() string {
	return actionGET
}

func (r *updateLauncherRequest) needsContent() bool {
	return false
}
