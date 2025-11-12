package repository

import (
	"fmt"

	"github.com/Pik-9/SilphScope/src/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CreateUnghostedBranch(repo *git.Repository, repositoryPath string, commit *object.Commit, parent *object.Commit, deltas []FileDelta) (string, error) {
	var authors utils.Set[AuthorEmail] = make(map[AuthorEmail]bool)

	for _, delta := range deltas {
		authors = authors.Union(delta.GetAllOldAuthors())
	}

	branchName := fmt.Sprintf("unghost-%s", commit.Hash.String())

	worktree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Hash:   parent.Hash,
		Create: true,
		Force:  false,
		Branch: plumbing.ReferenceName("refs/heads/" + branchName),
	})
	if err != nil {
		return "", err
	}

	// Initial state with all formatted lines removed:
	for _, fd := range deltas {
		err := fd.WriteUnalteredPart(repositoryPath, worktree)
		if err != nil {
			return "", err
		}
	}

	_, err = worktree.Commit(commit.Message, &git.CommitOptions{
		All:               true,
		AllowEmptyCommits: true,
		Author:            &commit.Committer,
		Committer:         &commit.Committer,
	})
	if err != nil {
		return "", err
	}

	// Changes from every author
	var authorsCollective utils.Set[AuthorEmail] = make(map[AuthorEmail]bool, len(authors))
	for author := range authors {
		authorsCollective.Add(author)
		for _, fd := range deltas {
			err = fd.WriteAuthorPart(authorsCollective, repositoryPath, worktree)
			if err != nil {
				return "", err
			}
		}

		_, err = worktree.Commit(commit.Message, &git.CommitOptions{
			All:               false,
			AllowEmptyCommits: true,
			Author: &object.Signature{
				Name:  author.AuthorName,
				Email: author.Author,
				When:  commit.Author.When,
			},
			Committer: &commit.Committer,
		})
		if err != nil {
			return "", err
		}
	}

	return branchName, nil
}
