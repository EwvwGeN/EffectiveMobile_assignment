package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("Usage: go run main.go <config_file_path>")
	}

	data, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatal("cant read config file")
	}

	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("cant unmarshal config file")
	}
	RewriteNestedKeys(config)
	DeleteMaps(config)
	envFile, err := os.Create(".env")
	if err != nil {
		log.Fatal("cant create .env file")
	}
	defer envFile.Close()
	var keys []string
	for k, _ := range config {
		keys = append(keys, k.(string))
	}
	sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
	})
	for _, key := range keys {
		fmt.Fprintf(envFile, "%s=%v\n", strings.ToUpper(key), config[key])
	}
}

func RewriteNestedKeys(part map[interface{}]interface{}) {
	for i, v := range part {
		switch v := v.(type) {
		case map[interface{}]interface{}:
			RewriteNestedKeys(v)
			for ii, vv := range v {
				part[i.(string)+"_"+ii.(string)] = vv
			}
		}
	}
}

func DeleteMaps(config map[interface{}]interface{}) {
	for i, v := range config {
		switch v.(type) {
		case map[interface{}]interface{}:
			delete(config, i)
		}
	}
}