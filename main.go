package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	dcardAPIBase    string = "https://www.dcard.tw/_api/"
	dcardAPISexPost string = dcardAPIBase + "/forums/sex/posts?popular=false"
)

var (
	currForum string
)

func getArticle(firstID int, lastID int) (int, int) {
	var article Articles
	url := dcardAPISexPost

	if firstID != 0 {
		url = fmt.Sprintf("%s&after=%d", dcardAPISexPost, firstID)
	} else if lastID != 0 {
		url = fmt.Sprintf("%s&before=%d", dcardAPISexPost, lastID)
	}

	firstID = 0

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalln(err)
		log.Fatalln(resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal([]byte(body), &article); err != nil {
		log.Fatalln(err)
	}

	for _, v := range article {
		fmt.Printf("[%v](%v) -> %v: %v\n", v.ID, v.CreatedAt, len(v.Media), v.Title)

		if firstID == 0 {
			firstID = v.ID
		}

		lastID = v.ID
	}

	return firstID, lastID
}

func main() {
	currForum = "sex"
	log.Printf("Forum: %v, URL: %v\n", currForum, dcardAPISexPost)

	firstID, lastID := getArticle(0, 0)

	scanner := bufio.NewScanner(os.Stdin)
	quit := false

	for !quit {
		fmt.Println("n: Next, p: Previous, v: View, d: Download, q/quit/exit: Quit")
		fmt.Printf("Dcard:%v> ", currForum)

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		tokens := strings.Split(line, " ")
		cmd := tokens[0]
		args := tokens[1:]

		switch cmd {
		case "q", "quit", "exit":
			quit = true
		case "n":
			fmt.Println("Next Page")
			firstID, lastID = getArticle(0, lastID)
		case "p":
			fmt.Println("Previous Page")
			firstID, lastID = getArticle(firstID, 0)
		case "v":
			fmt.Println("View Article")
		case "d":
			if len(args) == 0 {
				fmt.Println("No article specified. Try input 'd 1' to get media file")
				continue
			}

			fmt.Printf("Download article media: %v\n", args)
		}
	}
}
