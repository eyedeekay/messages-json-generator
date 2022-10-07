package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

func combine(parent, input map[string]map[string]string) map[string]map[string]string {
	for k, v := range input {
		parent[k] = v
	}
	return parent
}

func main() {
	a := flag.String("file-a", "", "base .json file")
	b := flag.String("file-b", "", ".json file to add to file-a")
	flag.Parse()
	content, err := ioutil.ReadFile(*a)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload map[string]map[string]string
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	contentb, err := ioutil.ReadFile(*b)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payloadb map[string]map[string]string
	err = json.Unmarshal(contentb, &payloadb)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	val := combine(payload, payloadb)
	log.Println(val)
	bytes, err := json.Marshal(val)
	if err != nil {
		log.Fatal("Error during Marshal(): ", err)
	}
	if err := ioutil.WriteFile("messages.combined.json", bytes, 0644); err != nil {
		log.Fatal("Error during Write(): ", err)
	}

}
