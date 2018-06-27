// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"
	"strings"
)

func validateResponse(r remoteResponse) error {
	return r.validate()
}

func (r *AuthResponse) validate() error {
	msg := "Auth/verify response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil auth/verify response", msg)
	}
	if r.Token == "" {
		return fmt.Errorf("%s got empty auth/verify token", msg)
	}
	if len(r.EntitlementIDs) == 0 {
		return fmt.Errorf("%s no entitlement ids present", msg)
	}
	hasQc := false
	for _, v := range r.EntitlementIDs {
		if v == qcEntitlmentID {
			hasQc = true
			break
		}
	}
	if !hasQc {
		return fmt.Errorf("%s user does not have QC access", msg)
	}
	return nil
}

func (r *EntitlementInfoResponse) validate() error {
	msg := "Entitlement info response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil entitlement info response", msg)
	}
	if len(r.Projects) == 0 {
		return fmt.Errorf("%s no project info present", msg)
	}
	if len(r.Branches) == 0 {
		return fmt.Errorf("%s no branch info present", msg)
	}
	return nil
}

func (r *BuildInfoResponse) validate() error {
	msg := "Build info response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil build info response", msg)
	}
	if len(r.Projects) == 0 {
		return fmt.Errorf("%s no project info present", msg)
	}
	if len(r.Branches) == 0 {
		return fmt.Errorf("%s no branch info present", msg)
	}
	return nil
}

func (r *BranchInfoResponse) validate() error {
	msg := "Branch info response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil branch info response", msg)
	}
	if len(r.LaunchinfoList) == 0 {
		return fmt.Errorf("%s no launch info present", msg)
	}
	return nil
}

func (r *LaunchArgsResponse) validate() error {
	msg := "Launch args response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil launch args response", msg)
	}
	return nil
}

func (r *GameCodeResponse) validate() error {
	msg := "Game code response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil game code response", msg)
	}
	if r.Gamecode == "" {
		return fmt.Errorf("%s got empty game code response", msg)
	}
	return nil
}

func (r *ServerStatusResponse) validate() error {
	msg := "Server status response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil server status response", msg)
	}
	if r.Platform.Message == "" {
		return fmt.Errorf("%s got empty platform message response", msg)
	}
	if !strings.EqualFold(r.Platform.Message, "success") {
		return fmt.Errorf("%s returned a non-successful platform message, returned: %s",
			msg, r.Platform.Message)
	}
	return nil
}

func (r *UpdateQCResponse) validate() error {
	msg := "QC update check response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil update qc response", msg)
	}
	if len(r.Hashes) == 0 {
		return fmt.Errorf("%s got empty hashes when checking for update", msg)
	}
	return nil
}

func (r *UpdateLauncherResponse) validate() error {
	msg := "QC update launcher response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil update launcher response", msg)
	}
	return nil
}

func (r *EntitlementCheckAPIResponse) validate() error {
	msg := "Entitlement check API response failed validation:"
	if r == nil {
		return fmt.Errorf("%s got nil entitlement check API response", msg)
	}
	return nil
}
