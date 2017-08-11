// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"github.com/lxn/walk"
	wd "github.com/lxn/walk/declarative"
)

const tabExpTitle = "QC Experimental Settings"

type QCLSettingsTab struct {
	wd.TabPage
	DataBinder *walk.DataBinder
}

func newQCExperimentalSettingsTab(expSettings *QCExperimentalSettings) *QCLSettingsTab {
	qcExperimentalSettingsTab := &QCLSettingsTab{}
	var neMaxFPS, neMaxFPSMinimized *walk.NumberEdit
	ttMaxFPS, ttMaxFPSMinimized := "Experimental: Limit or 'cap' the maximum FPS during the game",
		"Experimental: Limit or 'cap' the maximum FPS when QC is minimized"
	height, width, fmin, fmax := 20, 77, 0, 2000
	maxFPSCmp, maxFPSMinCmp := wd.Composite{
		Layout:     wd.Grid{Columns: 2, MarginsZero: true},
		ColumnSpan: 2,
		Children: []wd.Widget{
			wd.NumberEdit{
				AssignTo:    &neMaxFPS,
				Enabled:     wd.Bind("cbUseMaxFPSLimit.Checked"),
				Value:       wd.Bind("MaxFPSLimit"),
				MinValue:    float64(fmin),
				MaxValue:    float64(fmax),
				MinSize:     wd.Size{Height: height, Width: width},
				MaxSize:     wd.Size{Height: height, Width: width},
				Suffix:      " FPS",
				ToolTipText: ttMaxFPS,
			},
		},
	}, wd.Composite{
		Layout:     wd.Grid{Columns: 2, MarginsZero: true},
		ColumnSpan: 2,
		Children: []wd.Widget{
			wd.NumberEdit{
				AssignTo:    &neMaxFPSMinimized,
				Enabled:     wd.Bind("cbUseMaxFPSLimitMinimized.Checked"),
				Value:       wd.Bind("MaxFPSLimitMinimized"),
				MinValue:    float64(fmin),
				MaxValue:    float64(fmax),
				MinSize:     wd.Size{Height: height, Width: width},
				MaxSize:     wd.Size{Height: height, Width: width},
				Suffix:      " FPS",
				ToolTipText: ttMaxFPSMinimized,
			},
		},
	}
	tabPage := wd.TabPage{
		Title:  tabExpTitle,
		Layout: wd.HBox{},
		DataBinder: wd.DataBinder{
			AssignTo:       &qcExperimentalSettingsTab.DataBinder,
			DataSource:     expSettings,
			ErrorPresenter: wd.ToolTipErrorPresenter{},
		},
		Children: []wd.Widget{
			wd.GroupBox{
				Title:  tabExpTitle,
				Layout: wd.Grid{Columns: 2},
				Children: []wd.Widget{
					wd.Label{Text: "These experimental settings may be removed from the game at any time!",
						TextColor: walk.RGB(210, 0, 0), ColumnSpan: 2},
					wd.Label{Text: "Use at your own risk, these are undocumented and may have zero/unpredictable effects!",
						TextColor: walk.RGB(210, 0, 0), ColumnSpan: 2},
					wd.VSpacer{Size: 1, ColumnSpan: 2},
					wd.CheckBox{
						Name:        "cbUseMaxFPSLimit",
						ToolTipText: ttMaxFPS,
						Text:        `Set Max FPS Limit`,
						Checked:     wd.Bind("UseMaxFPSLimit"),
					},
					maxFPSCmp,
					wd.CheckBox{
						Text:        `Set Max FPS Limit (when QC is minimized)`,
						Name:        "cbUseMaxFPSLimitMinimized",
						ToolTipText: ttMaxFPSMinimized,
						Checked:     wd.Bind("UseMaxFPSLimitMinimized"),
					},
					maxFPSMinCmp,
					wd.CheckBox{
						ColumnSpan:  2,
						Text:        `Use FPS smoothing`,
						ToolTipText: "Experimental: use FPS smoothing (may have no effect)",
						Checked:     wd.Bind("UseFPSSmoothing"),
					},
					wd.VSpacer{
						ColumnSpan: 2,
					},
				},
			},
		},
	}
	qcExperimentalSettingsTab.TabPage = tabPage
	return qcExperimentalSettingsTab
}
