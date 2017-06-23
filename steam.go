// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

func getSteamInstallPath() string {
	installed, steamBasePath := getSteamRegistryInfo()
	if installed {
		return buildFullSteamPath(steamBasePath)
	}
	return ""
}

func getSteamRegistryInfo() (bool, string) {
	w64key, w64err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Wow6432Node\Valve\Steam`, registry.QUERY_VALUE)
	if w64err != nil {
		// couldn't find 64bit, check 32bit
		logger.Debugw("getSteamRegistryInfo: error getting 64-bit Steam info", "error", w64err)
		w32key, w32err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Valve\Steam`, registry.QUERY_VALUE)
		if w32err != nil {
			logger.Debugw("getSteamRegistryInfo: error getting 32-bit Steam info", "error", w32err)
			return false, ""
		}
		defer w32key.Close()
		s, _, werr := w32key.GetStringValue("InstallPath")
		if werr == nil {
			return s != "", s
		}
		logger.Errorw("getSteamRegistryInfo: error extracting 32-bit InstallPath value", "error", werr)
	} else {
		defer w64key.Close()
		s, _, werr := w64key.GetStringValue("InstallPath")
		if werr == nil {
			return s != "", s
		}
		logger.Errorw("getSteamRegistryInfo: error extracting 64-bit InstallPath value", "error", werr)
	}
	return false, ""
}

func buildFullSteamPath(steamBasePath string) string {
	return fmt.Sprintf("%s\\Steam.exe", steamBasePath)
}

func isSteamInstalled() bool {
	installed, _ := getSteamRegistryInfo()
	return installed
}

func addNonSteamGame(steamFilePath string) error {
	s := exec.Command(steamFilePath, "steam://AddNonSteamGame")
	if err := s.Start(); err != nil {
		logger.Errorw("addNonSteamGame: error launching Steam", "error", err)
		return err
	}
	return nil
}
