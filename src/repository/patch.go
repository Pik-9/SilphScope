package repository

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func getHash(s string, repo *git.Repository) plumbing.Hash {
	if s == "HEAD" {
		headRef, _ := repo.Reference(plumbing.HEAD, true)
		return headRef.Hash()
	} else {
		return plumbing.NewHash(s)
	}
}

func ExtractPatch(repositoryPath string, commitHash string) (*object.Patch, *object.Commit, *git.Repository, error) {
	repo, err := git.PlainOpen(repositoryPath)
	if err != nil {
		return nil, nil, nil, err
	}

	commit, err := repo.CommitObject(getHash(commitHash, repo))
	if err != nil {
		return nil, nil, repo, err
	}

	parent, err := commit.Parent(0)
	if err != nil {
		return nil, commit, repo, err
	}

	patch, err := commit.Patch(parent)
	if err != nil {
		return nil, commit, repo, err
	}

	return patch, commit, repo, nil
}
