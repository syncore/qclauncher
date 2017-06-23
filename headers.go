// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"
)

const (
	hkeyAccept                 = "Accept"
	hkeyAcceptEncoding         = "Accept-Encoding"
	hkeyAuthorization          = "Authorization"
	hkeyContentType            = "Content-Type"
	hkeyHost                   = "Host"
	hkeyUa                     = "User-Agent"
	hkeyXCdpApp                = "x-cdp-app"
	hkeyXCdpAppVer             = "x-cdp-app-ver"
	hkeyXCdpLibVer             = "x-cdp-lib-ver"
	hkeyXCdpPlatform           = "x-cdp-platform"
	hkeyXSrcFp                 = "x-src-fp"
	hvalServicesHost           = "services.bethesda.net"
	hvalBuildHost              = "buildinfo.cdp.bethesda.net"
	hvalAcceptAll              = "*/*"
	hvalAcceptEncodingIdentity = "identity"
	hvalAcceptEncodingAppJSON  = "application/json"
	hvalUserAgent              = "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36"
	hvalXCdpApp                = "Bethesda Launcher"
	hvalXCdpPlatform           = "Win/32"
)

var (
	genericBuildBaseHeaders = map[string][]string{
		hkeyHost:           []string{hvalBuildHost},
		hkeyAccept:         []string{hvalAcceptAll},
		hkeyUa:             []string{hvalUserAgent},
		hkeyAcceptEncoding: []string{hvalAcceptEncodingIdentity}}
	xcdpHeaders = map[string]string{
		hkeyXCdpApp:      hvalXCdpApp,
		hkeyXCdpPlatform: hvalXCdpPlatform,
	}
	genericExtraHeaders = localRequestExtraHeaders{xcdp: false, auth: false}
)

type localRequestHeader interface {
	build() error
	getExtra() localRequestExtraHeaders
	getBase() (headers map[string][]string)
}

type headerMapping struct {
	values map[string][]string
}

type localRequestExtraHeaders struct {
	xcdp bool
	auth bool
}

type requestHeaderAuth struct{ headerMapping }
type requestHeaderVerify struct{ headerMapping }
type requestHeaderBuildInfo struct{ headerMapping }
type requestHeaderBranchInfo struct{ headerMapping }
type requestHeaderLaunchArgs struct{ headerMapping }
type requestHeaderGameCode struct{ headerMapping }
type requestHeaderServerStatus struct{ headerMapping }
type requestHeaderUpdateQC struct{ headerMapping }
type requestHeaderUpdateLauncher struct{ headerMapping }

func (h *requestHeaderAuth) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderAuth.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h requestHeaderAuth) getExtra() localRequestExtraHeaders {
	return localRequestExtraHeaders{xcdp: true, auth: false}
}

func (h requestHeaderAuth) getBase() (headers map[string][]string) {
	return map[string][]string{
		hkeyHost:           []string{hvalServicesHost},
		hkeyAcceptEncoding: []string{hvalAcceptEncodingIdentity},
		hkeyAccept:         []string{hvalAcceptEncodingAppJSON},
		hkeyContentType:    []string{hvalAcceptEncodingAppJSON},
		hkeyUa:             []string{hvalUserAgent},
	}
}

func (h *requestHeaderVerify) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderVerify.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h requestHeaderVerify) getExtra() localRequestExtraHeaders {
	return localRequestExtraHeaders{xcdp: true, auth: true}
}

func (h requestHeaderVerify) getBase() (headers map[string][]string) {
	return map[string][]string{
		hkeyHost:           []string{hvalServicesHost},
		hkeyAcceptEncoding: []string{hvalAcceptEncodingIdentity},
		hkeyAccept:         []string{hvalAcceptEncodingAppJSON},
		hkeyContentType:    []string{hvalAcceptEncodingAppJSON},
		hkeyUa:             []string{hvalUserAgent},
	}
}

func (h *requestHeaderBuildInfo) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderBuildInfo.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h requestHeaderBuildInfo) getBase() (headers map[string][]string) {
	return genericBuildBaseHeaders
}

func (h requestHeaderBuildInfo) getExtra() localRequestExtraHeaders {
	return genericExtraHeaders
}

func (h *requestHeaderBranchInfo) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderBranchInfo.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h requestHeaderBranchInfo) getBase() (headers map[string][]string) {
	return genericBuildBaseHeaders
}

func (h requestHeaderBranchInfo) getExtra() localRequestExtraHeaders {
	return genericExtraHeaders
}

func (h *requestHeaderLaunchArgs) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderLaunchArgs.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h requestHeaderLaunchArgs) getBase() (headers map[string][]string) {
	return genericBuildBaseHeaders
}

func (h requestHeaderLaunchArgs) getExtra() localRequestExtraHeaders {
	return genericExtraHeaders
}

func (h *requestHeaderGameCode) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderGameCode.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h requestHeaderGameCode) getBase() (headers map[string][]string) {
	return map[string][]string{
		hkeyHost:           []string{hvalServicesHost},
		hkeyAcceptEncoding: []string{hvalAcceptEncodingIdentity},
		hkeyAccept:         []string{hvalAcceptAll},
		hkeyUa:             []string{hvalUserAgent},
	}
}

func (h *requestHeaderGameCode) getExtra() localRequestExtraHeaders {
	return localRequestExtraHeaders{xcdp: true, auth: true}
}

func (h *requestHeaderUpdateQC) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderUpdateQC.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h *requestHeaderServerStatus) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderServerStatus.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h requestHeaderServerStatus) getBase() (headers map[string][]string) {
	return map[string][]string{
		hkeyHost:           []string{hvalServicesHost},
		hkeyAcceptEncoding: []string{hvalAcceptEncodingIdentity},
		hkeyAccept:         []string{hvalAcceptAll},
		hkeyUa:             []string{hvalUserAgent},
	}
}

func (h requestHeaderServerStatus) getExtra() localRequestExtraHeaders {
	return genericExtraHeaders
}

func (h requestHeaderUpdateQC) getBase() (headers map[string][]string) {
	return map[string][]string{}
}

func (h requestHeaderUpdateQC) getExtra() localRequestExtraHeaders {
	return genericExtraHeaders
}

func (h *requestHeaderUpdateLauncher) build() error {
	headers, err := getAll(h)
	if err != nil {
		logger.Errorw("requestHeaderUpdateLauncher.build: error building header", "error", err)
		return err
	}
	h.headerMapping = headers
	return nil
}

func (h requestHeaderUpdateLauncher) getBase() (headers map[string][]string) {
	return map[string][]string{}
}

func (h requestHeaderUpdateLauncher) getExtra() localRequestExtraHeaders {
	return genericExtraHeaders
}

func getAll(h localRequestHeader) (headerMapping, error) {
	all := h.getBase()
	e := h.getExtra()
	if e.xcdp {
		for k, v := range xcdpHeaders {
			all[k] = []string{v}
		}
		// config vars set during initilization
		all[hkeyXSrcFp] = []string{ConfSrcFp}
		all[hkeyXCdpAppVer] = []string{ConfXAppVer}
		all[hkeyXCdpLibVer] = []string{ConfXLibVer}
	}
	if e.auth {
		authToken, err := getAuthToken()
		if err != nil {
			return headerMapping{}, fmt.Errorf("Unable to get auth token when building headers, error: %s", err)
		}
		all[hkeyAuthorization] = []string{fmt.Sprintf("Token %s", authToken)}
	}
	return headerMapping{values: all}, nil
}
