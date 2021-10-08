package font

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Rune struct {
	Width  int
	Height int
	Number int
	Data   []int
}

type Font struct {
	File       string
	Name       string
	Copyright  string
	Pointsize  int
	Height     int
	Ascent     int
	Inleading  int
	Exleading  int
	Italic     bool
	Underline  bool
	Strikeout  bool
	Weight     int
	Charset    int
	Characters []Rune
}

func get_hash(text string) (string, string) {
	tokens := strings.SplitN(text, " ", 2)
	var key, value string
	value = text

	if len(tokens) > 1 {
		key = strings.ToUpper(tokens[0])
		value = tokens[1]
	} else {
		key = ""
	}

	return key, value
}

func bool_str(value string) (bool, error) {

	value = strings.ToUpper(value)
	if value == "NO" || value == "FALSE" || value == "0" {
		return false, nil
	}
	if value == "YES" || value == "TRUE" || value == "1" {
		return true, nil
	}

	return false, errors.New("Cannot convert, Not a boolean value")
}

func (F *Font) header(key, value string) error {
	var err error
	switch key {
	case "FACENAME":
		F.Name = value
		break
	case "COPYRIGHT":
		F.Copyright = value
		break
	case "POINTSIZE":
		F.Pointsize, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("Cannot convert pointsize, not an integer")
		}
		break
	case "HEIGHT":
		F.Height, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("Cannot convert height, not an integer")
		}
		break
	case "ASCENT":
		F.Ascent, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("Cannot convert ascent, not an integer")
		}
		break
	case "INLEADING":
		F.Inleading, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("Cannot convert inleading, not an integer")
		}
		break
	case "EXLEADING":
		F.Exleading, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("Cannot convert exleading, not an integer")
		}
		break
	case "ITALIC":
		F.Italic, err = bool_str(value)
		if err != nil {
			return errors.New("Cannot convert italic, not a bool")
		}
		break
	case "UNDERLINE":
		F.Underline, err = bool_str(value)
		if err != nil {
			return errors.New("Cannot convert underline, not a bool")
		}
		break
	case "STRIKEOUT":
		F.Strikeout, err = bool_str(value)
		if err != nil {
			return errors.New("Cannot convert strikeout, not a bool")
		}
		break
	case "WEIGHT":
		F.Weight, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("Cannot convert weight, not an integer")
		}
		break
	case "CHARSET":
		F.Charset, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("Cannot convert charset, not an integer")
		}
		F.Charset++
		break
	}
	return nil
}

func (font *Font) Load(filename string) error {
	font.File = filename

	stream, err := os.Open(font.File)
	defer stream.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	scanner := bufio.NewScanner(stream)
	var Rune_Number = 0
	var in_header = true

	// walking the textfile line by line
	for scanner.Scan() {
		var line = scanner.Text()
		line = strings.TrimSpace(line)
		// ignore emnpty lines
		if len(line) == 0 {
			Rune_Number = -1
			continue
		}
		// ignore comments
		if line[0] == '#' {
			//fmt.Println("Comment ->" + line)
			continue
		}

		//get the next key/value combo
		key, value := get_hash(line)

		// preform in header text matching
		if in_header == true {
			// load the header values
			font.header(key, value)
			// charset is the last header value and is required.
			if font.Charset > 0 {
				in_header = false
				// resize character set
				for i := 0; i < font.Charset; i++ {
					var nr Rune
					nr.Number = i
					font.Characters = append(font.Characters, nr)
				}
			}
		} else {
			// header scan

			switch key {
			case "CHAR":
				Rune_Number, err = strconv.Atoi(value)
				if err != nil {
					return errors.New("Cannot convert character index, not an integer")
				}
				if Rune_Number < 0 || Rune_Number >= font.Charset {
					return errors.New("Character outside of index")
				}
				break
			case "WIDTH":
				font.Characters[Rune_Number].Width, err = strconv.Atoi(value)
				if err != nil {
					return errors.New("Cannot convert width index, not an integer")
				}
				break
			case "":
				var data_len = len(value)
				var data_int int
				if data_len == 0 {
					continue
				}

				for pos, n := range value {
					if pos >= font.Characters[Rune_Number].Width {

						return errors.New("font character is wider than defined")
					}
					if n == '.' {
						data_int = 0
					}
					if n == 'x' {
						data_int = 1
					}
					font.Characters[Rune_Number].Data = append(font.Characters[Rune_Number].Data, data_int)

					if err != nil {
						return errors.New("Cannot convert data index, not an integer")
					}
				}
				font.Characters[Rune_Number].Height++

				break

			} // end switch

			// we are in the character/rune section

		} //end else
	} //end for
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}
