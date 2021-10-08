package term

import (
	"fmt"
	"term21/cast"
	"term21/font"
)

type Cursor struct {
	X              int
	Y              int
	Saved_X        int
	Saved_Y        int
	Visible        bool
	Blink          bool
	Blink_Interval int
	Character      int
}

type Video_Buffer struct {
	Width  int
	Height int
	Memory []int
	Length int64
}

type Viewport struct {
	Top    int
	Bottom int
	Width  int
	Height int
}

type Pallette_State struct {
	Foreground         int
	Background         int
	Default_Foreground int
	Default_Background int
}

type Display_Flags struct {
	Mode          int
	Reverse_Video int
	Bold          int
	Text_Mode     int
	Autowrap      int
	Scroll        int
	Pending_Wrap  int
}

type Term struct {
	Stream   cast.CastStream
	Buffer   Video_Buffer
	Cursor   Cursor
	Viewport Viewport
	Palette  Pallette_State
	Flags    Display_Flags
	Font     font.Font
}

type Config struct {
	Cast_File string
	Font_File string
}

func (ds *Term) Init(config Config) error {
	ds.Stream.Load(config.Cast_File)
	var width = ds.Stream.Header.Width
	var height = ds.Stream.Header.Height
	ds.Buffer.Width = width
	ds.Buffer.Height = height
	ds.Viewport.Width = width
	ds.Viewport.Height = height
	ds.Buffer.Memory = make([]int, width*height*3)
	ds.Buffer.Length = int64(len(ds.Buffer.Memory))

	// reset buffer to FG,BG,CHAR
	// 3 bytes per x position
	for i := range ds.Buffer.Memory {
		ds.Buffer.Memory[i] = 0
	}

	if err := ds.Font.Load(config.Font_File); err != nil {
		fmt.Println(err)
		return err
	}

	ds.Clear(0, 0, 7)
	ds.Write_At("Sample text", 0, 0, 0, 0)
	ds.Write_At("Not text", 0, 1, 0, 0)
	ds.Print()

	return nil
}

//write character in screen buffer
func (term *Term) Char_At(c, x, y, bg, fg int) {
	var pos int
	pos = x + int(y)*term.Buffer.Width

	term.Buffer.Memory[pos] = c
	term.Buffer.Memory[pos+1] = bg
	term.Buffer.Memory[pos+2] = fg
}

// write text in screen buffer
func (term *Term) Write_At(text string, x, y, bg, fg int) {
	var pos int

	for i := 0; i < len(text); i++ {
		pos = x + i + int(y)*term.Buffer.Width
		if pos >= int(term.Buffer.Length) {
			break
		}
		term.Buffer.Memory[pos] = int(text[i])
		term.Buffer.Memory[pos+1] = bg
		term.Buffer.Memory[pos+2] = fg
	}
}

// clear entire screen buffer
func (term *Term) Clear(c, bg, fg int) {
	var pos int
	for y := 0; y < term.Buffer.Height; y++ {
		for x := 0; x < term.Buffer.Width; x++ {
			pos = x + int(y)*term.Buffer.Width
			term.Buffer.Memory[pos] = c
			term.Buffer.Memory[pos+1] = bg
			term.Buffer.Memory[pos+2] = fg
		}
	}
}

// print the curent screen buffer, visible and non visible
func (term *Term) Print() {
	var pos, c int
	//var bg, fg int
	for y := 0; y < term.Buffer.Height; y++ {
		for x := 0; x < term.Buffer.Width; x++ {
			pos = x + int(y)*term.Buffer.Width
			c = term.Buffer.Memory[pos]
			//bg = term.Buffer.Memory[pos+1]
			//fg = term.Buffer.Memory[pos+2]
			fmt.Printf("%c", c)
		}
		fmt.Print("\n")
	}
}

func (ds *Term) Play(time float64) {

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
