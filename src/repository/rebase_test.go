package repository

import (
	"testing"
	"time"

	"github.com/Pik-9/SilphScope/src/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestRebaseUnghostedBranchOnto(t *testing.T) {
	repoPath := tempSetup(t)

	_, commit, _, repo, err := ExtractPatch(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}

	reformatCommitHash := commit.Hash.String()

	wtree, _ := repo.Worktree()
	utils.WriteLines(repoPath+"/Gamma.txt", []string{"A", "B", "c1+", "c2+", "D"})
	wtree.Commit("Add David's line.", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Email: david.Author,
			Name:  david.AuthorName,
			When:  time.Now(),
		},
	})

	patch, commit, parent, repo, err := ExtractPatch(repoPath, reformatCommitHash)
	if err != nil {
		t.Error(err)
	}

	deltas, err := NewFileDeltas(commit, parent, patch)
	if err != nil {
		t.Error(err)
	}

	unghostedBranchName, err := CreateUnghostedBranch(repo, repoPath, commit, parent, deltas)
	if err != nil {
		t.Error(err)
	}

	output, err := RebaseUnghostedBranchOnto(unghostedBranchName, "master", repoPath)
	if err != nil {
		t.Error(err)
		t.Fatal(output)
	}

	expectedGammaAuthors := []AuthorEmail{alice, bob, zorin, zorin, david}

	_, commit, _, repo, err = ExtractPatch(repoPath, "HEAD")
	gammaBlame, err := git.Blame(commit, "Gamma.txt")
	if err != nil {
		t.Error(err)
	}

	for index, author := range expectedGammaAuthors {
		if !author.Authored(gammaBlame.Lines[index]) {
			t.Errorf("Author mismatch in Gamma.txt:%d: %s != %s", index, author, AuthorFromLine(gammaBlame.Lines[index]))
		}
	}
}
