package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	dcardAPIBase   string = "https://www.dcard.tw/_api/"
	dcardAPIForums string = dcardAPIBase + "forums"
	dcardAPIPost   string = dcardAPIBase + "posts"
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

func getPostMeta(firstID int, lastID int) (int, int, PostMeta) {
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

	for k, v := range postMeta {
		fmt.Printf("[%v](%v) -> %v: %v\n", k, v.CreatedAt, len(v.Media), v.Title)

		if firstID == 0 {
			firstID = v.ID
		}

		lastID = v.ID
	}

	return firstID, lastID, postMeta
}

func getPost(postID int) {
	var post Post
	url := fmt.Sprintf("%s/%d", dcardAPIPost, postID)

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
	if err := json.Unmarshal([]byte(body), &post); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("========================================")
	fmt.Printf("Title: %v\n", post.Title)
	fmt.Println("========================================")
	fmt.Printf("Content: \n%v\n", post.Content)
	fmt.Println("========================================")
}

func main() {
	currForum := getForums()
	firstID, lastID, postMeta := getPostMeta(0, 0)

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
			firstID, lastID, postMeta = getPostMeta(0, 0)
		case "n":
			fmt.Println("Next Page")
			firstID, lastID, postMeta = getPostMeta(0, lastID)
		case "p":
			fmt.Println("Previous Page")
			firstID, lastID, postMeta = getPostMeta(firstID, 0)
		case "v":
			if len(args) == 0 {
				fmt.Println("No post specified. Try input 'v 1' to view article")
				continue
			}

			i, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				continue
			}

			if i > len(postMeta) || i < 0 {
				fmt.Println("Not a valid post ID. Try input 'v 1' to view article")
				continue
			}

			getPost(postMeta[i].ID)
		case "d":
			if len(args) == 0 {
				fmt.Println("No post specified. Try input 'd 1' to get media file")
				continue
			}

			fmt.Printf("Download post media: %v\n", args)
		}
	}
}
