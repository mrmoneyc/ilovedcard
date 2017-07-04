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
	dcardAPIBase   string = "https://www.dcard.tw/_api/"
	dcardAPIForums string = dcardAPIBase + "forums"
)

var (
	dcardAPIPostMeta string
	currForum        string
)

func getForums() string {
	var forum Forums

	resp, err := http.Get(dcardAPIForums)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalln(err)
		log.Fatalln(resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal([]byte(body), &forum); err != nil {
		log.Fatalln(err)
	}

	for k, v := range forum {
		fmt.Printf("[%v] %v (%v)\n", k, v.Name, v.Alias)
	}

	fmt.Print("Enter forum ID: ")
	var forumID int
	fmt.Scan(&forumID)
	if forumID > len(forum) {
		log.Fatalln("Forum not found.")
	}

	dcardAPIPostMeta = fmt.Sprintf("%sforums/%s/posts?popular=%v", dcardAPIBase, forum[forumID].Alias, false)

	return forum[forumID].Alias
}

func getPostMeta(firstID int, lastID int) (int, int) {
	var postMeta PostMeta
	url := dcardAPIPostMeta

	if firstID != 0 {
		url = fmt.Sprintf("%s&after=%d", dcardAPIPostMeta, firstID)
	} else if lastID != 0 {
		url = fmt.Sprintf("%s&before=%d", dcardAPIPostMeta, lastID)
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
	if err := json.Unmarshal([]byte(body), &postMeta); err != nil {
		log.Fatalln(err)
	}

	for _, v := range postMeta {
		fmt.Printf("[%v](%v) -> %v: %v\n", v.ID, v.CreatedAt, len(v.Media), v.Title)

		if firstID == 0 {
			firstID = v.ID
		}

		lastID = v.ID
	}

	return firstID, lastID
}

func main() {
	currForum := getForums()
	firstID, lastID := getPostMeta(0, 0)

	scanner := bufio.NewScanner(os.Stdin)
	quit := false

	for !quit {
		fmt.Println("n: Next, p: Previous, v: View, d: Download, f: Change forum, q/quit/exit: Quit")
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
		case "f":
			fmt.Println("Change forum")
			currForum = getForums()
			firstID, lastID = getPostMeta(0, 0)
		case "n":
			fmt.Println("Next Page")
			firstID, lastID = getPostMeta(0, lastID)
		case "p":
			fmt.Println("Previous Page")
			firstID, lastID = getPostMeta(firstID, 0)
		case "v":
			fmt.Println("View Post")
		case "d":
			if len(args) == 0 {
				fmt.Println("No post specified. Try input 'd 1' to get media file")
				continue
			}

			fmt.Printf("Download post media: %v\n", args)
		}
	}
}
