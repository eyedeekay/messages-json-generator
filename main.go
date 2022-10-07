package main

import (
	"encoding/json"
	"io"
	"os"
	"sort"

	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/msoap/html2data"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/html"
)

func readHtmlFromFile(fileName string) (string, error) {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func unSpace(unclean string) string {
	i := 0
	for {
		unclean = strings.ReplaceAll(unclean, "  ", " ")
		i++
		if i == 100 {
			break
		}
	}
	return unclean
}

func Parse(file string) map[string]map[string]string {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(file)
		panic(err)
	}
	l := html.NewLexer(parse.NewInputBytes(text))
	htmlData := html2data.FromFile(file)
	var retData = make(map[string]map[string]string)
	for {
		//tt, data := l.Next()
		tt, _ := l.Next()
		switch tt {
		case html.ErrorToken:
			if l.Err() != io.EOF {
				fmt.Println("Error on line", ":", l.Err())
			}
			return retData
		case html.StartTagToken:
			for {
				ttAttr, dataAttr := l.Next()
				if ttAttr != html.AttributeToken {
					break
				}
				peek := string(dataAttr)
				split := strings.SplitN(peek, "=", 2)
				if strings.Contains(string(dataAttr), ".png") || strings.Contains(string(dataAttr), ".js") || strings.Contains(string(dataAttr), "src=") {
					break
				} else {
					if len(split) == 2 {
						value := strings.ReplaceAll(strings.ReplaceAll(split[0], " ", ""), "\"", "")
						key := strings.ReplaceAll(strings.ReplaceAll(split[1], " ", ""), "\"", "")
						if value == "id" {
							var selector = make(map[string]string)
							selector["ids"] = "#" + key + ":not(:has(img))"
							value2, err := htmlData.GetData(selector)
							if err != nil {
								panic(err)
							}
							log.Println(value2)
							if len(value2) == 1 {
								for _, v := range value2 {
									if len(strings.Join(v, "")) > 0 {
										if !strings.Contains(strings.Join(v, ""), "\n") {
											message := unSpace(strings.ReplaceAll(strings.Join(v, " "), "\n", " "))
											description := message
											unit := make(map[string]string)
											unit["message"] = message
											unit["description"] = description
											retData[key] = unit
											break
										}
									}
								}

							}
						}
					}
				}
			}
		}
	}
}

func combine(parent, input map[string]map[string]string) map[string]map[string]string {
	for k, v := range input {
		parent[k] = v
	}
	return parent
}

func main() {
	var finalMap = make(map[string]map[string]string)
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".html") {
				fmt.Println(path, info.Size())
				finalMap = combine(finalMap, Parse(path))
				jsonStr, err := json.Marshal(finalMap)
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile("messages.ex.json", jsonStr, 0644)
				if err != nil {
					log.Fatal(err)
				}
				js := javascript(finalMap)
				if err := ioutil.WriteFile("messages.js", []byte(js), 0644); err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func javascript(finalMap map[string]map[string]string) string {
	javascript := `
function contentUpdateById(id, message) {
	let infoTitle = document.getElementById(id);
	let messageContent = chrome.i18n.getMessage(message);
	if (infoTitle === null) {
		console.log('content error', id, messageContent);
		return;
	}
	infoTitle.textContent = messageContent;
}
`
	keys := make([]string, 0)
	for k, _ := range finalMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		javascript += fmt.Sprintf("contentUpdateById('%s', '%s');\n", k, k)
	}
	return javascript
}
