package mixins

import (
	"github.com/anaminus/gxui"
	"github.com/anaminus/gxui/math"
	"github.com/anaminus/gxui/mixins/base"
)

type Cell struct {
	x, y, w, h int
}

func (c Cell) AtColumn(x int) bool {
	return c.x <= x && c.x+c.w >= x
}

func (c Cell) AtRow(y int) bool {
	return c.y <= y && c.y+c.h >= y
}

type TableLayoutOuter interface {
	base.ContainerOuter
}

type TableLayout struct {
	base.Container

	outer TableLayoutOuter

	grid    map[gxui.Control]Cell
	rows    int
	columns int

	rowWeight      []int
	colWeight      []int
	rowTotalWeight int
	colTotalWeight int

	desiredSize math.Size
}

func (l *TableLayout) Init(outer TableLayoutOuter, theme gxui.Theme) {
	l.Container.Init(outer, theme)
	l.outer = outer
	l.grid = make(map[gxui.Control]Cell)
	l.desiredSize = math.Size{-1, -1}

	// Interface compliance test
	_ = gxui.TableLayout(l)
}

func (l *TableLayout) LayoutChildren() {
	s := l.outer.Size().Contract(l.outer.Padding())
	o := l.outer.Padding().LT()

	cw, ch := float64(s.W), float64(s.H)

	ctweight, rtweight := float64(l.colTotalWeight), float64(l.rowTotalWeight)

	var cr math.Rect

	for _, c := range l.outer.Children() {
		cm := c.Control.Margin()
		cell := l.grid[c.Control]

		var xs float64
		{
			var weight float64
			for i := 0; i < cell.x; i++ {
				weight += float64(l.colWeight[i])
			}
			xs = weight / ctweight
		}
		var ys float64
		{
			var weight float64
			for i := 0; i < cell.y; i++ {
				weight += float64(l.rowWeight[i])
			}
			ys = weight / rtweight
		}
		var ws float64
		{
			var weight float64
			for i := cell.x; i < cell.x+cell.w; i++ {
				weight += float64(l.colWeight[i])
			}
			ws = weight / ctweight
		}
		var hs float64
		{
			var weight float64
			for i := cell.y; i < cell.y+cell.h; i++ {
				weight += float64(l.rowWeight[i])
			}
			hs = weight / rtweight
		}

		x, y := cw*xs, ch*ys
		w, h := x+cw*ws, y+ch*hs

		cr = math.CreateRect(int(x), int(y), int(w), int(h)).Contract(cm)

		c.Layout(cr.Offset(o).Canon())
	}
}

func (l *TableLayout) DesiredSize(min, max math.Size) math.Size {
	size := l.desiredSize
	if size.W < 0 {
		size.W = max.W
	}
	if size.H < 0 {
		size.H = max.H
	}
	return size
}

func (l *TableLayout) SetDesiredSize(size math.Size) {
	if size != l.desiredSize {
		l.desiredSize = size
		l.outer.Relayout()
	}
}
}

func (l *TableLayout) SetGrid(columns, rows int) {
	if l.columns != columns {
		if l.columns > columns {
			for c := l.columns; c > columns; c-- {
				for _, cell := range l.grid {
					if cell.AtColumn(c) {
						panic("Can't remove column with cells")
					}
				}
				l.columns--
			}
			l.colWeight = l.colWeight[:l.columns]
			l.colTotalWeight = 0
			for _, w := range l.colWeight {
				l.colTotalWeight += w
			}
		} else {
			l.columns = columns
			a := make([]int, l.columns-len(l.colWeight))
			for i := range a {
				a[i] = 1
			}
			l.colWeight = append(l.colWeight, a...)
			l.colTotalWeight += len(a)
		}
	}

	if l.rows != rows {
		if l.rows > rows {
			for r := l.rows; r > rows; r-- {
				for _, cell := range l.grid {
					if cell.AtRow(r) {
						panic("Can't remove row with cells")
					}
				}
				l.rows--
			}
			l.rowWeight = l.rowWeight[:l.rows]
			l.rowTotalWeight = 0
			for _, w := range l.rowWeight {
				l.rowTotalWeight += w
			}
		} else {
			l.rows = rows
			a := make([]int, l.rows-len(l.rowWeight))
			for i := range a {
				a[i] = 1
			}
			l.rowWeight = append(l.rowWeight, a...)
			l.rowTotalWeight += len(a)
		}
	}

	if l.rows != rows || l.columns != columns {
		l.LayoutChildren()
	}
}

func (l *TableLayout) SetChildAt(x, y, w, h int, child gxui.Control) *gxui.Child {
	if x+w > l.columns || y+h > l.rows {
		panic("Cell is out of grid")
	}

	for _, c := range l.grid {
		if c.x+c.w > x && c.x < x+w && c.y+c.h > y && c.y < y+h {
			panic("Cell already has a child")
		}
	}

	l.grid[child] = Cell{x, y, w, h}
	return l.Container.AddChild(child)
}

func (l *TableLayout) RemoveChild(child gxui.Control) {
	delete(l.grid, child)
	l.Container.RemoveChild(child)
}

func (l *TableLayout) SetColumnWeight(col, weight int) {
	if col < 0 || col >= len(l.rowWeight) {
		panic("Column is out of grid")
	}
	if weight < 0 {
		panic("Weight cannot be less than 0")
	}
	if l.colWeight[col] != weight {
		w := l.colWeight[col]
		l.colWeight[col] = weight
		l.colTotalWeight += -w + weight
		l.LayoutChildren()
	}
}

func (l *TableLayout) SetRowWeight(row, weight int) {
	if row < 0 || row >= len(l.rowWeight) {
		panic("Row is out of grid")
	}
	if weight < 0 {
		panic("Weight cannot be less than 0")
	}
	if l.rowWeight[row] != weight {
		w := l.rowWeight[row]
		l.rowWeight[row] = weight
		l.rowTotalWeight += -w + weight
		l.LayoutChildren()
	}
}
