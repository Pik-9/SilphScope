package main

import (
	"flag"
	"log"

	"github.com/Pik-9/SilphScope/src/repository"
	"github.com/Pik-9/SilphScope/src/strategy"
)

func main() {
	commitHash := flag.String("c", "HEAD", "The commit to unghost.")
	repoPath := flag.String("r", ".", "Path to repository.")
	stratStr := flag.String("u", "author", `Strategy to unghost [author|commit].
  Strategy author will create one commit per author while discarding original commit message and date.
  Strategy commit will recreate every commit with its message and date.`)
	flag.Parse()

	strat, err := strategy.New(*stratStr)

	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println("Unghosting commit", *commitHash, "in repo", *repoPath, "while using strategy", strat)

	patch, commit, parent, repo, err := repository.ExtractPatch(*repoPath, *commitHash)
	if err != nil {
		log.Fatal(err)
	}

	deltas, err := repository.NewFileDeltas(commit, parent, patch)
	if err != nil {
		log.Fatal(err)
	}

	unghostedBranchName, err := repository.CreateUnghostedBranch(repo, *repoPath, commit, parent, deltas)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully unghosted commit %s in %s at branch %s.\n", commit.Hash.String(), *repoPath, unghostedBranchName)
}
