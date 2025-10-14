package main

import (
	"flag"
	"fmt"
	"log"
)

var str2strat = map[string]Strategy{
	"author": Author,
	"commit": Commit,
}

var strat2str = map[Strategy]string{
	Author: "author",
	Commit: "commit",
}

func main() {
	commit := flag.String("c", "HEAD", "The commit to unghost.")
	repo := flag.String("r", ".", "Path to repository.")
	stratStr := flag.String("u", "author", `Strategy to unghost [author|commit].
  Strategy author will create one commit per author while discarding original commit message and date.
  Strategy commit will recreate every commit with its message and date.`)
	flag.Parse()

	strat, ok := str2strat[*stratStr]

	if !ok {
		log.Fatalf("%s is not a valid strategy.", *stratStr)
	}

	fmt.Println("Unghosting commit", *commit, "in repo", *repo, "while using strategy", strat2str[strat])
}
