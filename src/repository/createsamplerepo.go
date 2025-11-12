package repository

import (
	"github.com/Pik-9/SilphScope/src/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"crypto/rand"
	"os"
	"time"
)

func CreateSampleRepo(repoPath string) error {
	baseFile := []string{"A", "B", "C", "D", "E"}
	formattedAlpha := []string{"A", "b1+", "b2+", "b3+", "C", "d1+"}
	formattedBeta := []string{"x1+", "x2+", "A", "b1+", "D"}
	formattedGamma := []string{"A", "B", "c1+", "c2+"}

	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		return err
	}

	wtree, err := repo.Worktree()
	if err != nil {
		return err
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

	utils.WriteLines(repoPath+"/Alpha.txt", baseFile[0:1])
	utils.WriteLines(repoPath+"/Beta.txt", baseFile[0:1])
	utils.WriteLines(repoPath+"/Gamma.txt", baseFile[0:1])
	wtree.AddGlob("*")
	commit("Alice", "alice@testing.com")

	utils.WriteLines(repoPath+"/Alpha.txt", baseFile[0:2])
	utils.WriteLines(repoPath+"/Beta.txt", baseFile[0:2])
	utils.WriteLines(repoPath+"/Gamma.txt", baseFile[0:2])
	wtree.AddGlob("*")
	commit("Bob", "bob@testing.com")

	utils.WriteLines(repoPath+"/Alpha.txt", baseFile[0:3])
	utils.WriteLines(repoPath+"/Beta.txt", baseFile[0:3])
	wtree.AddGlob("*")
	commit("Claire", "claire@testing.com")

	utils.WriteLines(repoPath+"/Alpha.txt", baseFile[0:4])
	utils.WriteLines(repoPath+"/Beta.txt", baseFile[0:4])
	wtree.AddGlob("*")
	commit("David", "david@testing.com")

	data := make([]byte, 1024)
	rand.Read(data)
	file, _ := os.Create(repoPath + "/blob_eleonore.bin")
	file.Write(data)
	file.Close()
	wtree.AddGlob("*")

	utils.WriteLines(repoPath+"/Alpha.txt", baseFile[0:5])
	wtree.AddGlob("*")
	commit("Eleonore", "eleonore@testing.com")

	utils.WriteLines(repoPath+"/Alpha.txt", formattedAlpha)
	utils.WriteLines(repoPath+"/Beta.txt", formattedBeta)
	utils.WriteLines(repoPath+"/Gamma.txt", formattedGamma)
	data = make([]byte, 512)
	rand.Read(data)
	file, _ = os.Create(repoPath + "/blob_zorin.bin")
	file.Write(data)
	file.Close()
	wtree.AddGlob("*")
	commit("Zorin", "zorin@reformatting.com")
	baseFile = []string{"A", "B", "C", "D", "E"}
	formattedAlpha = []string{"A", "b1+", "b2+", "b3+", "C", "d1+"}
	formattedBeta = []string{"x1+", "x2+", "A", "b1+", "D"}
	formattedGamma = []string{"A", "B", "c1+", "c2+"}

	return nil
}
