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
	embeddedLogoPath = "../../resources/img/qclauncher.png"
	mainWindowWidth  = 300
	mainWindowHeight = 160
	loggedInAs       = "Logged in as"
)

type QCLMainWindow struct {
	*walk.MainWindow
	TrayIcon *walk.NotifyIcon
	Options  *QCLMainWindowOptions
	Binder   *walk.DataBinder
}

type QCLMainWindowOptions struct {
	MinimizeToTray bool
	CanLaunch      bool
	SignedInName   string
}

var qclauncherMainWindow *QCLMainWindow

func LoadUI(cfg *Configuration) {
	signedInName := formatSignedInName(cfg.Core.Username)
	m := newMainWindow(cfg, &QCLMainWindowOptions{
		MinimizeToTray: cfg.Launcher.MinimizeToTray,
		SignedInName:   signedInName,
		CanLaunch:      FileExists(GetDataFilePath())})

	if qclauncherMainWindow == nil {
		qclauncherMainWindow = m
	}
	m.Run()
}

func (qm *QCLMainWindow) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_SIZE:
		if wParam == 0 { // SIZE_RESTORED
			qm.setMainWindowSize()
			if qm.Options.MinimizeToTray {
				qm.showTrayIcon(false)
			}
		}
		if wParam == 1 { // SIZE_MINIMIZED
			if qm.Options.MinimizeToTray {
				qm.MainWindow.SetVisible(false)
				qm.showTrayIcon(true)
			}
		}
	case win.WM_DESTROY:
		qm.exitFromMainWindow()
	}
	return qm.MainWindow.WndProc(hwnd, msg, wParam, lParam)
}

func newMainWindow(cfg *Configuration, opts *QCLMainWindowOptions) *QCLMainWindow {
	var logoView *walk.ImageView
	var mwBinder *walk.DataBinder
	icon := getAppIcon()
	logoImg, err := loadAppLogo(embeddedLogoPath)
	if err != nil {
		logger.Errorw(fmt.Sprintf("%s: error loading logo image", GetCaller()), "error", err)
	}
	mainWindow := &QCLMainWindow{}
	mainWindow.Options = opts
	binder := wd.DataBinder{
		AssignTo:       &mwBinder,
		DataSource:     mainWindow.Options,
		ErrorPresenter: wd.ToolTipErrorPresenter{},
	}
	if err := (wd.MainWindow{
		AssignTo:             &mainWindow.MainWindow,
		Title:                title,
		MinSize:              wd.Size{Width: mainWindowWidth, Height: mainWindowHeight},
		MaxSize:              wd.Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Size:                 wd.Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Layout:               wd.VBox{},
		DataBinder:           binder,
		UseCustomWindowStyle: true,
		CustomWindowStyle:    win.WS_CAPTION | win.WS_SYSMENU | win.WS_MINIMIZEBOX,
		Icon:                 icon,
		Children: []wd.Widget{
			wd.ImageView{
				AssignTo:    &logoView,
				Image:       logoImg,
				ToolTipText: title,
				MinSize:     wd.Size{Width: logoImg.Size().Width, Height: logoImg.Size().Height},
				OnMouseDown: func(x, y int, button walk.MouseButton) {
					if button != walk.LeftButton {
						return
					}
					ShowInfoMsg(fmt.Sprintf("QCLauncher %.2f", version),
						fmt.Sprintf("%s\nhttp://github.com/syncore/qclauncher\nSupport Quake!", title), mainWindow.MainWindow)
				},
			},
			wd.Composite{
				Layout: wd.HBox{MarginsZero: true},
				Children: []wd.Widget{
					wd.PushButton{
						Text:        "Play",
						ToolTipText: "Launch Quake Champions",
						OnClicked: func() {
							if err := Launch(); err != nil {
								logger.Errorw(fmt.Sprintf("%s: launch error", GetCaller()),
									"error", err)
								if IsErrAlreadyRunning(err) || IsErrHashMismatch(err) || IsErrAuthFailed(err) {
									ShowErrorMsg("Error", err.Error(), mainWindow.MainWindow)
								} else {
									ShowErrorMsg("Error", UILaunchErrorMsg, mainWindow.MainWindow)
								}
							}
						},
						Enabled: wd.Bind("CanLaunch"),
					},
					wd.PushButton{
						Text:        "Configure",
						ToolTipText: "Configure your settings and account information",
						OnClicked: func() {
							configureSettings()
						},
					},
					wd.PushButton{
						Text:        "Exit",
						ToolTipText: "Exit QCLauncher",
						OnClicked:   func() { mainWindow.exitFromMainWindow() },
					},
				},
			},
			wd.VSpacer{},
			wd.Composite{
				Layout: wd.Grid{MarginsZero: true},
				Children: []wd.Widget{
					wd.Label{Text: wd.Bind("SignedInName"), Row: 0, Column: 0},
					wd.HSpacer{Row: 0, Column: 1},
					wd.Label{Text: fmt.Sprintf("v%.2f", version), Row: 0, Column: 2},
				},
			},
		},
	}.Create()); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error during creation of configuration UI", GetCaller()), "error", err)
	}
	// for overriding WndProc
	if err := walk.InitWrapperWindow(mainWindow); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error during creation of configuration UI wrapper window", GetCaller()), "error", err)
	}
	mainWindow.setTrayIcon(icon, cfg)
	mainWindow.setMainWindowSize()
	mainWindow.Binder = mwBinder
	return mainWindow
}

func (qm *QCLMainWindow) setTrayIcon(icon *walk.Icon, cfg *Configuration) {
	trayIcon, err := walk.NewNotifyIcon()
	if err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error creating system tray icon", GetCaller()), "error", err)
	}
	if err := trayIcon.SetIcon(icon); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error setting system tray icon", GetCaller()), "error", err)
	}
	if err := trayIcon.SetToolTip(title); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error setting system tray icon tooltip", GetCaller()), "error", err)
	}
	trayIcon.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		qm.showTrayIcon(false)
		qm.restore(true)
	})
	actionPlay := walk.NewAction()
	if err := actionPlay.SetEnabled(qm.Options.CanLaunch); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error setting play action status", GetCaller()), "error", err)
	}
	if err := actionPlay.SetText("P&lay"); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error setting play action", GetCaller()), "error", err)
	}
	actionPlay.Triggered().Attach(func() {
		if err := Launch(); err != nil {
			logger.Errorw(fmt.Sprintf("%s: error occurred while executing the launch process.", GetCaller()),
				"error", err)
			if IsErrAlreadyRunning(err) || IsErrHashMismatch(err) || IsErrAuthFailed(err) {
				ShowErrorMsg("Error", err.Error(), nil)
			} else {
				ShowErrorMsg("Error", UILaunchErrorMsg, nil)
			}
		}
	})
	if err := trayIcon.ContextMenu().Actions().Add(actionPlay); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error adding play action", GetCaller()), "error", err)
	}
	actionConfigure := walk.NewAction()
	if err := actionConfigure.SetText("C&onfigure"); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error setting configure action", GetCaller()), "error", err)
	}
	actionConfigure.Triggered().Attach(func() { configureSettings() })
	if err := trayIcon.ContextMenu().Actions().Add(actionConfigure); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error adding configure action", GetCaller()), "error", err)
	}
	actionExit := walk.NewAction()
	if err := actionExit.SetText("E&xit"); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error setting exit action", GetCaller()), "error", err)
	}
	actionExit.Triggered().Attach(func() { qm.exitFromMainWindow() })
	if err := trayIcon.ContextMenu().Actions().Add(actionExit); err != nil {
		logger.FatalUIw(fmt.Sprintf("%s: Fatal error adding exit action", GetCaller()), "error", err)
	}
	qm.TrayIcon = trayIcon
}

func (qm *QCLMainWindow) cleanupTrayIcon() {
	if qm == nil || qm.TrayIcon == nil {
		return
	}
	qm.showTrayIcon(false)
	if err := qm.TrayIcon.Dispose(); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error cleaning up system tray icon resources", GetCaller()), "error", err)
	}
}

func (qm *QCLMainWindow) showTrayIcon(visible bool) {
	if qm == nil || qm.TrayIcon == nil {
		return
	}
	if err := qm.TrayIcon.SetVisible(visible); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error setting system tray icon visibility", GetCaller()), "error", err)
	}
}

func (qm *QCLMainWindow) minimize(enqueue bool) {
	if enqueue {
		// some minimize events (i.e. post-QC launch auto minimize to tray) need to be enqueued at a later time
		// using the synchronization mechanism (mutex) so that the window will re-paint properly when restored
		qm.Synchronize(func() {
			win.ShowWindow(qm.Handle(), win.SW_MINIMIZE)
		})
	} else {
		win.ShowWindow(qm.Handle(), win.SW_MINIMIZE)
	}
}

func (qm *QCLMainWindow) restore(enqueue bool) {
	if enqueue {
		qm.Synchronize(func() {
			qm.SwitchToThisWindow() // tray
		})
	} else {
		win.ShowWindow(qm.Handle(), win.SW_RESTORE) // or: win.SendMessage(qm.Handle(), win.WM_SYSCOMMAND, win.SC_RESTORE, 0)
	}
}

func (qm *QCLMainWindow) setMainWindowSize() {
	if err := qm.MainWindow.SetSize(walk.Size{Width: mainWindowWidth, Height: mainWindowHeight}); err != nil {
		logger.Errorw(fmt.Sprintf("%s: error setting main window size", GetCaller()), "error", err)
	}
}

func (qm *QCLMainWindow) setSignedInName(username string) {
	if qm == nil || qm.Options == nil {
		return
	}
	defer qm.MainWindow.SetSuspended(false) // re-draw to clear overlapping 'logged in as' text
	signedInName := formatSignedInName(username)
	qm.Options.SignedInName = signedInName
	qm.MainWindow.SetSuspended(true)
	if err := qm.refreshBoundSettings(); err != nil {
		logger.Error(fmt.Sprintf("%s: %s", GetCaller(), err))
	}
}

func (qm *QCLMainWindow) enableLaunchButton(enabled bool) {
	if qm == nil || qm.Options == nil {
		return
	}
	qm.Options.CanLaunch = enabled
	if err := qm.refreshBoundSettings(); err != nil {
		logger.Error(fmt.Sprintf("%s: %s", GetCaller(), err))
	}
}

func (qm *QCLMainWindow) enableLaunchTrayAction(enabled bool) {
	if qm == nil || qm.TrayIcon == nil {
		return
	}
	playAction := qm.TrayIcon.ContextMenu().Actions().At(0)
	if err := playAction.SetEnabled(enabled); err != nil {
		logger.Error(fmt.Sprintf("%s: unable to change launch tray action status: %s", GetCaller(), err))
	}
}

func (qm *QCLMainWindow) updateMinimizeSettings(minimizeToTray bool) {
	if qm == nil || qm.Options == nil {
		return
	}
	qm.Options.MinimizeToTray = minimizeToTray
	if err := qm.refreshBoundSettings(); err != nil {
		logger.Error(fmt.Sprintf("%s: %s", GetCaller(), err))
	}
}

func formatSignedInName(username string) string {
	if username == "" {
		// allocate additional space (tabs) in string to get this to re-draw properly; TODO: examine this more, if it's an issue
		return "Not logged in.\t\t\t\t\t"
	}
	return fmt.Sprintf("%s %s", loggedInAs, username)
}

func (qm *QCLMainWindow) refreshBoundSettings() error {
	if err := qm.Binder.Reset(); err != nil {
		return fmt.Errorf("%s: Error resetting data binder: %s", GetCaller(), err)
	}
	return nil
}

func (qm *QCLMainWindow) exitFromMainWindow() {
	qm.cleanupTrayIcon()
	qm.MainWindow.Dispose()
	exitFromUI()
}
