package basic

import (
	"github.com/anaminus/gxui"
	"github.com/anaminus/gxui/mixins"
)

func CreateTableLayout(theme *Theme) gxui.TableLayout {
	l := &mixins.TableLayout{}
	l.Init(l, theme)
	return l
}
