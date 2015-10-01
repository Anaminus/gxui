package mixins

import (
	"github.com/anaminus/gxui"
	"github.com/anaminus/gxui/math"
	"github.com/anaminus/gxui/mixins/base"
)

type MenuOverlayOuter interface {
	base.ContainerOuter
}

type menuTarget struct {
	area        math.Rect
	orientation gxui.Orientation
}

type MenuOverlay struct {
	base.Container
	outer   MenuOverlayOuter
	targets map[gxui.Control]menuTarget
	brush   gxui.Brush
	pen     gxui.Pen
}

func (o *MenuOverlay) Init(outer MenuOverlayOuter, theme gxui.Theme) {
	o.Container.Init(outer, theme)
	o.outer = outer
	o.targets = make(map[gxui.Control]menuTarget)

	// Interface compliance test
	_ = gxui.MenuOverlay(o)
}

func (o *MenuOverlay) LayoutChildren() {
	for _, child := range o.outer.Children() {
		bounds := o.outer.Size().Rect().Contract(o.outer.Padding())
		cs := child.Control.DesiredSize(
			math.ZeroSize,
			bounds.Size().Contract(child.Control.Margin()).Max(math.ZeroSize),
		)
		ttarget := o.targets[child.Control]
		target := ttarget.area
		target.Min = gxui.WindowToChild(target.Min, o.outer)
		target.Max = gxui.WindowToChild(target.Max, o.outer)

		if ttarget.orientation == gxui.Vertical {
			// Swap major an minor axes.
			target.Min.X, target.Min.Y = target.Min.Y, target.Min.X
			target.Max.X, target.Max.Y = target.Max.Y, target.Max.X
			bounds.Min.X, bounds.Min.Y = bounds.Min.Y, bounds.Min.X
			bounds.Max.X, bounds.Max.Y = bounds.Max.Y, bounds.Max.X
			cs.W, cs.H = cs.H, cs.W
		}

		var off math.Point
		// Minor; constrain within boundary.
		off.Y = target.Min.Y
		if off.Y > bounds.Max.Y-cs.H {
			off.Y = bounds.Max.Y - cs.H
		} else if off.Y < bounds.Min.Y {
			off.Y = bounds.Min.Y
		}
		// Major; avoid overlapping target if possible.
		off.X = target.Max.X
		if off.X+cs.W <= bounds.Max.X {
			goto layout
		}
		off.X = target.Min.X - cs.W
		if off.X >= bounds.Min.X {
			goto layout
		}
		off.X = target.Max.X
		if off.X > bounds.Max.X-cs.W {
			off.X = bounds.Max.X - cs.W
		} else if off.X < bounds.Min.X {
			off.X = bounds.Min.X
		}

	layout:
		if ttarget.orientation == gxui.Vertical {
			// Swap axes back.
			cs.W, cs.H = cs.H, cs.W
			off.X, off.Y = off.Y, off.X
		}
		child.Layout(cs.Rect().Offset(off))
	}
}

func (o *MenuOverlay) AddMenu(child gxui.Control, target math.Rect, orientation gxui.Orientation) *gxui.Child {
	o.targets[child] = menuTarget{
		area:        target,
		orientation: orientation,
	}
	return o.outer.AddChild(child)
}

func (o *MenuOverlay) RemoveMenu(child gxui.Control) {
	delete(o.targets, child)
	o.outer.RemoveChild(child)
}

func (o *MenuOverlay) Clear() {
	o.targets = make(map[gxui.Control]menuTarget, len(o.targets))
	o.outer.RemoveAll()
}

func (o *MenuOverlay) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (o *MenuOverlay) Brush() gxui.Brush {
	return o.brush
}

func (o *MenuOverlay) SetBrush(brush gxui.Brush) {
	if o.brush != brush {
		o.brush = brush
		o.Redraw()
	}
}

func (o *MenuOverlay) Pen() gxui.Pen {
	return o.pen
}

func (o *MenuOverlay) SetPen(pen gxui.Pen) {
	if o.pen != pen {
		o.pen = pen
		o.Redraw()
	}
}

func (o *MenuOverlay) Paint(c gxui.Canvas) {
	if !o.IsVisible() {
		return
	}
	for i, v := range o.outer.Children() {
		if v.Control.IsVisible() {
			c.Push()
			c.DrawRoundedRect(v.Bounds().Expand(math.Spacing{3, 3, 3, 3}), 3, 3, 3, 3, o.pen, o.brush)
			c.AddClip(v.Control.Size().Rect().Offset(v.Offset))
			o.outer.PaintChild(c, v, i)
			c.Pop()
		}
	}
}

func (o *MenuOverlay) MouseDown(ev gxui.MouseEvent) {
	o.InputEventHandler.MouseDown(ev)
	if !o.Container.ContainsPoint(ev.Point) {
		o.Clear()
	}
}

func (o *MenuOverlay) ContainsPoint(p math.Point) bool {
	if len(o.Children()) > 0 {
		return o.IsVisible() && o.Size().Rect().Contains(p)
	}
	return false
}
