package raster

import (
	"image"
	"image/color"
	"term21/core"
	"term21/theme"
)

func Draw_block(img *image.Paletted, src_img image.Image, r theme.Slice, transparency theme.Transparent) {
	var src_width, src_height int
	var x2, y2 int

	src_width = r.Src.Max.X - r.Src.Min.X + 1
	src_height = r.Src.Max.Y - r.Src.Min.Y + 1

	y2 = 0
	t_color := color.RGBA{transparency.R, transparency.G, transparency.B, transparency.A}
	for y := r.Dst.Min.Y; y <= r.Dst.Max.Y; y++ {
		if y2 >= src_height {
			break
		}
		x2 = 0
		for x := r.Dst.Min.X; x <= r.Dst.Max.X; x++ {
			if x2 >= src_width {
				continue
			}
			c := src_img.At(x2+r.Src.Min.X, y2+r.Src.Min.Y)

			if c == t_color {
				x2++
				continue
			}
			img.Set(x, y, c)
			x2++
		}
		y2++
		y2 %= src_height
	}
}

func Draw_tile(img *image.Paletted, src_img image.Image, r theme.Slice, transparency theme.Transparent) {
	var src_width, src_height int
	var x2, y2 int

	src_width = r.Src.Max.X - r.Src.Min.X + 1
	src_height = r.Src.Max.Y - r.Src.Min.Y + 1

	y2 = 0
	t_color := color.RGBA{transparency.R, transparency.G, transparency.B, transparency.A}
	for y := r.Dst.Min.Y; y <= r.Dst.Max.Y; y++ {
		x2 = 0
		for x := r.Dst.Min.X; x <= r.Dst.Max.X; x++ {
			x2 %= src_width
			c := src_img.At(x2+r.Src.Min.X, y2+r.Src.Min.Y)

			if c == t_color {
				x2++
				continue
			}
			img.Set(x, y, c)
			x2++
		}
		y2++
		y2 %= src_height
	}
}

func Draw_scale(img *image.Paletted, src_img image.Image, r theme.Slice, transparency theme.Transparent) {
	var src_width, src_height int
	var x2, y2 int

	src_width = r.Src.Max.X - r.Src.Min.X + 1
	src_height = r.Src.Max.Y - r.Src.Min.Y + 1

	t_color := color.RGBA{transparency.R, transparency.G, transparency.B, transparency.A}
	dest_width := r.Dst.Max.X - r.Dst.Min.X
	dest_height := r.Dst.Max.Y - r.Dst.Min.Y
	core.Print(dest_height)
	for y := 0; y <= dest_height; y++ {
		fy := float32(y) / float32(dest_height)
		for x := 0; x <= dest_width; x++ {
			fx := float32(x) / float32(dest_width)
			x2 = int(fx * float32(src_width))
			y2 = int(fy * float32(src_height))
			c := src_img.At(x2+r.Src.Min.X, y2+r.Src.Min.Y)

			if c == t_color {
				//continue
			}

			img.Set(x+r.Dst.Min.X, y+r.Dst.Min.Y, c)
		}
	}
}
