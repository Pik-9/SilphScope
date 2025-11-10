package repository

import (
	"errors"
	"fmt"

	"github.com/Pik-9/SilphScope/src/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CreateUnghostedBranch(repo *git.Repository, parent *object.Commit, deltas []FileDelta) (string, error) {
	var authors utils.Set[AuthorEmail] = make(map[AuthorEmail]bool)

	for _, delta := range deltas {
		authors = authors.Union(delta.GetAllOldAuthors())
	}

	fmt.Println("Authors:", authors)

	branchName := fmt.Sprintf("unghost-%s", parent.Hash.String())

	worktree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	// TODO: Make sure there are no unstashed changes.
	err = worktree.Reset(&git.ResetOptions{
		Commit: parent.Hash,
	})
	if err != nil {
		return "", err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Create: true,
		Force:  false,
		Branch: plumbing.ReferenceName("refs/heads/" + branchName),
	})
	if err != nil {
		return "", err
	}

	return "", errors.New("Not implemented, yet.")
}
