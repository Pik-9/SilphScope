package repository

import (
	"errors"
	"strings"

	"github.com/go-git/go-git/v5"
)

type AuthorEmail struct {
	Author, AuthorName string
}

func (ae *AuthorEmail) Authored(line *git.Line) bool {
	return line.Author == ae.Author && line.AuthorName == ae.AuthorName
}

type FileDelta struct {
	From, To         *git.BlameResult
	FromPath, ToPath string
	CommitHash       string
}

func NewFileDelta(repositoryPath string, commitHash string) ([]FileDelta, error) {
	return nil, errors.New("Not implemented, yet.")
}

func compareLines(a, b string) bool {
	l1 := strings.TrimLeft(a, " \t")
	l2 := strings.TrimLeft(b, " \t")
	return l1 == l2
}

func (fd *FileDelta) findFirstOccurenceInTo(line string, startingAt int) int {
	for ii := startingAt; ii < len(fd.To.Lines); ii++ {
		if compareLines(line, fd.To.Lines[ii].Text) {
			return ii
		}
	}
	return -1
}

func (fd *FileDelta) isLineNew(index int) bool {
	if index >= len(fd.To.Lines) {
		return false
	}

	return fd.To.Lines[index].Hash.String() == fd.CommitHash
}

func (fd *FileDelta) AssignOldAuthors() {
	b2a := make(map[int]int)
	pointer := 0
	for index, origLine := range fd.From.Lines {
		found := fd.findFirstOccurenceInTo(origLine.Text, pointer)
		if found > -1 {
			pointer = found
			b2a[found] = index
		}
	}

	leftPointer := -1
	for rightPointer := range fd.To.Lines {
		correspondingFromLine, alreadyAssigned := b2a[rightPointer]
		if alreadyAssigned {
			leftPointer = correspondingFromLine + 1
		} else {
			b2a[rightPointer] = leftPointer
		}
	}

	for index := range fd.To.Lines {
		corrLine := b2a[index]
		if fd.isLineNew(index) && corrLine > -1 && corrLine < len(fd.From.Lines) {
			fd.To.Lines[index].Author = fd.From.Lines[index].Author
			fd.To.Lines[index].AuthorName = fd.From.Lines[index].AuthorName
			fd.To.Lines[index].Date = fd.From.Lines[index].Date
		}
	}
}

func (fd *FileDelta) ContentFilteredForUsers(users []AuthorEmail) []string {
	ret := make([]string, 0)
	for index, line := range fd.To.Lines {
		if !fd.isLineNew(index) {
			ret = append(ret, line.Text)
			continue
		}

		for _, user := range users {
			if user.Authored(line) {
				ret = append(ret, line.Text)
				break
			}
		}
	}
	return ret
}

func (fd *FileDelta) GetAllAuthors() []AuthorEmail {
	authorSet := make(map[AuthorEmail]bool)
	for _, line := range fd.To.Lines {
		author := AuthorEmail{
			Author:     line.Author,
			AuthorName: line.AuthorName,
		}
		authorSet[author] = true
	}

	ret := make([]AuthorEmail, 0, len(authorSet))
	for k := range authorSet {
		ret = append(ret, k)
	}

	return ret
}
