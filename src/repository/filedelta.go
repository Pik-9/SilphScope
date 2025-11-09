package repository

import (
	"errors"

	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
)

type FileDelta struct {
	From, To         *git.BlameResult
	FromPath, ToPath string
	CommitHash       string
}

func NewFileDeltas(repositoryPath string, commitHash string) ([]FileDelta, error) {
	patch, commit, parent, _, err := ExtractPatch(repositoryPath, commitHash)
	if err != nil {
		return nil, err
	}

	filePath := func(file diff.File) string {
		if file == nil {
			return ""
		}
		return file.Path()
	}

	ret := make([]FileDelta, 0, len(patch.FilePatches()))
	for _, fp := range patch.FilePatches() {
		fromPath, toPath := fp.Files()
		fd := FileDelta{
			CommitHash: commitHash,
			FromPath:   filePath(fromPath),
			ToPath:     filePath(toPath),
		}

		if fp.IsBinary() {
			return nil, errors.New("The commit must not contain any binary files.")
		}

		if fromPath != nil {
			fd.From, err = git.Blame(parent, fromPath.Path())
			if err != nil {
				return nil, err
			}
		}

		if toPath != nil {
			fd.To, err = git.Blame(commit, toPath.Path())
			if err != nil {
				return nil, err
			}
		}

		ret = append(ret, fd)
	}

	return ret, nil
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

func (fd *FileDelta) UnghostAuthors() []AuthorEmail {
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

	ret := make([]AuthorEmail, 0, len(fd.To.Lines))

	for index := range fd.To.Lines {
		corrLine := b2a[index]
		if fd.isLineNew(index) && corrLine > -1 && corrLine < len(fd.From.Lines) {
			ret = append(ret, AuthorEmail{
				Author:     fd.To.Lines[corrLine].Author,
				AuthorName: fd.To.Lines[corrLine].AuthorName,
			})
		} else {
			ret = append(ret, AuthorEmail{
				Author:     fd.To.Lines[index].Author,
				AuthorName: fd.To.Lines[index].AuthorName,
			})
		}
	}

	return ret
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

func (fd *FileDelta) GetAllOldAuthors() Set[AuthorEmail] {
	authorSet := make(map[AuthorEmail]bool)
	for _, line := range fd.To.Lines {
		author := AuthorEmail{
			Author:     line.Author,
			AuthorName: line.AuthorName,
		}
		authorSet[author] = true
	}

	return authorSet
}
