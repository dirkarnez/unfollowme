package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/jsonquery"
)

func main() {
	dat, err := os.ReadFile("seed.txt")
	checkErr(err)
	str := string(dat)

	re := regexp.MustCompile(`-H\s+'([^:]+):\s+([^']+)'`)
	matched := re.FindAllStringSubmatch(str, -1)
	//fmt.Println(len(matched))

	client := &http.Client{}

	var url = "https://i.instagram.com/api/v1/friendships/49612955961/followers/?count=12"

	type Entry struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	jsonEntry := []Entry{}

	var count = 0
	var continued = true
	for continued {
		req, err := http.NewRequest("GET", url, nil)
		checkErr(err)

		for _, header := range matched {
			req.Header[header[1]] = []string{header[2]}
		}

		resp, err := client.Do(req)
		checkErr(err)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		doc, err := jsonquery.Parse(strings.NewReader(string(body)))
		checkErr(err)

		list, err := jsonquery.QueryAll(doc, "/users/*")
		if err != nil {
			log.Println(err)
		}

		for _, value := range list {
			pk, err := jsonquery.Query(value, "/pk")
			if err != nil {
				log.Println(err)
			}

			username, err := jsonquery.Query(value, "/username")
			if err != nil {
				log.Println(err)
			}

			jsonEntry = append(jsonEntry, Entry{ID: fmt.Sprintf("%s", pk.Value()), Name: fmt.Sprintf("%s", username.Value())})
		}

		count = count + len(list)

		nextMaxId, err := getNextMaxId(doc)
		if err != nil || len(nextMaxId) == 0 {
			continued = false
		} else {
			url = fmt.Sprintf(`https://i.instagram.com/api/v1/friendships/49612955961/followers/?count=12&max_id=%s`, nextMaxId)
		}
	}

	bytes, err := json.MarshalIndent(jsonEntry, "", "\t")
	checkErr(err)

	err = ioutil.WriteFile(fmt.Sprintf("%s.txt", strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")), bytes, 0644)
	checkErr(err)

	fmt.Println("Done,", count)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getNextMaxId(doc *jsonquery.Node) (string, error) {
	next_max_id, err := jsonquery.Query(doc, "/next_max_id")
	if err != nil {
		return "", fmt.Errorf("next_max_id is nil #1")
	}

	if next_max_id == nil {
		return "", fmt.Errorf("next_max_id is nil #2")
	}

	val := next_max_id.Value()
	if val == nil {
		return "", fmt.Errorf("next_max_id is nil #3")
	}
	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("next_max_id not a string")
	}
	return s, nil
}

// func getHasMore(doc *jsonquery.Node) (bool, error) {
// 	hasMore, err := jsonquery.Query(doc, "/has_more")
// 	if err != nil {
// 		return false, fmt.Errorf("has_more is nil #1")
// 	}

// 	if hasMore == nil {
// 		return false, fmt.Errorf("has_more is nil #2")
// 	}

// 	val := hasMore.Value()
// 	if val == nil {
// 		return false, fmt.Errorf("has_more is nil #3")
// 	}
// 	s, ok := val.(bool)
// 	if !ok {
// 		return false, fmt.Errorf("has_more not a bool")
// 	}
// 	return s, nil
// }
