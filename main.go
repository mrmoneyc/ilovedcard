package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	dcardAPIBase   string = "https://www.dcard.tw/_api/"
	dcardAPIForums string = dcardAPIBase + "forums"
	dcardAPIPost   string = dcardAPIBase + "posts"
)

var (
	dcardAPIPostMeta string
	currForum        string
	dlPath           string
	numOfWorker      int
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

func getPostMeta(firstID int, lastID int) (int, int, []PostMeta) {
	var postMeta []PostMeta
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
		fmt.Printf("[#%v] %v pics (%v): %v\n", k, len(v.Media), v.CreatedAt, v.Title)

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
	fmt.Printf("(%v Comments)\n", post.CommentCount)
}

func getPostMedia(postMeta PostMeta) {
	var wg sync.WaitGroup
	ch := make(chan string)
	destDir := fmt.Sprintf("%v/%v/%v_%v", dlPath, postMeta.ForumAlias, postMeta.ID, postMeta.Title)

	for i := 0; i < numOfWorker; i++ {
		wg.Add(1)
		go downloadWorker(&wg, ch, destDir)
	}

	for _, v := range postMeta.Media {
		ch <- v.URL
	}
}

func showComments(postID int) {
	var comments Comments
	url := fmt.Sprintf("%s/%d/comments", dcardAPIPost, postID)

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

	if err := json.Unmarshal([]byte(body), &comments); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("---------------------------------------------------------")
	for index, comment := range comments {
		userDescription := fmt.Sprintf("%v%v", comment.School, comment.Department)
		if comment.Anonymous {
			userDescription = "匿名"
		}
		genderDescription := "男"
		if comment.Gender == "F" {
			genderDescription = "女"
		}
		fmt.Printf("B%v %v(%v):\n\n%v \n", index+1, userDescription, genderDescription, comment.Content)
		fmt.Println("---------------------------------------------------------")
	}
}

func downloadWorker(wg *sync.WaitGroup, ch chan string, destDir string) {
	defer wg.Done()

	for dlURL := range ch {
		slURL := strings.Split(dlURL, "/")
		fileName := slURL[len(slURL)-1]
		log.Printf("Get %v from %v", fileName, dlURL)

		resp, err := http.Get(dlURL)
		if err != nil {
			log.Println(err)
			continue
		}
		defer resp.Body.Close()

		if err := os.MkdirAll(destDir, 0755); err != nil {
			log.Println(err)
			continue
		}

		f, err := os.Create(filepath.Join(destDir, fileName))
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		f.Close()
	}
}

func main() {
	numOfWorker = 1
	homeDir := os.Getenv("HOME")
	dlPath = fmt.Sprintf("%v/%v/%v", homeDir, "Downloads", "ILoveDcard")

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

			postID := postMeta[i].ID
			getPost(postID)
			showComments(postID)
		case "d":
			if len(args) == 0 {
				fmt.Println("No post specified. Try input 'd 1' to get media file")
				continue
			}

			for _, v := range args {
				i, err := strconv.Atoi(v)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if i > len(postMeta) || i < 0 {
					fmt.Println("Not a valid post ID. Try input 'v 1' to view article")
					continue
				}

				getPostMedia(postMeta[i])
			}

			fmt.Printf("Download post media: %v\n", args)
		}
	}
}
