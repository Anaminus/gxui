// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/anaminus/gxui"
	"github.com/anaminus/gxui/math"
	"github.com/anaminus/gxui/mixins/outer"
	"sort"
)

type LinearLayoutOuter interface {
	gxui.Container
	outer.Sized
}

type LinearLayout struct {
	outer               LinearLayoutOuter
	direction           gxui.Direction
	sizeMode            gxui.SizeMode
	horizontalAlignment gxui.HorizontalAlignment
	verticalAlignment   gxui.VerticalAlignment
}

func (l *LinearLayout) Init(outer LinearLayoutOuter) {
	l.outer = outer
}

type childSize struct {
	child *gxui.Child
	major int
	minor int
}

type childSizes []*childSize

func (c childSizes) Len() int {
	return len(c)
}
func (c childSizes) Less(i, j int) bool {
	return c[j].major < c[i].major
}
func (c childSizes) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (l *LinearLayout) LayoutChildren() {
	children := l.outer.Children()
	if len(children) == 0 {
		return
	}

	s := l.outer.Size().Contract(l.outer.Padding()).Max(math.ZeroSize)
	o := l.outer.Padding().LT()
	sizes := make(childSizes, len(children))
	var parentSize, contentSize, major int
	if l.direction.Orientation().Horizontal() {
		parentSize = s.W
		for i, c := range children {
			size := c.Control.DesiredSize(math.ZeroSize, s.Contract(c.Control.Margin()).Max(math.ZeroSize)).Expand(c.Control.Margin())
			sizes[i] = &childSize{
				child: c,
				major: size.W,
				minor: size.H,
			}
			contentSize += size.W
		}
		if l.direction.RightToLeft() {
			switch l.horizontalAlignment {
			case gxui.AlignLeft:
				if contentSize > parentSize {
					major = parentSize
				} else {
					major = contentSize
				}
			case gxui.AlignCenter:
				major = (parentSize-contentSize)/2 + contentSize
				if major < 0 {
					major = 0
				}
			case gxui.AlignRight:
				major = parentSize
			}
		} else {
			switch l.horizontalAlignment {
			case gxui.AlignLeft:
				major = 0
			case gxui.AlignCenter:
				major = (parentSize - contentSize) / 2
				if major < 0 {
					major = 0
				}
			case gxui.AlignRight:
				major = parentSize - contentSize
				if major < 0 {
					major = 0
				}
			}
		}
	} else {
		parentSize = s.H
		for i, c := range children {
			size := c.Control.DesiredSize(math.ZeroSize, s.Contract(c.Control.Margin()).Max(math.ZeroSize)).Expand(c.Control.Margin())
			sizes[i] = &childSize{
				child: c,
				major: size.H,
				minor: size.W,
			}
			contentSize += size.H
		}
		if l.direction.BottomToTop() {
			switch l.verticalAlignment {
			case gxui.AlignTop:
				if contentSize > parentSize {
					major = parentSize
				} else {
					major = contentSize
				}
			case gxui.AlignMiddle:
				major = (parentSize-contentSize)/2 + contentSize
				if major < 0 {
					major = 0
				}
			case gxui.AlignBottom:
				major = parentSize
			}
		} else {
			switch l.verticalAlignment {
			case gxui.AlignTop:
				major = 0
			case gxui.AlignMiddle:
				major = (parentSize - contentSize) / 2
				if major < 0 {
					major = 0
				}
			case gxui.AlignBottom:
				major = parentSize - contentSize
				if major < 0 {
					major = 0
				}
			}
		}
	}

	// If the the combined size of the children is greater than the size of
	// the parent, reduce the sizes of the children until there is no more
	// overflow.
	overflow := contentSize - parentSize
	if overflow > 0 {
		sort.Sort(sizes)
		largest := make(childSizes, 0, len(sizes))
		for overflow > 0 {
			largest = append(largest[:0], sizes[0])
			maxSize := sizes[0].major
			if maxSize <= 0 {
				// Finish if all objects are too small.
				break
			}
			goalSize := 0
			// Populate array with any other objects that have the same size.
			for i := 1; i < len(sizes); i++ {
				c := sizes[i]
				if c.major < maxSize {
					goalSize = c.major
					break
				}
				largest = append(largest, c)
			}

			// Collapse largest objects evenly, but only down to the size of
			// the next largest. This will allow the next largest object to be
			// counted in the next iteration.
			dist := maxSize - goalSize
			if dist > overflow {
				dist = overflow
			}
			for _, c := range largest {
				c.major -= dist / len(largest)
			}
			rem := dist % len(largest)
			if rem > 0 {
				// Distribute remainder to objects lower in the list, so that
				// the list remains sorted.
				for i := len(largest) - 1; i >= 0; i-- {
					if rem <= 0 {
						break
					}
					largest[i].major -= 1
					rem -= 1
				}
			}
			overflow -= dist
		}
	}

	if l.direction.Orientation().Horizontal() {
		for _, c := range sizes {
			c.child.Control.SetSize(math.Size{c.major, c.minor}.Contract(c.child.Control.Margin()).Max(math.ZeroSize))
		}
	} else {
		for _, c := range sizes {
			c.child.Control.SetSize(math.Size{c.minor, c.major}.Contract(c.child.Control.Margin()).Max(math.ZeroSize))
		}
	}

	for _, c := range children {
		cm := c.Control.Margin()
		cs := c.Control.Size()

		// Calculate minor-axis alignment
		var minor int
		switch l.direction.Orientation() {
		case gxui.Horizontal:
			switch l.verticalAlignment {
			case gxui.AlignTop:
				minor = cm.T
			case gxui.AlignMiddle:
				minor = (s.H - cs.H) / 2
			case gxui.AlignBottom:
				minor = s.H - cs.H
			}
		case gxui.Vertical:
			switch l.horizontalAlignment {
			case gxui.AlignLeft:
				minor = cm.L
			case gxui.AlignCenter:
				minor = (s.W - cs.W) / 2
			case gxui.AlignRight:
				minor = s.W - cs.W
			}
		}

		// Peform layout
		switch l.direction {
		case gxui.LeftToRight:
			major += cm.L
			c.Offset = math.Point{X: major, Y: minor}.Add(o)
			major += cs.W
			major += cm.R
		case gxui.RightToLeft:
			major -= cm.R
			c.Offset = math.Point{X: major - cs.W, Y: minor}.Add(o)
			major -= cs.W
			major -= cm.L
		case gxui.TopToBottom:
			major += cm.T
			c.Offset = math.Point{X: minor, Y: major}.Add(o)
			major += cs.H
			major += cm.B
		case gxui.BottomToTop:
			major -= cm.B
			c.Offset = math.Point{X: minor, Y: major - cs.H}.Add(o)
			major -= cs.H
			major -= cm.T
		}
	}
}

func (l *LinearLayout) DesiredSize(min, max math.Size) math.Size {
	if l.sizeMode.Fill() {
		return max
	}

	bounds := min.Rect()
	children := l.outer.Children()

	horizontal := l.direction.Orientation().Horizontal()
	offset := math.Point{X: 0, Y: 0}
	for _, c := range children {
		cs := c.Control.DesiredSize(math.ZeroSize, max)
		cm := c.Control.Margin()
		cb := cs.Expand(cm).Rect().Offset(offset)
		if horizontal {
			offset.X += cb.W()
		} else {
			offset.Y += cb.H()
		}
		bounds = bounds.Union(cb)
	}

	return bounds.Size().Expand(l.outer.Padding()).Clamp(min, max)
}

func (l *LinearLayout) Direction() gxui.Direction {
	return l.direction
}

func (l *LinearLayout) SetDirection(d gxui.Direction) {
	if l.direction != d {
		l.direction = d
		l.outer.Relayout()
	}
}

func (l *LinearLayout) SizeMode() gxui.SizeMode {
	return l.sizeMode
}

func (l *LinearLayout) SetSizeMode(mode gxui.SizeMode) {
	if l.sizeMode != mode {
		l.sizeMode = mode
		l.outer.Relayout()
	}
}

func (l *LinearLayout) HorizontalAlignment() gxui.HorizontalAlignment {
	return l.horizontalAlignment
}

func (l *LinearLayout) SetHorizontalAlignment(alignment gxui.HorizontalAlignment) {
	if l.horizontalAlignment != alignment {
		l.horizontalAlignment = alignment
		l.outer.Relayout()
	}
}

func (l *LinearLayout) VerticalAlignment() gxui.VerticalAlignment {
	return l.verticalAlignment
}

func (l *LinearLayout) SetVerticalAlignment(alignment gxui.VerticalAlignment) {
	if l.verticalAlignment != alignment {
		l.verticalAlignment = alignment
		l.outer.Relayout()
	}
}
