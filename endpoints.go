// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import "fmt"

const (
	updateQCEndpoint       = "https://qc.syncore.org/launcher/v2/checkforupdate"
	updateLauncherEndpoint = "https://qc.syncore.org/qcl_latest_version.json"
)

func getAuthEndpoint() string {
	return fmt.Sprintf("%s/cdp-user/auth", ConfBaseSvc)
}

func getVerifyEndpoint() string {
	return fmt.Sprintf("%s/cdp-user/verify/.json", ConfBaseSvc)
}

func getServerStatusEndpoint() string {
	return fmt.Sprintf("%s/status/ext-server-status?product_id=5", ConfBaseSvc)
}

func getGameCodeEndpoint(projectID int) string {
	return fmt.Sprintf("%s/cdp-user/projects/%d/gamecode/.json", ConfBaseSvc, projectID)
}

func getBuildInfoEndpoint() string {
	return fmt.Sprintf("%s/projects/get_from_entitlement/%d/.json", ConfBaseBi, qcEntitlmentID)
}

func getBranchInfoEndpoint(projectID, branchID int) string {
	return fmt.Sprintf("%s/projects/%d/branches/%d/.json", ConfBaseBi, projectID, branchID)
}

func getLaunchArgsEndpoint(projectID int) string {
	return fmt.Sprintf("%s/projects/%d/.json", ConfBaseBi, projectID)
}
