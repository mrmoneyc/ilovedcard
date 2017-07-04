package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	quit := false

	for !quit {
		fmt.Println("n: Next, p: Previous, v: View, d: Download, q/quit/exit: Quit")

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
		case "p":
			fmt.Println("Previous Page")
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
