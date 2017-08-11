// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"

	"github.com/lxn/walk"
	wd "github.com/lxn/walk/declarative"
)

type Language struct {
	LangCode string
	Name     string
}

const tabCoreTitle = "QC Core Settings"

func newQCCoreSettingsTab(qcCoreSettings *QCCoreSettings) *QCLSettingsTab {
	qcCoreSettingsTab := &QCLSettingsTab{}
	tabPage := wd.TabPage{
		Title:  tabCoreTitle,
		Layout: wd.HBox{},
		DataBinder: wd.DataBinder{
			AssignTo:       &qcCoreSettingsTab.DataBinder,
			DataSource:     qcCoreSettings,
			ErrorPresenter: wd.ToolTipErrorPresenter{},
		},
		Children: []wd.Widget{
			wd.GroupBox{
				Title:  tabCoreTitle,
				Layout: wd.Grid{Columns: 3},
				Children: []wd.Widget{
					wd.VSpacer{ColumnSpan: 3, Size: 1},
					wd.Label{
						ColumnSpan: 2,
						Text:       "QC Username:",
					},
					wd.LineEdit{
						ColumnSpan:  2,
						Text:        wd.Bind("Username"),
						ToolTipText: `Enter your Bethesda.net username`,
					},
					wd.Label{
						ColumnSpan: 2,
						Text:       "QC Password:",
					},
					wd.LineEdit{
						ColumnSpan:   2,
						Text:         wd.Bind("Password"),
						ToolTipText:  `Enter your Bethesda.net password`,
						PasswordMode: true,
					},
					wd.Label{
						ColumnSpan: 2,
						Text:       "QC Language:",
					},
					wd.ComboBox{
						ColumnSpan:    2,
						Editable:      false,
						Value:         wd.Bind("Language"),
						ToolTipText:   `Select the language for the in-game QC interface`,
						BindingMember: "LangCode",
						DisplayMember: "Name",
						Model: []*Language{
							{"en", "English"},
							{"es", "Español"},
							{"es-419", "Español (Latinoamérica)"},
							{"de", "Deutsch"},
							{"fr", "Français"},
							{"it", "Italiano"},
							{"pl", "Polski"},
							{"pt", "Português (Brasil)"},
							{"ru", "Русский"},
						},
					},
					wd.Label{
						ColumnSpan: 2,
						Text:       "QC EXE Location",
					},
					wd.PushButton{
						ColumnSpan:  2,
						Text:        "Select QC EXE",
						ToolTipText: "Select your Quake Champions.exe file location",
						OnClicked: func() {
							qcFilePathDialog := &walk.FileDialog{}
							qcFilePathDialog.Filter = "Quake Champions Exe File (QuakeChampions.exe)|QuakeChampions.exe*.*"
							qcFilePathDialog.Title = "Select your QuakeChampions.exe file"
							qcDefaultDir := "C:\\Program Files (x86)\\Bethesda.net Launcher\\games\\client\\bin\\pc"
							if dirExists(qcDefaultDir) {
								qcFilePathDialog.InitialDirPath = qcDefaultDir
							}
							if accepted, err := qcFilePathDialog.ShowOpen(nil); err != nil {
								logger.Errorw(fmt.Sprintf("%s: error submitting data to binder when saving QC filepath", GetCaller()),
									"error", err)
								return
							} else if !accepted {
								return
							}
							qcCoreSettings.FilePath = qcFilePathDialog.FilePath
						},
					},
					wd.Label{ColumnSpan: 2, Text: qcCoreSettings.FilePath},
					wd.VSpacer{ColumnSpan: 2, Size: 4},
				},
			},
		},
	}
	qcCoreSettingsTab.TabPage = tabPage
	return qcCoreSettingsTab
}
