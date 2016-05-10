package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tweethour"
)

func main() {

	h, err := tweethour.NewHistogram()

	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("Enter username or q to quit: ")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		text = strings.TrimSpace(text)

		if text == "q" {
			return
		}

		fmt.Printf("Processing for %s ...\n", text)

		username := text

		tbyHour, err := h.Get(username)

		if err != nil {
			fmt.Println(err)
			continue
		}

		count := 0
		for _, r := range tbyHour {
			fmt.Printf("hour: %s tweets: %d\n", r.From, r.Count)
			count = count + r.Count

		}

		fmt.Printf("User: %s, Total tweets today(UTC time): %d\n", username, count)

	}
}
