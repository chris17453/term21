package core

import (
	"encoding/json"
	"fmt"
	"log"
)

func Print(data interface{}) {
	JSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("\n%s\n", string(JSON))

}
