package theme

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Theme struct {
	Directory          string
	Name               string    `yaml:"name"`
	Font_Face          string    `yaml:"font_face"`
	Font_Size          int64     `yaml:"font_size"`
	Default_Foreground int64     `yaml:"default_foreground"`
	Default_Background int64     `yaml:"default_background"`
	Foreground         int64     `yaml:"foreground"`
	Background         int64     `yaml:"background"`
	Margin             Margin    `yaml:"margin"`
	Texture            []Texture `yaml:"texture"`
	Text               []Text    `yaml:"text"`
	Pixel_Width        int64
	Pixel_Height       int64
	Terminal           Terminal
}
type Text struct {
	Name       string `yaml:"name"`
	Value      string `yaml:"value"`
	Foreground int    `yaml:"foreground"`
	Background int    `yaml:"background"`
	Column     int    `yaml:"col`
	Row        int    `yaml:"row`
	Font_Face  string `yaml:"font_face"`
	Font_Size  int64  `yaml:"font_size"`
}
type Terminal struct {
	Columns int64 `yaml:"columns"`
	Rows    int64 `yaml:"rows"`
	Width   int64 `yaml:"width"`
	Height  int64 `yaml:"height"`
}

type Margin struct {
	Left   int64 `yaml:"left"`
	Top    int64 `yaml:"top"`
	Right  int64 `yaml:"right"`
	Bottom int64 `yaml:"bottom"`
}
type Bounds struct {
	Left   int64 `yaml:"left"`
	Top    int64 `yaml:"top"`
	Right  int64 `yaml:"right"`
	Bottom int64 `yaml:"bottom"`
}
type Transparent struct {
	R     uint8 `yaml:"r"`
	G     uint8 `yaml:"g"`
	B     uint8 `yaml:"b"`
	A     uint8 `yaml:"a"`
	Index uint8 `yaml:"index"`
}

type Slice struct {
	Src       image.Rectangle
	Dst       image.Rectangle
	Copy_Mode int
}

type Sprite struct {
	Name         string      `yaml:"name"`
	Source_Inner Bounds      `yaml:"src_inner"`
	Source_Outer Bounds      `yaml:"src_outer"`
	Dest         Bounds      `yaml:"dest"`
	Rotation     int64       `yaml:"rotation"`
	CopyMode     string      `yaml:"copy_mode"`
	Transparent  Transparent `yaml:"transparent"`
	Slices       []Slice
}

type Texture struct {
	Image  image.Image
	Source string   `yaml:"source"`
	Sprite []Sprite `yaml:"sprite"`
}

type copy_mode int64

const (
	Copy_Block = 0
	Copy_Scale = 1
	Copy_Tile  = 2
)

func (s *Sprite) transform() {
	var new_slice Slice

	switch s.CopyMode {
	case "block":
		var src = image.Rectangle{image.Point{int(s.Source_Outer.Left), int(s.Source_Outer.Top)},
			image.Point{int(s.Source_Outer.Right), int(s.Source_Outer.Bottom)}}
		var dst = image.Rectangle{image.Point{int(s.Dest.Left), int(s.Dest.Top)},
			image.Point{int(s.Dest.Right), int(s.Dest.Bottom)}}

		new_slice = Slice{src, dst, Copy_Block}
		s.Slices = append(s.Slices, new_slice)
		break

	case "slice_9_scale":
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(1), s.dst_rect_slice_9(1), Copy_Block})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(2), s.dst_rect_slice_9(2), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(3), s.dst_rect_slice_9(3), Copy_Block})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(4), s.dst_rect_slice_9(4), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(5), s.dst_rect_slice_9(5), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(6), s.dst_rect_slice_9(6), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(7), s.dst_rect_slice_9(7), Copy_Block})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(8), s.dst_rect_slice_9(8), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(9), s.dst_rect_slice_9(9), Copy_Block})

		break
	case "slice_9_tile":
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(1), s.dst_rect_slice_9(1), Copy_Block})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(2), s.dst_rect_slice_9(2), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(3), s.dst_rect_slice_9(3), Copy_Block})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(4), s.dst_rect_slice_9(4), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(5), s.dst_rect_slice_9(5), Copy_Scale})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(6), s.dst_rect_slice_9(6), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(7), s.dst_rect_slice_9(7), Copy_Block})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(8), s.dst_rect_slice_9(8), Copy_Tile})
		s.Slices = append(s.Slices, Slice{s.src_rect_slice_9(9), s.dst_rect_slice_9(9), Copy_Block})

		break
	}
}

func (s *Sprite) dst_rect_slice_9(p int) image.Rectangle {
	var p1, p2 image.Point
	var x1, x2, x3, x4, x5, x6 int
	var y1, y2, y3, y4, y5, y6 int
	var width, height int
	var x_spacer, y_spacer int

	x1 = int(-s.Source_Outer.Left + s.Source_Outer.Left)
	x2 = int(-s.Source_Outer.Left + s.Source_Inner.Left - 1)
	x3 = int(-s.Source_Outer.Left + s.Source_Inner.Left)
	x4 = int(-s.Source_Outer.Left + s.Source_Inner.Right)
	x5 = int(-s.Source_Outer.Left + s.Source_Inner.Right + 1)
	x6 = int(-s.Source_Outer.Left + s.Source_Outer.Right)

	y1 = int(-s.Source_Outer.Top + s.Source_Outer.Top)
	y2 = int(-s.Source_Outer.Top + s.Source_Inner.Top - 1)
	y3 = int(-s.Source_Outer.Top + s.Source_Inner.Top)
	y4 = int(-s.Source_Outer.Top + s.Source_Inner.Bottom)
	y5 = int(-s.Source_Outer.Top + s.Source_Inner.Bottom + 1)
	y6 = int(-s.Source_Outer.Top + s.Source_Outer.Bottom)

	width = int(s.Dest.Right - s.Dest.Left)
	height = int(s.Dest.Bottom - s.Dest.Top)
	x_spacer = width - x6
	y_spacer = height - y6

	x1 += int(s.Dest.Left)
	x2 += int(s.Dest.Left)
	x3 += int(s.Dest.Left)
	x4 += int(s.Dest.Left) + x_spacer
	x5 += int(s.Dest.Left) + x_spacer
	x6 += int(s.Dest.Left) + x_spacer

	y1 += int(s.Dest.Top)
	y2 += int(s.Dest.Top)
	y3 += int(s.Dest.Top)
	y4 += int(s.Dest.Top) + y_spacer
	y5 += int(s.Dest.Top) + y_spacer
	y6 += int(s.Dest.Top) + y_spacer

	switch p {
	case 1:
		p1 = image.Point{x1, y1}
		p2 = image.Point{x2, y2}
		break

	case 2:
		p1 = image.Point{x3, y1}
		p2 = image.Point{x4, y2}
		break

	case 3:
		p1 = image.Point{x5, y1}
		p2 = image.Point{x6, y2}
		break

	case 4:
		p1 = image.Point{x1, y3}
		p2 = image.Point{x2, y4}
		break

	case 5:
		p1 = image.Point{x3, y3}
		p2 = image.Point{x4, y4}
		break

	case 6:
		p1 = image.Point{x5, y3}
		p2 = image.Point{x6, y4}
		break

	case 7:
		p1 = image.Point{x1, y5}
		p2 = image.Point{x2, y6}
		break

	case 8:
		p1 = image.Point{x3, y5}
		p2 = image.Point{x4, y6}
		break

	case 9:
		p1 = image.Point{x5, y5}
		p2 = image.Point{x6, y6}
		break
	}
	return image.Rectangle{p1, p2}

}

func (s *Sprite) src_rect_slice_9(p int) image.Rectangle {
	var p1, p2 image.Point
	var x1, x2, x3, x4, x5, x6 int
	var y1, y2, y3, y4, y5, y6 int

	x1 = int(s.Source_Outer.Left)
	x2 = int(s.Source_Inner.Left - 1)
	x3 = int(s.Source_Inner.Left)
	x4 = int(s.Source_Inner.Right)
	x5 = int(s.Source_Inner.Right + 1)
	x6 = int(s.Source_Outer.Right)

	y1 = int(s.Source_Outer.Top)
	y2 = int(s.Source_Inner.Top - 1)
	y3 = int(s.Source_Inner.Top)
	y4 = int(s.Source_Inner.Bottom)
	y5 = int(s.Source_Inner.Bottom + 1)
	y6 = int(s.Source_Outer.Bottom)
	switch p {

	case 1:
		p1 = image.Point{x1, y1}
		p2 = image.Point{x2, y2}
		break

	case 2:
		p1 = image.Point{x3, y1}
		p2 = image.Point{x4, y2}
		break

	case 3:
		p1 = image.Point{x5, y1}
		p2 = image.Point{x6, y2}
		break

	case 4:
		p1 = image.Point{x1, y3}
		p2 = image.Point{x2, y4}
		break

	case 5:
		p1 = image.Point{x3, y3}
		p2 = image.Point{x4, y4}
		break

	case 6:
		p1 = image.Point{x5, y3}
		p2 = image.Point{x6, y4}
		break

	case 7:
		p1 = image.Point{x1, y5}
		p2 = image.Point{x2, y6}
		break

	case 8:
		p1 = image.Point{x3, y5}
		p2 = image.Point{x4, y6}
		break

	case 9:
		p1 = image.Point{x5, y5}
		p2 = image.Point{x6, y6}
		break
	}
	return image.Rectangle{p1, p2}

}

func (t *Theme) get_directory(theme_name string) error {
	var home_dir = os.Getenv("HOME")
	var term21_dir = os.Getenv("TERM21_PATH")

	//check the term21 defined path
	var theme_dir string
	if term21_dir != "" {
		theme_dir = filepath.Join(term21_dir, "themes", theme_name)
		if _, err := os.Stat(theme_dir); !os.IsNotExist(err) {
			t.Directory = term21_dir
			return nil
		}
	}
	//core.Print(theme_dir)
	// check the user config path
	theme_dir = filepath.Join(home_dir, ".config", "term21", "themes", theme_name)
	if _, err := os.Stat(theme_dir); !os.IsNotExist(err) {
		t.Directory = filepath.Join(home_dir, ".config", "term21")
		return nil
	}
	//core.Print(theme_dir)
	// check the global config path
	theme_dir = filepath.Join("/etc", "term21", "themes", theme_name)
	if _, err := os.Stat(theme_dir); !os.IsNotExist(err) {
		t.Directory = filepath.Join("/etc", "term21")
		return nil
	}
	//core.Print(theme_dir)
	return errors.New("Term21 Theme not found")
}

func (c *Theme) Load(theme_name string) error {

	err := c.get_directory(theme_name)

	if err != nil {
		return err
	}
	file := filepath.Join(c.Directory, "themes", theme_name, "theme.t21")
	//fmt.Println("File:" + file)

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}
	return nil

}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func (tx *Texture) load_Image(directory string) error {

	file_path := filepath.Join(directory, tx.Source)
	existingImageFile, err := os.Open(file_path)
	if err != nil {

		return err
	}

	defer existingImageFile.Close()
	var imageType string
	tx.Image, imageType, err = image.Decode(existingImageFile)
	if err != nil {
		return err
	}
	//fmt.Println("IMAGE TYPE")

	//fmt.Println(imageType)
	existingImageFile.Seek(0, 0)

	switch imageType {
	case "png":
		tx.Image, err = png.Decode(existingImageFile)
		break
	case "jpg":
		tx.Image, err = jpeg.Decode(existingImageFile)
		break
	case "gif":
		tx.Image, err = gif.Decode(existingImageFile)
		break
	default:
		return errors.New("Unknown image format")
	}

	if err != nil {
		fmt.Println("Cannot load image")
		return err
	}
	return nil
}

func (r *Bounds) adjust(width, height int64) {
	if r.Left < 0 {
		r.Left += width
	}
	if r.Top < 0 {
		r.Top += height
	}
	if r.Right < 0 {
		r.Right += width
	}
	if r.Bottom < 0 {
		r.Bottom += height
	}
	///fmt.Println(r)
}

// the image which is the base for all renderin is bases on the max bounds of all bounding layers
func (t *Theme) Init(columns, rows, width, height int64) {
	t.Terminal.Columns = columns
	t.Terminal.Rows = rows
	t.Terminal.Width = width
	t.Terminal.Height = height

	var max_x, max_y int64
	var min_x, min_y int64
	max_x = t.Margin.Left + t.Margin.Right + t.Terminal.Width
	max_y = t.Margin.Top + t.Margin.Bottom + t.Terminal.Height

	for _, tx := range t.Texture {
		for _, s := range tx.Sprite {

			if s.Dest.Right > max_x {
				max_x = s.Dest.Right
			}

			if s.Dest.Bottom > max_y {
				max_y = s.Dest.Bottom
			}

			if s.Dest.Left < min_x {
				min_x = s.Dest.Left
			}

			if s.Dest.Top < min_y {
				min_y = s.Dest.Top
			}

		}
		if min_x < 0 {
			max_x += Abs(min_x)
			min_y = 0
		}
		if min_y < 0 {
			max_y += Abs(min_y)
			min_y = 0
		}

	}
	t.Pixel_Width = max_x
	t.Pixel_Height = max_y
	//fmt.Println(max_x, max_y)
	for t_key, _ := range t.Texture {
		for s_key, _ := range t.Texture[t_key].Sprite {

			t.Texture[t_key].Sprite[s_key].Dest.adjust(t.Pixel_Width, t.Pixel_Height)
			t.Texture[t_key].Sprite[s_key].Source_Inner.adjust(t.Pixel_Width, t.Pixel_Height)
			t.Texture[t_key].Sprite[s_key].Source_Outer.adjust(t.Pixel_Width, t.Pixel_Height)
			t.Texture[t_key].Sprite[s_key].transform()

		}
	}

	// load texture images
	for key, _ := range t.Texture {

		if err := t.Texture[key].load_Image(filepath.Join(t.Directory, "themes", t.Name)); err != nil {
			fmt.Println("ERR")
			fmt.Println(err)
		}
	}
}
