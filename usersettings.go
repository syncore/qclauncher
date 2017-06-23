// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"errors"
	"fmt"
	"os"

	"strings"

	"github.com/lxn/walk"
	wd "github.com/lxn/walk/declarative"
)

var title = fmt.Sprintf("QCLauncher %.2f by syncore", version)

const embeddedLogoPath = "../../resources/img/qclauncher.png"

func OpenSettings() {
	var mainDialog *walk.Dialog
	var logoView *walk.ImageView
	userSettings := &UserSettings{}
	logoImg, err := loadAppLogo(embeddedLogoPath)
	if err != nil {
		logger.Errorw("openSettings: error loading logo image", "error", err)
		return
	}

	if _, err := (wd.Dialog{
		AssignTo:  &mainDialog,
		Title:     title,
		FixedSize: true,
		MinSize:   wd.Size{Width: 300, Height: 150},
		MaxSize:   wd.Size{Width: 300, Height: 150},
		Size:      wd.Size{Width: 300, Height: 150},
		Layout:    wd.VBox{},
		Icon:      getAppIcon(),
		Children: []wd.Widget{
			wd.ImageView{
				AssignTo: &logoView,
				Image:    logoImg,
				MinSize:  wd.Size{Width: logoImg.Size().Width, Height: logoImg.Size().Height},
			},
			wd.Label{
				Text: "Your QC info has not been set.",
			},
			wd.PushButton{
				Text: "Set QC Information",
				OnClicked: func() {
					if cmd, err := runSettingsCollectionDialog(mainDialog, userSettings); err != nil {
						logger.Errorw("openSettings: error running config dialog", "error", err)
					} else if cmd == walk.DlgCmdOK {
						if err := saveUserSettings(userSettings); err != nil {
							logger.Errorw("openSettings: error saving data from config dialog", "error", err)
							walk.MsgBox(mainDialog, "Error",
								fmt.Sprintf("Unable to save QC data. Check %s for more information. Exiting.", LogFile),
								walk.MsgBoxIconError)
							os.Exit(1)
						}
					} else {
						return
					}
					exitMsg := "Now exiting"
					launchSteam := false
					if userSettings.SetAsNonSteamGame {
						exitMsg = "Now launching Steam to add as a non-Steam game and then exiting QCLauncher"
						launchSteam = true
					}
					walk.MsgBox(mainDialog, "Success", fmt.Sprintf("QC account information saved. %s. Re-run to play. To reset, delete the %s file.",
						exitMsg, DataFile),
						walk.MsgBoxIconInformation)
					if launchSteam {
						err := addNonSteamGame(getSteamInstallPath())
						if err != nil {
							ShowErrorMsg("Error", "Unable to launch Steam")
							return
						}
					}
					os.Exit(0)
				},
			},
		},
	}.Run(nil)); err != nil {
		logger.Fatalw("openSettings: Fatal error opening configuration UI", "error", err)
	}
}

func runSettingsCollectionDialog(owner walk.Form, userSettings *UserSettings) (int, error) {
	var qcInfoDlg *walk.Dialog
	var qcInfoBinder *walk.DataBinder
	var qcInfoDlgOkBtn, qcInfoDlgCancelBtn, qcFilePathBtn *walk.PushButton
	var hasSteam = isSteamInstalled()
	return wd.Dialog{
		AssignTo:      &qcInfoDlg,
		Title:         "Enter QC Info",
		Icon:          getAppIcon(),
		DefaultButton: &qcInfoDlgOkBtn,
		CancelButton:  &qcInfoDlgCancelBtn,
		DataBinder: wd.DataBinder{
			AssignTo:       &qcInfoBinder,
			DataSource:     userSettings,
			ErrorPresenter: wd.ToolTipErrorPresenter{},
		},
		MinSize: wd.Size{Width: 300, Height: 300},
		Layout:  wd.VBox{},
		Children: []wd.Widget{
			wd.Composite{
				Layout: wd.Grid{Columns: 2},
				Children: []wd.Widget{
					wd.Label{
						Text: "QC Username:",
					},
					wd.LineEdit{
						ColumnSpan: 2,
						Text:       wd.Bind("Username"),
					},
					wd.Label{
						ColumnSpan: 2,
						Text:       "QC Password:",
					},
					wd.LineEdit{
						ColumnSpan:   2,
						Text:         wd.Bind("Password"),
						PasswordMode: true,
					},
					wd.LineEdit{
						ColumnSpan:   2,
						Text:         wd.Bind("PasswordConfirm"),
						PasswordMode: true,
					},
					wd.Label{
						ColumnSpan: 2,
						Text:       "QC EXE Location",
					},
					wd.PushButton{
						ColumnSpan: 2,
						AssignTo:   &qcFilePathBtn,
						Text:       "Select QC EXE",
						OnClicked: func() {
							qcFilePathDialog := &walk.FileDialog{}
							qcFilePathDialog.Filter = "Quake Champions Exe File (QuakeChampions.exe)|QuakeChampions.exe*.*"
							qcFilePathDialog.Title = "Select your QuakeChampions.exe file"
							qcDefaultDir := "C:\\Program Files (x86)\\Bethesda.net Launcher\\games\\client\\bin\\pc"
							if dirExists(qcDefaultDir) {
								qcFilePathDialog.InitialDirPath = qcDefaultDir
							}
							if accepted, err := qcFilePathDialog.ShowOpen(owner); err != nil {
								logger.Errorw("runSettingsCollectionDialog: error submitting data to binder when saving QC filepath",
									"error", err)
								return
							} else if !accepted {
								return
							}
							userSettings.FilePath = qcFilePathDialog.FilePath
						},
					},
					wd.VSpacer{
						ColumnSpan: 2,
						Size:       10,
					},
					wd.Label{
						ColumnSpan: 1,
						Visible:    hasSteam,
						Text:       "Add as a non-Steam Game (for Steam overlay)",
					},
					wd.CheckBox{
						ColumnSpan: 1,
						Visible:    hasSteam,
						Checked:    wd.Bind("SetAsNonSteamGame"),
					},
				},
			},
			wd.Composite{
				Layout: wd.HBox{},
				Children: []wd.Widget{
					wd.HSpacer{},
					wd.PushButton{
						AssignTo: &qcInfoDlgOkBtn,
						Text:     "Save",
						OnClicked: func() {
							if err := qcInfoBinder.Submit(); err != nil {
								logger.Errorw("runSettingsCollectionDialog: error submitting data to binder when saving",
									"error", err)
								return
							}
							if err := validateSaveData(userSettings); err != nil {
								walk.MsgBox(owner, "Save Error", err.Error(), walk.MsgBoxIconError)
								return
							}
							qcInfoDlg.Accept()
						},
					},
					wd.PushButton{
						AssignTo:  &qcInfoDlgCancelBtn,
						Text:      "Cancel",
						OnClicked: func() { qcInfoDlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}

func saveUserSettings(us *UserSettings) error {
	if err := createEmptyFile(GetDataFilePath(), true); err != nil {
		logger.Errorw("saveUserSettings: error initializing empty data file", "error", err)
		return err
	}
	if err := saveUserCredentials(&QCUserCredentials{Username: us.Username, Password: us.Password, Token: ""}); err != nil {
		logger.Errorw("saveUserSettings: error saving user credentials", "error", err)
		return err
	}
	if err := saveQCOptions(&QCOptions{QCFilePath: us.FilePath}); err != nil {
		logger.Errorw("saveUserSettings: error saving QC options", "error", err)
		return err
	}
	if err := updateLastCheckTime(updateAll, 0); err != nil {
		logger.Errorw("saveUserSettings: error setting last update check time", "error", err)
		return err
	}
	return nil
}

func validateSaveData(us *UserSettings) error {
	if us == nil {
		return errors.New("All data must be specified")
	}
	if us.Username == "" {
		return errors.New("QC username must be specified")
	}
	if us.Password == "" || us.PasswordConfirm == "" {
		return errors.New("QC password must be specified")
	}
	if us.Password != us.PasswordConfirm {
		return errors.New("Passwords must match")
	}
	if us.FilePath == "" {
		return errors.New("QC EXE location must be specified")
	}
	if !strings.Contains(us.FilePath, "QuakeChampions.exe") {
		return errors.New("Invalid QC EXE was specified")
	}
	return newLauncherClient(10).verifyCredentials(us.Username, us.Password)
}
