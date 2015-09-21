package gxui

import (
	"github.com/anaminus/gxui/math"
)

type TableLayout interface {
	Control

	Parent

	SetGrid(rows, columns int)
	// Add child at cell {x, y} with size of {w, h}
	SetChildAt(x, y, w, h int, child Control) *Child
	RemoveChild(child Control)
	SetColumnWeight(col, weight int)
	SetRowWeight(row, weight int)
	// Give the table a constant size. If an axis is less than 0, then that
	// axis will instead fill the available space.
	SetDesiredSize(math.Size)
}
