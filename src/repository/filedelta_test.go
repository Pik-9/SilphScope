package repository

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var baseFile []string
var formattedAlpha []string
var formattedBeta []string
var formattedGamma []string
var repoPath string

func init() {
	repoPath = "../../testing_tmp/repo/"
	baseFile = []string{"A", "B", "C", "D", "E"}
	formattedAlpha = []string{"A", "b1+", "b2+", "b3+", "C", "d1+"}
	formattedBeta = []string{"x1+", "x2+", "A", "b1+", "D"}
	formattedGamma = []string{"A", "B", "c1+", "c2+"}
}

func setup() {
	err := os.MkdirAll(repoPath, 0750)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		log.Fatal(err)
	}

	wtree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	commit := func(author string, authorEmail string) error {
		cmt, err := wtree.Commit(author, &git.CommitOptions{
			All:               false,
			AllowEmptyCommits: false,
			Parents:           []plumbing.Hash{},
			SignKey:           nil,
			Signer:            nil,
			Amend:             false,
			Committer:         nil,
			Author: &object.Signature{
				Name:  author,
				Email: authorEmail,
				When:  time.Now(),
			},
		})

		if err != nil {
			return err
		}

		_, err = repo.CommitObject(cmt)

		if err != nil {
			return err
		}

		return nil
	}

	WriteLines(repoPath+"Alpha.txt", baseFile[0:1])
	WriteLines(repoPath+"Beta.txt", baseFile[0:1])
	WriteLines(repoPath+"Gamma.txt", baseFile[0:1])
	wtree.AddGlob("*")
	commit("Alice", "alice@testing.com")

	WriteLines(repoPath+"Alpha.txt", baseFile[0:2])
	WriteLines(repoPath+"Beta.txt", baseFile[0:2])
	WriteLines(repoPath+"Gamma.txt", baseFile[0:2])
	wtree.AddGlob("*")
	commit("Bob", "bob@testing.com")

	WriteLines(repoPath+"Alpha.txt", baseFile[0:3])
	WriteLines(repoPath+"Beta.txt", baseFile[0:3])
	wtree.AddGlob("*")
	commit("Claire", "claire@testing.com")

	WriteLines(repoPath+"Alpha.txt", baseFile[0:4])
	WriteLines(repoPath+"Beta.txt", baseFile[0:4])
	wtree.AddGlob("*")
	commit("David", "david@testing.com")

	WriteLines(repoPath+"Alpha.txt", baseFile[0:5])
	wtree.Add(repoPath + "Alpha.txt")
	wtree.AddGlob("*")
	commit("Eleonore", "eleonore@testing.com")

	WriteLines(repoPath+"Alpha.txt", formattedAlpha)
	WriteLines(repoPath+"Beta.txt", formattedBeta)
	WriteLines(repoPath+"Gamma.txt", formattedGamma)
	wtree.AddGlob("*")
	commit("Zorin", "zorin@reformatting.com")
}

func tearDown() {
	os.RemoveAll(repoPath)
}

func TestStuff(t *testing.T) {
	setup()
	defer tearDown()
}
