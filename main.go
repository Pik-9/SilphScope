package main

import (
	"flag"
	"fmt"
	"github.com/Pik-9/SilphScope/src/strategy"
	"log"
)

func main() {
	commit := flag.String("c", "HEAD", "The commit to unghost.")
	repo := flag.String("r", ".", "Path to repository.")
	stratStr := flag.String("u", "author", `Strategy to unghost [author|commit].
  Strategy author will create one commit per author while discarding original commit message and date.
  Strategy commit will recreate every commit with its message and date.`)
	flag.Parse()

	strat, err := strategy.New(*stratStr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unghosting commit", *commit, "in repo", *repo, "while using strategy", strat)
}
