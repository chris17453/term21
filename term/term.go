package term

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"regexp"
	"term21/cast"
	"term21/core"
	"term21/font"
	"term21/raster"
	"term21/theme"

	"github.com/pborman/ansi"
)

type Cursor struct {
	X              int64
	Y              int64
	Saved_X        int64
	Saved_Y        int64
	Visible        bool
	Blink          bool
	Blink_Interval int64
	Character      int64
}

type Video_Buffer struct {
	Width   int64
	Height  int64
	Display []int
	FG      []uint8
	BG      []uint8
	Length  int64
}

type Palette_State struct {
	Foreground         int64
	Background         int64
	Default_Foreground int64
	Default_Background int64
}

type Display_Flags struct {
	Mode          int64
	Reverse_Video int64
	Bold          int64
	Text_Mode     int64
	Autowrap      int64
	Scroll        int64
	Pending_Wrap  int64
}

type Term struct {
	Stream  cast.CastStream
	Buffer  Video_Buffer
	Cursor  Cursor
	Palette Palette_State
	Flags   Display_Flags
	Font    font.Font
	Theme   theme.Theme
	Img     *image.Paletted
	Columns int64
	Rows    int64
}

type Config struct {
	Cast_File string
	Theme     string
}

func (ds *Term) Init(config Config) error {
	ds.Stream.Load(config.Cast_File)

	if err := ds.Theme.Load(config.Theme); err != nil {
		return err
	}

	if err := ds.Font.Load(filepath.Join(ds.Theme.Directory, "fonts", ds.Theme.Font_Face)); err != nil {
		//fmt.Println(err)
		return err
	}

	var width = ds.Stream.Header.Width
	var height = ds.Stream.Header.Height

	ds.Theme.Init(width,
		height,
		width*ds.Font.Width,
		height*ds.Font.Height)

	ds.Buffer.Width = width
	ds.Buffer.Height = height
	ds.Buffer.Display = make([]int, width*height)
	ds.Buffer.FG = make([]uint8, width*height)
	ds.Buffer.BG = make([]uint8, width*height)
	ds.Buffer.Length = int64(len(ds.Buffer.Display))

	ds.Clear(0, 0, 7)

	/*
		ds.Clear(0, 0, 7)
		ds.Write_At("Sample text", 0, 0, 0, 0)
		ds.Write_At("Not text", 3, 1, 0, 0)
	*/
	//core.Print(ds.Theme)
	//ds.Print()

	return nil
}

//write character in screen buffer
func (term *Term) Char_At(c int, x, y int64, bg, fg uint8) {
	var pos int64
	pos = x + y*term.Buffer.Width

	term.Buffer.Display[pos] = c
	term.Buffer.BG[pos] = bg
	term.Buffer.FG[pos] = fg
}

// write text in screen buffer
func (term *Term) Write_At(text string, x, y int64, bg, fg uint8) {
	var pos, i int64
	var text_len int64
	text_len = int64(len(text))

	for i = 0; i < text_len; i++ {
		pos = x + i + y*term.Buffer.Width
		if pos >= term.Buffer.Length {
			break
		}
		term.Buffer.Display[pos] = int(text[i])
		term.Buffer.BG[pos] = bg
		term.Buffer.FG[pos] = fg
	}
}

// clear entire screen buffer
func (term *Term) Clear(c int, bg, fg uint8) {
	var pos, x, y int64
	for y = 0; y < term.Buffer.Height; y++ {
		for x = 0; x < term.Buffer.Width; x++ {
			pos = x + y*term.Buffer.Width
			term.Buffer.Display[pos] = c
			term.Buffer.FG[pos] = bg
			term.Buffer.BG[pos] = fg
		}
	}
}

// print the curent screen buffer, visible and non visible
func (term *Term) Print() {
	var pos, x, y int64
	var c int
	//var bg, fg int

	for x = 0; x < term.Buffer.Width+2; x++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
	for y = 0; y < term.Buffer.Height; y++ {
		fmt.Print("|")
		for x = 0; x < term.Buffer.Width; x++ {
			pos = x + y*term.Buffer.Width
			c = term.Buffer.Display[pos]
			//bg = term.Buffer.Memory[pos+1]
			//fg = term.Buffer.Memory[pos+2]
			if c < 32 {
				c = 32
			}
			fmt.Printf("%c", c)
		}
		fmt.Print("|\n")
	}
	for x = 0; x < term.Buffer.Width+2; x++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
}

// print the curent screen buffer, visible and non visible
func (term *Term) Draw_Screen() {
	var pos, x, y int64
	var c int
	var bg, fg uint8

	fmt.Print("\n")
	for y = 0; y < term.Buffer.Height; y++ {
		for x = 0; x < term.Buffer.Width; x++ {
			pos = x + y*term.Buffer.Width
			c = term.Buffer.Display[pos]
			fg = term.Buffer.FG[pos]
			bg = term.Buffer.BG[pos]
			if c < 32 {
				continue
			}
			for _, char := range term.Font.Characters {
				if int(char.Number) == c {
					var fx, fy int64
					for fy = 0; fy < char.Height; fy++ {
						for fx = 0; fx < char.Width; fx++ {
							var img_pos = term.Theme.Margin.Left + x*term.Font.Width + fx + (y*term.Font.Height+fy+term.Theme.Margin.Left)*term.Theme.Pixel_Width

							if char.Data[fx+fy*char.Width] == 0 {
								term.Img.Pix[img_pos] = 0
								fg = fg + 1
							} else {
								bg = bg - 1
								term.Img.Pix[img_pos] = 15
							}
						}
					}
					break
				}
			}
		}
	}
}

func (ds *Term) Play(time float64) {

}

func XTerm_Palette() color.Palette {
	rgb := [768]uint8{
		0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 0,
		0, 0, 128, 128, 0, 128, 0, 128, 128, 192, 192, 192,
		128, 128, 128, 255, 0, 0, 0, 255, 0, 255, 255, 0,

		0, 0, 255, 255, 0, 255, 0, 255, 255, 255, 255, 255,
		// xterm palette
		0, 0, 0, 0, 0, 95, 0, 0, 135, 0, 0, 175, 0, 0, 215, 0, 0, 255,
		0, 95, 0, 0, 95, 95, 0, 95, 135, 0, 95, 175, 0, 95, 215, 0, 95, 255,
		0, 135, 0, 0, 135, 95, 0, 135, 135, 0, 135, 175, 0, 135, 215, 0, 135, 255,
		0, 175, 0, 0, 175, 95, 0, 175, 135, 0, 175, 175, 0, 175, 215, 0, 175, 255,
		0, 215, 0, 0, 215, 95, 0, 215, 135, 0, 215, 175, 0, 215, 215, 0, 215, 255,
		0, 255, 0, 0, 255, 95, 0, 255, 135, 0, 255, 175, 0, 255, 215, 0, 255, 255,
		95, 0, 0, 95, 0, 95, 95, 0, 135, 95, 0, 175, 95, 0, 215, 95, 0, 255,
		95, 95, 0, 95, 95, 95, 95, 95, 135, 95, 95, 175, 95, 95, 215, 95, 95, 255,
		95, 135, 0, 95, 135, 95, 95, 135, 135, 95, 135, 175, 95, 135, 215, 95, 135, 255,
		95, 175, 0, 95, 175, 95, 95, 175, 135, 95, 175, 175, 95, 175, 215, 95, 175, 255,
		95, 215, 0, 95, 215, 95, 95, 215, 135, 95, 215, 175, 95, 215, 215, 95, 215, 255,
		95, 255, 0, 95, 255, 95, 95, 255, 135, 95, 255, 175, 95, 255, 215, 95, 255, 255,
		135, 0, 0, 135, 0, 95, 135, 0, 135, 135, 0, 175, 135, 0, 215, 135, 0, 255,
		135, 95, 0, 135, 95, 95, 135, 95, 135, 135, 95, 175, 135, 95, 215, 135, 95, 255,
		135, 135, 0, 135, 135, 95, 135, 135, 135, 135, 135, 175, 135, 135, 215, 135, 135, 255,
		135, 175, 0, 135, 175, 95, 135, 175, 135, 135, 175, 175, 135, 175, 215, 135, 175, 255,
		135, 215, 0, 135, 215, 95, 135, 215, 135, 135, 215, 175, 135, 215, 215, 135, 215, 255,
		135, 255, 0, 135, 255, 95, 135, 255, 135, 135, 255, 175, 135, 255, 215, 135, 255, 255,
		175, 0, 0, 175, 0, 95, 175, 0, 135, 175, 0, 175, 175, 0, 215, 175, 0, 255,
		175, 95, 0, 175, 95, 95, 175, 95, 135, 175, 95, 175, 175, 95, 215, 175, 95, 255,
		175, 135, 0, 175, 135, 95, 175, 135, 135, 175, 135, 175, 175, 135, 215, 175, 135, 255,
		175, 175, 0, 175, 175, 95, 175, 175, 135, 175, 175, 175, 175, 175, 215, 175, 175, 255,
		175, 215, 0, 175, 215, 95, 175, 215, 135, 175, 215, 175, 175, 215, 215, 175, 215, 255,
		175, 255, 0, 175, 255, 95, 175, 255, 135, 175, 255, 175, 175, 255, 215, 175, 255, 255,
		215, 0, 0, 215, 0, 95, 215, 0, 135, 215, 0, 175, 215, 0, 215, 215, 0, 255,
		215, 95, 0, 215, 95, 95, 215, 95, 135, 215, 95, 175, 215, 95, 215, 215, 95, 255,
		215, 135, 0, 215, 135, 95, 215, 135, 135, 215, 135, 175, 215, 135, 215, 215, 135, 255,
		215, 175, 0, 215, 175, 95, 215, 175, 135, 215, 175, 175, 215, 175, 215, 215, 175, 255,
		215, 215, 0, 215, 215, 95, 215, 215, 135, 215, 215, 175, 215, 215, 215, 215, 215, 255,
		215, 255, 0, 215, 255, 95, 215, 255, 135, 215, 255, 175, 215, 255, 215, 215, 255, 255,
		255, 0, 0, 255, 0, 95, 255, 0, 135, 255, 0, 175, 255, 0, 215, 255, 0, 255,
		255, 95, 0, 255, 95, 95, 255, 95, 135, 255, 95, 175, 255, 95, 215, 255, 95, 255,
		255, 135, 0, 255, 135, 95, 255, 135, 135, 255, 135, 175, 255, 135, 215, 255, 135, 255,
		255, 175, 0, 255, 175, 95, 255, 175, 135, 255, 175, 175, 255, 175, 215, 255, 175, 255,
		255, 215, 0, 255, 215, 95, 255, 215, 135, 255, 215, 175, 255, 215, 215, 255, 215, 255,
		255, 255, 0, 255, 255, 95, 255, 255, 135, 255, 255, 175, 255, 255, 215, 255, 255, 255,
		8, 8, 8, 18, 18, 18, 28, 28, 28, 38, 38, 38, 48, 48, 48, 58, 58, 58, 68, 68, 68,
		78, 78, 78, 88, 88, 88, 98, 98, 98, 108, 108, 108, 118, 118, 118, 128, 128, 128,
		138, 138, 138, 148, 148, 148, 158, 158, 158, 168, 168, 168, 178, 178, 178, 188, 188, 188,
		198, 198, 198, 208, 208, 208, 218, 218, 218, 228, 228, 228, 238, 238, 238}
	var colors []color.Color

	var r, g, b, a uint8
	a = 0xFF
	for i := 0; i < 256; i++ {
		r = rgb[i*3+0]
		g = rgb[i*3+1]
		b = rgb[i*3+2]
		c := color.RGBA{r, g, b, a}
		colors = append(colors, c)
	}

	return color.Palette(colors)

}

func (term *Term) Term_Image() {
	var width, height int
	width = int(term.Theme.Pixel_Width)
	height = int(term.Theme.Pixel_Height)
	var palette = XTerm_Palette()
	rect := image.Rect(0, 0, width, height)
	img := image.NewPaletted(rect, palette)
	term.Img = img
}

func GetRxParams(rx *regexp.Regexp, str string) (pm map[string]string) {
	if !rx.MatchString(str) {
		return nil
	}
	p := rx.FindStringSubmatch(str)
	n := rx.SubexpNames()
	pm = map[string]string{}
	for i := range n {
		if i == 0 {
			continue
		}

		if n[i] != "" && p[i] != "" {
			pm[n[i]] = p[i]
		}
	}
	return
}

func (t *Term) GifStream() {
	var gif_file = "my.gif"

	// open output file
	fo, err := os.Create(gif_file)
	if err != nil {
		panic(err)
	}

	t.Term_Image()

	anim := gif.GIF{}
	Delay := 4
	//var pre string
	var data string
	data = ""
	for _, item := range t.Stream.Data {
		core.Print(item)
		// loop through all the loaded images
		for _, tx := range t.Theme.Texture {
			//loop through all the layers form this image
			for _, s := range tx.Sprite {
				// draw all of the sliced layers
				for _, sl := range s.Slices {
					switch sl.Copy_Mode {
					case theme.Copy_Block:
						raster.Draw_block(t.Img, tx.Image, sl, s.Transparent)
						break
					case theme.Copy_Tile:
						raster.Draw_tile(t.Img, tx.Image, sl, s.Transparent)
						break
					case theme.Copy_Scale:
						raster.Draw_scale(t.Img, tx.Image, sl, s.Transparent)
						break
					}
				}
			}
		}
		/*
			var ANSI_SINGLE, ANSI_CHAR_SET, ANSI_G0, ANSI_G1, BRACKET_PASTE, ANSI_TITLE, ANSI_OSC, ANSI_CSI_RE string
			//var    ANSI_CSI_RE,   string

			ANSI_SINGLE = "[\033](?P<SINGLE>[cDEHMZ6789>=i])" //ijkl arrow keys
			ANSI_CHAR_SET = "[\033]\\%([?P<CHAR>@G*])"
			ANSI_G0 = "[\033]\\((?P<G0>[B0UK])"
			ANSI_G1 = "[\033]\\)(?P<G1>[B0UK])"
			BRACKET_PASTE = "[\033]\\[(?P<BRACKET>20[0-1]~)"
			ANSI_OSC = "(?:\033\\]|\u009d)(?P<OSC>.*?)(?:\033\\\\|[\a\u009c])"
			ANSI_TITLE = "[\033][k](?P<TITLE>.*)[\033][\\\\]"
			ANSI_CSI_RE = "[\033]\\[(?P<CSI_CODE>(?:\\d|;|<|>|=|\\?)*)(?P<CSI_DATA>[a-zA-Z])\002?"

			ESC_SEQUENCES := []string{ANSI_SINGLE, ANSI_CHAR_SET, ANSI_G0, ANSI_G1, BRACKET_PASTE, ANSI_TITLE, ANSI_OSC, ANSI_CSI_RE, "(?P<TEXT>.?)*"}
			ANSI_REGEX := "" + strings.Join(ESC_SEQUENCES, "|") + ""
		*/

		data += item.Data

		t.Draw_Screen()

		t.Img.SetColorIndex(2, 2, 1)

		anim.Delay = append(anim.Delay, Delay)
		anim.Image = append(anim.Image, t.Img)
		//		break
	}

	//fmt.Print(data)

	b_data := []byte(data)
	for {

		out, s, err := ansi.Decode(b_data)
		if out == nil || err != nil {
			break
		}
		b_data = out
		if err != nil {
			fmt.Print("ADDED")
			break
		}
		//core.Print(remb)
		core.Print(s)
	}
	gif.EncodeAll(fo, &anim)
	if err := fo.Close(); err != nil {
		panic(err)
	}

}

/*
  ds_init                    (unsigned int width,unsigned int height,unsigned int foreground,unsigned int background)
  ds_text_mode_on            (ds *state)
  ds_text_mode_off           (ds *state)
  ds_autowrap_on             (ds *state)
  ds_autowrap_on             (ds *state)
  ds_set_scroll_region       (ds *state,unsigned int top,unsigned int bottom)
  ds_show_cursor             (ds *state)
  ds_hide_cursor             (ds *state)
  ds_check_bounds            (ds *state)
  ds_cursor_up               (ds *state,unsigned int distance)
  ds_cursor_down             (ds *state,unsigned int distance)
  ds_cursor_left             (ds *state,unsigned int distance)
  ds_cursor_right            (ds *state,unsigned int distance)
  ds_cursor_absolute_x       (ds *state,unsigned int x)
  ds_cursor_absolute_y       (ds *state,unsigned int y)
  ds_cursor_absolute         (ds *state,unsigned int x, unsigned int y)
  ds_cursor_save_position    (ds *state)
  ds_cursor_restore_position (ds *state)
  ds_cursor_get_position     (ds *state)
  ds_set_background          (ds *state,int color)
  ds_set_foreground          (ds *state,int color)
*/
