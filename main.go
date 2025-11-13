package main

import (
	"flag"
	"log"

	"github.com/Pik-9/SilphScope/src/repository"
)

func main() {
	commitHash := flag.String("c", "HEAD", "The commit to unghost.")
	repoPath := flag.String("r", ".", "Path to repository.")
	rebaseBranch := flag.String("b", "", "Branch to rebase onto.")
	flag.Parse()

	log.Default().Println("Unghosting commit", *commitHash, "in repo", *repoPath)

	patch, commit, parent, repo, err := repository.ExtractPatch(*repoPath, *commitHash)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Extracted a patch of changes across %d files", len(patch.FilePatches()))

	deltas, err := repository.NewFileDeltas(commit, parent, patch)
	if err != nil {
		log.Fatal(err)
	}

	unghostedBranchName, err := repository.CreateUnghostedBranch(repo, *repoPath, commit, parent, deltas)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully unghosted commit %s in %s at branch %s.\n", commit.Hash.String(), *repoPath, unghostedBranchName)

	if *rebaseBranch != "" {
		rebaseOutput, err := repository.RebaseUnghostedBranchOnto(unghostedBranchName, *rebaseBranch, *repoPath)
		if err != nil {
			log.Fatalf("Error while rebasing: %s\nManual intervention needed.\n\n%s\n", err, rebaseOutput)
		}
		log.Println("Rebased onto", *rebaseBranch)
	}
}
