package repository

import (
	"fmt"
	"log"

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

	log.Printf("Content of these authors has been affected: %v", authors)

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
	for authorIndex, author := range authors.ToArray() {
		authorsCollective.Add(author)
		for deltaIndex, fd := range deltas {
			err = fd.WriteAuthorPart(authorsCollective, repositoryPath, worktree)
			if err != nil {
				return "", err
			}

			log.Printf("Applying changes: %s (%d/%d) --> %s (%d/%d)",
				author,
				authorIndex+1,
				len(authors),
				fd.GetFilePath(),
				deltaIndex+1,
				len(deltas),
			)
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
