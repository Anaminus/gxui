package gxui

import (
	"github.com/anaminus/gxui/math"
)

type TableLayout interface {
	Control

	Parent

	Grid() (columns, row int)
	SetGrid(columns, rows int)
	// Add child at cell {x, y} with size of {w, h}
	SetChildAt(x, y, w, h int, child Control) *Child
	RemoveChild(child Control)
	SetColumnWeight(col, weight int)
	SetRowWeight(row, weight int)
	// Give the table a constant size. If an axis is less than 0, then that
	// axis will instead fill the available space.
	SetDesiredSize(math.Size)
	// Returns how the size is clamped on either axis.
	SizeClamped() (w bool, h bool)
	// Sets how the size is clamped on either axis.
	SetSizeClamped(w bool, h bool)
}
