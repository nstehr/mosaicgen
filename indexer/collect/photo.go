package collect

import (
	"image/color"
)

type Photo struct {
	Url      string
	ThumbUrl string
	Source   string
	Text     string
	Tag      string
	AvgColor color.RGBA
}
