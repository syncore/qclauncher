// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

var UseEntitlementAPI = true

func SetEntitlementAPI() {
	l := newLauncherClient(defTimeout)
	UseEntitlementAPI = l.getEntitlementAPIValue()
}

func (lc *launcherClient) getEntitlementAPIValue() bool {
	entitlement, err := lc.checkEntitlement()
	if err != nil {
		logger.Errorw("Error occurred while checking entitlement check API response, using default value of false", "error", err)
		return false
	}
	logger.Debugf("Use entitlement API: %v", entitlement)
	return entitlement.UseEntitlementAPI
}
