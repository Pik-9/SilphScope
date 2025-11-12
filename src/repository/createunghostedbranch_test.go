package repository

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
)

func Test_CreateUnghostedBranch(t *testing.T) {
	setup()
	t.Cleanup(tearDown)

	patch, commit, parentCommit, repo, err := ExtractPatch(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}
	deltas, err := NewFileDeltas(commit, parentCommit, patch)
	if err != nil {
		t.Error(err)
	}

	unghostedBranch, err := CreateUnghostedBranch(repo, repoPath, commit, parentCommit, deltas)
	if err != nil {
		t.Error(err)
	}

	unghostedBranchRef, err := repo.Reference(plumbing.NewBranchReferenceName(unghostedBranch), true)
	if err != nil {
		t.Error(err)
	}

	unghostedCommit, err := repo.CommitObject(unghostedBranchRef.Hash())
	if err != nil {
		t.Error(err)
	}

	branchDiff, err := commit.Patch(unghostedCommit)
	if err != nil {
		t.Error(err)
	}

	filePatches := branchDiff.FilePatches()
	if len(filePatches) != 0 {
		t.Errorf("There are %d differences:", len(filePatches))
		for _, fp := range filePatches {
			_, toFile := fp.Files()
			t.Error("  In file:", toFile.Path())
			var prefix string
			for _, chunk := range fp.Chunks() {
				if chunk.Type() == diff.Equal {
					continue
				}

				if chunk.Type() == diff.Add {
					prefix = "    +"
				}
				if chunk.Type() == diff.Delete {
					prefix = "    -"
				}

				t.Errorf("%s %s", prefix, chunk.Content())
			}
		}
	}

	expectedAlphaAuthors := []AuthorEmail{alice, bob, bob, bob, claire, david}
	expectedBetaAuthors := []AuthorEmail{zorin, zorin, alice, bob, david}
	expectedGammaAuthors := []AuthorEmail{alice, bob, zorin, zorin}

	alphaBlame, err := git.Blame(unghostedCommit, "Alpha.txt")
	if err != nil {
		t.Error(err)
	}
	betaBlame, err := git.Blame(unghostedCommit, "Beta.txt")
	if err != nil {
		t.Error(err)
	}
	gammaBlame, err := git.Blame(unghostedCommit, "Gamma.txt")
	if err != nil {
		t.Error(err)
	}

	for lineNumber, line := range alphaBlame.Lines {
		if !expectedAlphaAuthors[lineNumber].Authored(line) {
			t.Errorf("Dicrepancy in Alpha.txt:%d: %s != %s", lineNumber, expectedAlphaAuthors[lineNumber], AuthorFromLine(line))
		}
	}
	for lineNumber, line := range betaBlame.Lines {
		if !expectedBetaAuthors[lineNumber].Authored(line) {
			t.Errorf("Dicrepancy in Beta.txt:%d: %s != %s", lineNumber, expectedBetaAuthors[lineNumber], AuthorFromLine(line))
		}
	}
	for lineNumber, line := range gammaBlame.Lines {
		if !expectedGammaAuthors[lineNumber].Authored(line) {
			t.Errorf("Dicrepancy in Gamma.txt:%d: %s != %s", lineNumber, expectedGammaAuthors[lineNumber], AuthorFromLine(line))
		}
	}
}
