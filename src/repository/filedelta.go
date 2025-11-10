package repository

import (
	"errors"
	"fmt"

	"strings"

	"github.com/Pik-9/SilphScope/src/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type FileDelta struct {
	From, To         *git.BlameResult
	FromPath, ToPath string
	CommitAuthor     AuthorEmail
}

func NewFileDeltas(commit *object.Commit, parentCommit *object.Commit, patch diff.Patch) ([]FileDelta, error) {
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
			CommitAuthor: AuthorEmail{Author: commit.Author.Email, AuthorName: commit.Author.Name},
			FromPath:     filePath(fromPath),
			ToPath:       filePath(toPath),
		}

		if fp.IsBinary() {
			return nil, errors.New("The commit must not contain any binary files.")
		}

		var err error
		if fromPath != nil {
			fd.From, err = git.Blame(parentCommit, fromPath.Path())
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
	if index < 0 || index >= len(fd.To.Lines) {
		return false
	}

	return AuthorFromLine(fd.To.Lines[index]) == fd.CommitAuthor
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

	for index, toLine := range fd.To.Lines {
		corrLine := b2a[index]
		if corrLine < 0 || corrLine >= len(fd.From.Lines) {
			ret = append(ret, AuthorFromLine(toLine))
		} else {
			ret = append(ret, AuthorFromLine(fd.From.Lines[corrLine]))
		}
	}

	return ret
}

func (fd *FileDelta) ContentFilteredForUsers(users utils.Set[AuthorEmail]) []string {
	ret := make([]string, 0)
	unghostedLineAuthors := fd.UnghostAuthors()
	for index, line := range fd.To.Lines {
		if !fd.isLineNew(index) {
			ret = append(ret, line.Text)
			continue
		}

		if users.Contains(unghostedLineAuthors[index]) {
			ret = append(ret, line.Text)
		}
	}
	return ret
}

func (fd *FileDelta) GetAllOldAuthors() utils.Set[AuthorEmail] {
	authorSet := make(map[AuthorEmail]bool)
	for _, line := range fd.To.Lines {
		author := AuthorFromLine(line)
		authorSet[author] = true
	}

	return authorSet
}

func (fd *FileDelta) PrintUnghosted() {
	unghostedLineAuthors := fd.UnghostAuthors()
	for index, line := range fd.To.Lines {
		fmt.Printf("%d\t%s\t%s\n", index, unghostedLineAuthors[index], line.Text)
	}
}
