// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"

	"github.com/lxn/walk"
	wd "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

const (
	settingsWindowWidth  = 500
	settingsWindowHeight = 425
)

type QCLSettingsWindow struct {
	*walk.MainWindow
}

var qclauncherSettingsWindow *QCLSettingsWindow

func (sw *QCLSettingsWindow) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_CLOSE: // manually dispose of window resource (window frame 'X' button)
		sw.cleanup()
	}
	return sw.MainWindow.WndProc(hwnd, msg, wParam, lParam)
}

func newSettingsWindow(cfg *Configuration) *QCLSettingsWindow {
	settingsWindow := &QCLSettingsWindow{}
	settingsTabs := getSettingsTabs(cfg)
	settingsTabPages := getSettingsTabPages(settingsTabs)
	icon := getAppIcon()
	if err := (wd.MainWindow{
		AssignTo:             &settingsWindow.MainWindow,
		Title:                "Configuration",
		Icon:                 icon,
		UseCustomWindowStyle: true,
		CustomWindowStyle:    win.WS_CAPTION | win.WS_SYSMENU | win.WS_MINIMIZEBOX,
		MinSize:              wd.Size{Width: settingsWindowWidth, Height: settingsWindowHeight},
		MaxSize:              wd.Size{Width: settingsWindowWidth, Height: settingsWindowHeight},
		Size:                 wd.Size{Width: settingsWindowWidth, Height: settingsWindowHeight},
		Layout:               wd.VBox{},
		Children: []wd.Widget{
			wd.TabWidget{
				Pages: settingsTabPages,
			},
			wd.Composite{
				Layout: wd.HBox{},
				Children: []wd.Widget{
					wd.HSpacer{},
					wd.PushButton{
						Text:        "Save All",
						ToolTipText: "Save all settings and verify account with Bethesda",
						OnClicked: func() {
							if err := settingsTabsBinderSubmit(settingsTabs); err != nil {
								logger.Errorw(fmt.Sprintf("%s: error submitting all tabs' saved settings from binder", GetCaller()),
									"error", err)
								ShowErrorMsg("Save Error", err.Error(), settingsWindow.MainWindow)
								return
							}
							if err := saveConfiguration(cfg); err != nil {
								ShowErrorMsg("Save Error", err.Error(), settingsWindow.MainWindow)
								return
							}
							defer settingsWindow.close()
							if err := handlePostSave(cfg); err != nil {
								logger.Errorw(fmt.Sprintf("%s: error executing post-save function", GetCaller()), "error", err)
								return
							}
						},
					},
					wd.PushButton{
						Text:        "Cancel",
						ToolTipText: "Cancel",
						OnClicked: func() {
							settingsWindow.close()
						},
					},
					wd.PushButton{
						Text:        "Reset All",
						ToolTipText: "Clear all settings and account information",
						OnClicked: func() {
							result := walk.MsgBox(settingsWindow.MainWindow, "Reset All Settings",
								"Delete All Saved Settings?", walk.MsgBoxYesNo)
							if result == win.IDYES {
								defer settingsWindow.close()
								handlePostReset()
								return
							} else if result == win.IDNO || result == 0 {
								return
							}
						},
					},
				},
			},
		},
	}).Create(); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error during creation of settings window", GetCaller()), "error", err)
	}
	if qclauncherMainWindow != nil {
		if err := settingsWindow.SetOwner(qclauncherMainWindow.MainWindow); err != nil {
			logger.Errorw(fmt.Sprintf("%s: unable to set owner of settings window to main window", GetCaller()), "error", err)
		}
	}
	// for overriding WndProc
	if err := walk.InitWrapperWindow(settingsWindow); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error during creation of configuration UI wrapper window", GetCaller()), "error", err)
	}
	qclauncherSettingsWindow = settingsWindow
	return settingsWindow
}

func getSettingsTabs(cfg *Configuration) []*QCLSettingsTab {
	qcCoreSettingsTab := newQCCoreSettingsTab(cfg.Core)
	qcExperimentalSettingsTab := newQCExperimentalSettingsTab(cfg.Experimental)
	launcherSettingsTab := newLauncherSettingsTab(cfg.Launcher)
	return []*QCLSettingsTab{
		qcCoreSettingsTab,
		qcExperimentalSettingsTab,
		launcherSettingsTab,
	}
}

func getSettingsTabPages(q []*QCLSettingsTab) []wd.TabPage {
	tp := []wd.TabPage{}
	for _, t := range q {
		tp = append(tp, t.TabPage)
	}
	return tp
}

func settingsTabsBinderSubmit(tabs []*QCLSettingsTab) error {
	for _, t := range tabs {
		if err := t.DataBinder.Submit(); err != nil {
			logger.Errorw(fmt.Sprintf("%s: error submitting data to binder when saving settings", GetCaller()), "error", err)
			return err
		}
	}
	return nil
}

func handlePostSave(cfg *Configuration) error {
	// re-read to receive unencrypted credentials
	if cfg.Core.Username != "" || cfg.Core.Password != "" {
		qcs := &QCCoreSettings{}
		err := Get(qcs)
		if err != nil {
			logger.Errorw(fmt.Sprintf("%s: error re-reading non-default core settings", GetCaller()), "error", err)
			return err
		}
		qclauncherMainWindow.setSignedInName(qcs.Username)
	}
	qclauncherMainWindow.enableLaunchButton(true)
	qclauncherMainWindow.enableLaunchTrayAction(true)
	qclauncherMainWindow.updateMinimizeSettings(cfg.Launcher.MinimizeToTray)
	// done from 'Configure' context menu option from tray when window is minimized
	if qclauncherMainWindow.TrayIcon.Visible() && !cfg.Launcher.MinimizeToTray {
		qclauncherMainWindow.showTrayIcon(false)
		qclauncherMainWindow.restore(true)
	}
	ShowInfoMsg("Success", "Settings were saved successfully. Click \"Play\" to launch.", qclauncherSettingsWindow)
	return nil
}

func handlePostReset() {
	DeleteConfiguration(false)
	qclauncherMainWindow.updateMinimizeSettings(false)
	qclauncherMainWindow.enableLaunchButton(false)
	qclauncherMainWindow.enableLaunchTrayAction(false)
	qclauncherMainWindow.setSignedInName("")
	// done from 'Configure' context menu option from tray when window is minimized
	if !qclauncherMainWindow.Visible() {
		qclauncherMainWindow.showTrayIcon(false)
		qclauncherMainWindow.restore(true)
	}
	ShowInfoMsg("Success", "All settings were reset. Click \"Configure\" to set up.", qclauncherSettingsWindow)
}

func (sw *QCLSettingsWindow) close() {
	isCollectingSettings = false
	if sw == nil || sw.MainWindow == nil {
		return
	}
	if err := sw.MainWindow.Close(); err != nil { // automatically disposes (via 'Cancel'/'Reset' buttons)
		logger.Errorw(fmt.Sprintf("%s: error closing settings window", GetCaller()), "error", err)
	}
}

func (sw *QCLSettingsWindow) cleanup() {
	isCollectingSettings = false
	if sw == nil || sw.MainWindow == nil {
		return
	}
	sw.MainWindow.Dispose()
}
