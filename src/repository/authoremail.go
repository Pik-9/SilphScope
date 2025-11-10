package repository

import (
	"fmt"
	"github.com/go-git/go-git/v5"
)

type AuthorEmail struct {
	Author, AuthorName string
}

func AuthorFromLine(line *git.Line) AuthorEmail {
	return AuthorEmail{
		AuthorName: line.AuthorName,
		Author:     line.Author,
	}
}

func (ae *AuthorEmail) Authored(line *git.Line) bool {
	return line.Author == ae.Author && line.AuthorName == ae.AuthorName
}

func (ae AuthorEmail) String() string {
	return fmt.Sprintf("%s <%s>", ae.AuthorName, ae.Author)
}
