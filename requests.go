// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"

	"github.com/google/uuid"
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
	cfg, err := GetConfiguration()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting configuration for credential lookup when building auth request",
			GetCaller()), "error", err)
		return err
	}
	r.Username = cfg.Core.Username
	r.Password = cfg.Core.Password
	r.SessionID = uuid.New().String()
	header := &requestHeaderAuth{}
	err = header.build()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting auth request headers", GetCaller()), "error", err)
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
	r.SessionID = uuid.New().String()
	header := &requestHeaderVerify{}
	err := header.build()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting verify request headers", GetCaller()), "error", err)
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

type entitlementInfoRequest struct {
	EntitlementIDs []int `json:"entitlement_ids"`
	params         *requestParams
}

func (r *entitlementInfoRequest) build(addr string) error {
	r.EntitlementIDs = []int{0}
	header := &requestHeaderEntitlementInfo{}
	err := header.build()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting entitlement info request headers", GetCaller()), "error", err)
		return err
	}
	r.params = &requestParams{header: header.values, endpointAddr: addr}
	return nil
}

func (r *entitlementInfoRequest) getParams() *requestParams {
	return r.params
}

func (r *entitlementInfoRequest) expectedResponse() remoteResponseType {
	return rrEntitlementInfo
}

func (r *entitlementInfoRequest) action() string {
	return actionPOST
}

func (r *entitlementInfoRequest) needsContent() bool {
	return true
}

type preSaveVerifyRequest struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	Password  string `json:"password"`
	params    *requestParams
}

func (r *preSaveVerifyRequest) build(addr string) error {
	// Login credentials are passed in the struct prior to allowing save
	r.SessionID = uuid.New().String()
	header := &requestHeaderAuth{}
	err := header.build()
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error getting pre-save auth verification request headers",
			GetCaller()), "error", err)
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
		logger.Errorw(fmt.Sprintf("%s: error getting branch info request headers", GetCaller()), "error", err)
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
		logger.Errorw(fmt.Sprintf("%s: error getting launch args request headers", GetCaller()), "error", err)
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
		logger.Errorw(fmt.Sprintf("%s: error getting game code request headers", GetCaller()), "error", err)
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
		logger.Errorw(fmt.Sprintf("%s: error getting server status request headers", GetCaller()), "error", err)
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
		logger.Errorw(fmt.Sprintf("%s: error getting QC update request headers", GetCaller()), "error", err)
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
		logger.Errorw(fmt.Sprintf("%s: error getting launcher update request headers", GetCaller()), "error", err)
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
