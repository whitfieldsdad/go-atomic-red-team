package cmd

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func PrintJson(v interface{}) {
	blob, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(blob))
}

func PrintYaml(v interface{}) {
	blob, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(blob))
}
