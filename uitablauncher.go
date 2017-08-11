// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"

	"github.com/lxn/walk"
	wd "github.com/lxn/walk/declarative"
)

const tabLauncherTitle = "QCLauncher Settings"

var (
	launcherSettingsTabHeight = settingsWindowHeight - 132
	launcherSettingsTabWidth  = settingsWindowWidth - 20
)

func newLauncherSettingsTab(launcherSettings *LauncherSettings) *QCLSettingsTab {
	var hasSteam = isSteamInstalled()
	launcherSettingsTab := &QCLSettingsTab{}
	var cbAutoStartQC *walk.CheckBox
	tabPage := wd.TabPage{
		Title:  tabLauncherTitle,
		Layout: wd.VBox{},
		DataBinder: wd.DataBinder{
			AssignTo:       &launcherSettingsTab.DataBinder,
			DataSource:     launcherSettings,
			ErrorPresenter: wd.ToolTipErrorPresenter{},
		},
		Children: []wd.Widget{
			wd.GroupBox{
				MinSize: wd.Size{Height: launcherSettingsTabHeight, Width: launcherSettingsTabWidth},
				MaxSize: wd.Size{Height: launcherSettingsTabHeight, Width: launcherSettingsTabWidth},
				Title:   tabLauncherTitle,
				Layout:  wd.Grid{Columns: 1},
				Children: []wd.Widget{
					wd.CheckBox{
						AssignTo:    &cbAutoStartQC,
						Name:        "cbAutoStartQC",
						ToolTipText: "Skip the QCLauncher screen and auto start Quake Champions",
						Text:        `Auto-start QC and exit QCLauncher (skip QCLauncher UI completely)`,
						Enabled:     wd.Bind("!cbMinimizeOnLaunch.Checked && !cbExitOnLaunch.Checked"),
						Checked:     wd.Bind("AutoStartQC"),
						OnClicked: func() {
							if cbAutoStartQC.Checked() {
								ShowWarningMsg("Info",
									fmt.Sprintf(
										"The *next* time you run QCLauncher, QC will start directly and the QCLauncher UI will not be shown. To restore it, delete the %s file or start QCLauncher with qclauncher.exe -%s",
										DataFile, ShowMainWindowFlag), qclauncherSettingsWindow)
							}
						},
					},
					wd.CheckBox{
						Visible:     hasSteam,
						Text:        "Add as a non-Steam Game (for Steam overlay)",
						ToolTipText: `Open Steam after saving settings to add QC as a non-Steam game`,
						Checked:     wd.Bind("SetAsNonSteamGame"),
					},
					wd.CheckBox{
						Name:        "cbExitOnLaunch",
						ToolTipText: "Exit QCLauncher after launching QC. The QCLauncher UI will show on next launch",
						Text:        `Exit QCLauncher on game launch`,
						Enabled:     wd.Bind("!cbMinimizeOnLaunch.Checked && !cbAutoStartQC.Checked"),
						Checked:     wd.Bind("ExitOnLaunch"),
					},
					wd.CheckBox{
						Name:        "cbMinimizeOnLaunch",
						ToolTipText: "Minimize QCLauncher when launching Quake Champions",
						Text:        `Minimize QCLauncher on game launch`,
						Enabled:     wd.Bind("!cbAutoStartQC.Checked && !cbExitOnLaunch.Checked"),
						Checked:     wd.Bind("MinimizeOnLaunch"),
					},
					wd.CheckBox{
						Name:        "cbMinimizeToTray",
						ToolTipText: "Minimize QCLauncher to the system tray instead of the taskbar",
						Text:        `Minimize QCLauncher to system tray`,
						Checked:     wd.Bind("MinimizeToTray"),
					},
					wd.HSpacer{},
				},
			},
		},
	}
	launcherSettingsTab.TabPage = tabPage
	return launcherSettingsTab
}
