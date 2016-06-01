package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mozilla-services/go-sops/sops"

	"gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go-sops <inFile>")
		os.Exit(1)
	}
	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	encYamlMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(fileBytes, encYamlMap)
	if err != nil {
		log.Fatal(err)
	}

	sopsBytes, err := yaml.Marshal(encYamlMap["sops"])
	if err != nil {
		log.Fatal(err)
	}

	sopsData, err := sops.NewData(sopsBytes)
	if err != nil {
		log.Fatal(err)
	}

	orderedMap := make(yaml.MapSlice, 0)
	err = yaml.Unmarshal(fileBytes, &orderedMap)

	decOrderedMap := sopsData.DecryptMapSlice(orderedMap, "")
	out, err := yaml.Marshal(decOrderedMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(out))
}