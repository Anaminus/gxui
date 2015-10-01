package gxui

import (
	"github.com/anaminus/gxui/math"
)

type MenuOverlay interface {
	Control
	Parent

	// AddMenu adds a control to the overlay, positioned next to the given
	// target. Orientation determines which axis the menu will move on to
	// avoid overlapping the target.
	AddMenu(menu Control, target math.Rect, orientation Orientation) *Child
	RemoveMenu(menu Control)
	Clear()
}
