package repository

import (
	"fmt"

	"strings"

	"github.com/Pik-9/SilphScope/src/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type sourceLine struct {
	Line   string
	Author AuthorEmail
}

type FileDelta interface {
	ContentForUsers(utils.Set[AuthorEmail]) []byte
	GetAllOldAuthors() utils.Set[AuthorEmail]
	WriteUnalteredPart(string, *git.Worktree) error
	WriteAuthorPart(utils.Set[AuthorEmail], string, *git.Worktree) error
}

type TextFileDelta struct {
	From, To         *git.BlameResult
	FromPath, ToPath string
	Commit           *object.Commit
}

type BinaryFileDelta struct {
	FromPath, ToPath string
	Commit           *object.Commit
	Content          []byte
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
		var lret FileDelta
		toPath, fromPath := fp.Files()

		if fp.IsBinary() {
			fd := BinaryFileDelta{
				FromPath: filePath(fromPath),
				ToPath:   filePath(toPath),
				Commit:   commit,
				Content:  []byte{},
			}

			fileIter, _ := commit.Files()
			err := fileIter.ForEach(func(file *object.File) error {
				if file.Name == fd.ToPath {
					fd.Content = make([]byte, file.Blob.Size)
					reader, err := file.Blob.Reader()
					if err != nil {
						return err
					}
					defer reader.Close()

					_, err = reader.Read(fd.Content)
					if err != nil {
						return err
					}

					return storer.ErrStop
				}

				return nil
			})

			if err != nil {
				return nil, err
			}

			lret = fd
		} else {
			fd := TextFileDelta{
				Commit:   commit,
				FromPath: filePath(fromPath),
				ToPath:   filePath(toPath),
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

			lret = fd
		}

		ret = append(ret, lret)
	}

	return ret, nil
}

func (fd TextFileDelta) GetToLines() []sourceLine {
	if fd.To == nil {
		return []sourceLine{}
	}

	ret := make([]sourceLine, 0, len(fd.To.Lines))
	for _, line := range fd.To.Lines {
		author := AuthorFromLine(line)
		ret = append(ret, sourceLine{Line: line.Text, Author: author})
	}

	return ret
}

func (fd TextFileDelta) GetFromLines() []sourceLine {
	if fd.From == nil {
		return []sourceLine{}
	}

	ret := make([]sourceLine, 0, len(fd.From.Lines))
	for _, line := range fd.From.Lines {
		author := AuthorFromLine(line)
		ret = append(ret, sourceLine{Line: line.Text, Author: author})
	}

	return ret
}

func compareLines(a, b string) bool {
	l1 := strings.ReplaceAll(strings.TrimSpace(a), " ", "")
	l2 := strings.ReplaceAll(strings.TrimSpace(b), " ", "")
	return l1 == l2
}

func (fd TextFileDelta) findFirstOccurenceInTo(line string, startingAt int) int {
	toLines := fd.GetToLines()
	for ii := startingAt; ii < len(toLines); ii++ {
		if compareLines(line, toLines[ii].Line) {
			return ii
		}
	}
	return -1
}

func (fd TextFileDelta) isLineNew(index int) bool {
	toLines := fd.GetToLines()
	if index < 0 || index >= len(toLines) {
		return false
	}

	return fd.To.Lines[index].Hash == fd.Commit.Hash
}

func (fd TextFileDelta) UnghostAuthors() []AuthorEmail {
	fromLines := fd.GetFromLines()
	toLines := fd.GetToLines()

	b2a := make(map[int]int)
	pointer := 0
	for index, origLine := range fromLines {
		found := fd.findFirstOccurenceInTo(origLine.Line, pointer)
		if found > -1 {
			pointer = found
			b2a[found] = index
		}
	}

	leftPointer := -1
	for rightPointer := range toLines {
		correspondingFromLine, alreadyAssigned := b2a[rightPointer]
		if alreadyAssigned {
			leftPointer = correspondingFromLine + 1
		} else {
			b2a[rightPointer] = leftPointer
		}
	}

	ret := make([]AuthorEmail, 0, len(toLines))

	for index, toLine := range toLines {
		corrLine := b2a[index]
		if corrLine < 0 || corrLine >= len(fromLines) {
			ret = append(ret, toLine.Author)
		} else {
			ret = append(ret, fromLines[corrLine].Author)
		}
	}

	return ret
}

func (fd TextFileDelta) GetUnalteredLines() []string {
	ret := make([]string, 0)

	for index, line := range fd.GetToLines() {
		if !fd.isLineNew(index) {
			ret = append(ret, line.Line)
		}
	}

	return ret
}

func (fd TextFileDelta) ContentFilteredForUsers(users utils.Set[AuthorEmail]) []string {
	ret := make([]string, 0)
	unghostedLineAuthors := fd.UnghostAuthors()
	for index, line := range fd.GetToLines() {
		if !fd.isLineNew(index) {
			ret = append(ret, line.Line)
			continue
		}

		if users.Contains(unghostedLineAuthors[index]) {
			ret = append(ret, line.Line)
		}
	}
	return ret
}

func (fd TextFileDelta) GetAllOldAuthors() utils.Set[AuthorEmail] {
	authorSet := make(map[AuthorEmail]bool)
	for _, line := range fd.GetToLines() {
		authorSet[line.Author] = true
	}

	return authorSet
}

func (fd BinaryFileDelta) GetAllOldAuthors() utils.Set[AuthorEmail] {
	return utils.Unique([]AuthorEmail{{
		Author:     fd.Commit.Author.Email,
		AuthorName: fd.Commit.Author.Name,
	}})
}

func (fd TextFileDelta) PrintUnghosted() {
	unghostedLineAuthors := fd.UnghostAuthors()
	for index, line := range fd.GetToLines() {
		fmt.Printf("%d\t%s\t%s\n", index, unghostedLineAuthors[index], line.Line)
	}
}

func (fd TextFileDelta) ContentForUsers(users utils.Set[AuthorEmail]) []byte {
	ret := make([]byte, 0)
	for _, line := range fd.ContentFilteredForUsers(users) {
		ret = append(ret, []byte(line)...)
		ret = append(ret, '\n')
	}
	return ret
}

func (fd BinaryFileDelta) ContentForUsers(users utils.Set[AuthorEmail]) []byte {
	author := AuthorEmail{
		Author:     fd.Commit.Author.Email,
		AuthorName: fd.Commit.Author.Name,
	}

	if users.Contains(author) {
		return fd.Content
	} else {
		return []byte{}
	}
}

func (fd TextFileDelta) WriteUnalteredPart(repoPath string, worktree *git.Worktree) error {
	if fd.FromPath != "" && fd.FromPath != fd.ToPath {
		if fd.ToPath == "" {
			_, err := worktree.Remove(fd.FromPath)
			if err != nil {
				return err
			}
		}

		_, err := worktree.Move(fd.FromPath, fd.ToPath)
		if err != nil {
			return err
		}
	}

	err := utils.WriteLines(repoPath+"/"+fd.ToPath, fd.GetUnalteredLines())
	if err != nil {
		return err
	}

	_, err = worktree.Add(fd.ToPath)
	return err
}

func (fd BinaryFileDelta) WriteUnalteredPart(repoPath string, worktree *git.Worktree) error {
	if fd.FromPath != "" && fd.FromPath != fd.ToPath {
		if fd.ToPath == "" {
			_, err := worktree.Remove(fd.FromPath)
			if err != nil {
				return err
			}
		}

		_, err := worktree.Move(fd.FromPath, fd.ToPath)
		if err != nil {
			return err
		}
	}

	// Don't write any content, yet
	return nil
}

func (fd TextFileDelta) WriteAuthorPart(authors utils.Set[AuthorEmail], repoPath string, worktree *git.Worktree) error {
	err := utils.WriteBuffer(repoPath+"/"+fd.ToPath, fd.ContentForUsers(authors))
	if err != nil {
		return err
	}

	_, err = worktree.Add(fd.ToPath)
	return err
}

func (fd BinaryFileDelta) WriteAuthorPart(authors utils.Set[AuthorEmail], repoPath string, worktree *git.Worktree) error {
	err := utils.WriteBuffer(repoPath+"/"+fd.ToPath, fd.ContentForUsers(authors))
	if err != nil {
		return err
	}

	_, err = worktree.Add(fd.ToPath)
	return err
}
