package cast

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Env struct {
	Term  string `json:"TERM"`
	Shell string `json:"SHELL"`
}

type CastHeader struct {
	Width     int64   `json:"width"`
	Version   int64   `json:"version"`
	Title     string  `json:"title"`
	Height    int64   `json:"height"`
	Timestamp int64   `json:"timestamp"`
	Env       Env     `json:"env"`
	Duration  float64 `json:"duration"`
	Command   string  `json:"command"`
}

type CastData struct {
	Time float64
	Mode string
	Data string
}

type CastStream struct {
	File   string
	Header CastHeader
	Data   []CastData
}

// custom unmarshal function for this struc type
func (d *CastData) UnmarshalJSON(data []byte) error {

	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.Time = v[0].(float64)
	d.Mode = v[1].(string)
	d.Data = v[2].(string)
	return nil
}

// took 3 hours to fgigure this out. omg go go go.. why...
func (cs *CastStream) Load(filename string) error {
	cs.File = filename
	stream, err := os.Open(filename)
	defer stream.Close()
	if err != nil {
		return errors.New("CastStream file could not be opened")
	}
	parser := json.NewDecoder(stream)
	parser.Decode(&cs.Header)
	for {
		var cd CastData
		if err := parser.Decode(&cd); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		cs.Data = append(cs.Data, cd)
	}
	return nil
}
