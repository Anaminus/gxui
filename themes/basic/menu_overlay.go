package basic

import (
	"github.com/anaminus/gxui"
	"github.com/anaminus/gxui/math"
	"github.com/anaminus/gxui/mixins"
)

type MenuOverlay struct {
	mixins.MenuOverlay
	theme *Theme
}

func CreateMenuOverlay(theme *Theme) gxui.MenuOverlay {
	b := &MenuOverlay{}
	b.Init(b, theme)
	b.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	b.SetPadding(math.Spacing{L: 5, T: 5, R: 5, B: 5})
	b.SetPen(theme.MenuOverlayStyle.Pen)
	b.SetBrush(theme.MenuOverlayStyle.Brush)
	b.theme = theme
	return b
}
