package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
)

func PrintJson(v interface{}) {
	blob, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("JSON marshalling failed - cowardly refusing to continue: %s - %s\n", blob, err)
	}
	fmt.Println(string(blob))
}