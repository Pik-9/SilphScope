package repository

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var baseFile []string
var formattedAlpha []string
var formattedBeta []string
var formattedGamma []string
var repoPath string
var alice, bob, claire, david, eleonore, zorin AuthorEmail

func init() {
	repoPath = "../../testing_tmp/repo/"
	baseFile = []string{"A", "B", "C", "D", "E"}
	formattedAlpha = []string{"A", "b1+", "b2+", "b3+", "C", "d1+"}
	formattedBeta = []string{"x1+", "x2+", "A", "b1+", "D"}
	formattedGamma = []string{"A", "B", "c1+", "c2+"}

	alice = AuthorEmail{
		AuthorName: "Alice",
		Author:     "alice@testing.com",
	}

	bob = AuthorEmail{
		AuthorName: "Bob",
		Author:     "bob@testing.com",
	}

	claire = AuthorEmail{
		AuthorName: "Claire",
		Author:     "claire@testing.com",
	}

	david = AuthorEmail{
		AuthorName: "David",
		Author:     "david@testing.com",
	}

	eleonore = AuthorEmail{
		AuthorName: "Eleonore",
		Author:     "eleonore@testing.com",
	}

	zorin = AuthorEmail{
		AuthorName: "Zorin",
		Author:     "zorin@reformatting.com",
	}
}

func compareSlices[P comparable](t *testing.T, expected []P, actual []P) {
	if len(expected) != len(actual) {
		t.Errorf("Lenghts mismatch: %d != %d\n  --> %v != %v", len(expected), len(actual), expected, actual)
		return
	}

	for index, exp := range expected {
		if exp != actual[index] {
			t.Errorf("Failure at index %d: %v != %v", index, exp, actual[index])
		}
	}
}

func setup() {
	err := os.MkdirAll(repoPath, 0750)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		log.Fatal(err)
	}

	wtree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	commit := func(author string, authorEmail string) error {
		cmt, err := wtree.Commit(author, &git.CommitOptions{
			All:               false,
			AllowEmptyCommits: false,
			Parents:           []plumbing.Hash{},
			SignKey:           nil,
			Signer:            nil,
			Amend:             false,
			Committer:         nil,
			Author: &object.Signature{
				Name:  author,
				Email: authorEmail,
				When:  time.Now(),
			},
		})

		if err != nil {
			return err
		}

		_, err = repo.CommitObject(cmt)

		if err != nil {
			return err
		}

		return nil
	}

	WriteLines(repoPath+"Alpha.txt", baseFile[0:1])
	WriteLines(repoPath+"Beta.txt", baseFile[0:1])
	WriteLines(repoPath+"Gamma.txt", baseFile[0:1])
	wtree.AddGlob("*")
	commit("Alice", "alice@testing.com")

	WriteLines(repoPath+"Alpha.txt", baseFile[0:2])
	WriteLines(repoPath+"Beta.txt", baseFile[0:2])
	WriteLines(repoPath+"Gamma.txt", baseFile[0:2])
	wtree.AddGlob("*")
	commit("Bob", "bob@testing.com")

	WriteLines(repoPath+"Alpha.txt", baseFile[0:3])
	WriteLines(repoPath+"Beta.txt", baseFile[0:3])
	wtree.AddGlob("*")
	commit("Claire", "claire@testing.com")

	WriteLines(repoPath+"Alpha.txt", baseFile[0:4])
	WriteLines(repoPath+"Beta.txt", baseFile[0:4])
	wtree.AddGlob("*")
	commit("David", "david@testing.com")

	WriteLines(repoPath+"Alpha.txt", baseFile[0:5])
	wtree.Add(repoPath + "Alpha.txt")
	wtree.AddGlob("*")
	commit("Eleonore", "eleonore@testing.com")

	WriteLines(repoPath+"Alpha.txt", formattedAlpha)
	WriteLines(repoPath+"Beta.txt", formattedBeta)
	WriteLines(repoPath+"Gamma.txt", formattedGamma)
	wtree.AddGlob("*")
	commit("Zorin", "zorin@reformatting.com")
}

func tearDown() {
	os.RemoveAll(repoPath)
}

func TestFileDelta_GetAllOldAuthors(t *testing.T) {
	setup()
	defer tearDown()

	expectedAlpha := Unique([]AuthorEmail{
		{AuthorName: "Alice", Author: "alice@testing.com"},
		{AuthorName: "Zorin", Author: "zorin@reformatting.com"},
		// {AuthorName: "Bob", Author: "bob@testing.com"},
		{AuthorName: "Claire", Author: "claire@testing.com"},
		// {AuthorName: "David", Author: "david@testing.com"},
		// {AuthorName: "Eleonore", Author: "eleonore@testing.com"},
	})

	expectedBeta := Unique([]AuthorEmail{
		{AuthorName: "Zorin", Author: "zorin@reformatting.com"},
		{AuthorName: "Alice", Author: "alice@testing.com"},
		// {AuthorName: "Bob", Author: "bob@testing.com"},
		// {AuthorName: "Claire", Author: "claire@testing.com"},
		{AuthorName: "David", Author: "david@testing.com"},
		// {AuthorName: "Eleonore", Author: "eleonore@testing.com"},
	})

	expectedGamma := Unique([]AuthorEmail{
		{AuthorName: "Alice", Author: "alice@testing.com"},
		{AuthorName: "Bob", Author: "bob@testing.com"},
		{AuthorName: "Zorin", Author: "zorin@reformatting.com"},
		// {AuthorName: "Claire", Author: "claire@testing.com"},
		// {AuthorName: "David", Author: "david@testing.com"},
		// {AuthorName: "Eleonore", Author: "eleonore@testing.com"},
	})

	deltas, err := NewFileDeltas(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}

	oldAuthorsAlpha := deltas[0].GetAllOldAuthors()
	oldAuthorsBeta := deltas[1].GetAllOldAuthors()
	oldAuthorsGamma := deltas[2].GetAllOldAuthors()

	if !oldAuthorsAlpha.Equals(expectedAlpha) {
		t.Errorf("The authors for file Alpha don't match: %v != %v", oldAuthorsAlpha.ToArray(), expectedAlpha.ToArray())
	}

	if !oldAuthorsBeta.Equals(expectedBeta) {
		t.Errorf("The authors for file Beta don't match: %v != %v", oldAuthorsBeta.ToArray(), expectedBeta.ToArray())
	}

	if !oldAuthorsGamma.Equals(expectedGamma) {
		t.Errorf("The authors for file Gamma don't match: %v != %v", oldAuthorsGamma.ToArray(), expectedGamma.ToArray())
	}
}

func TestUnghostAuthors(t *testing.T) {
	setup()
	defer tearDown()

	expectedAlpha := []AuthorEmail{
		alice,
		bob,
		bob,
		bob,
		claire,
		david,
	}

	expectedBeta := []AuthorEmail{
		zorin,
		zorin,
		alice,
		bob,
		david,
	}

	expectedGamma := []AuthorEmail{
		alice,
		bob,
		zorin,
		zorin,
	}

	deltas, err := NewFileDeltas(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}

	newAuthorsAlpha := deltas[0].UnghostAuthors()
	newAuthorsBeta := deltas[1].UnghostAuthors()
	newAuthorsGamma := deltas[2].UnghostAuthors()

	compareSlices(t, expectedAlpha, newAuthorsAlpha)
	compareSlices(t, expectedBeta, newAuthorsBeta)
	compareSlices(t, expectedGamma, newAuthorsGamma)
}

func TestFileDelta_ContentFilteredForUsers(t *testing.T) {
	setup()
	defer tearDown()

	deltas, err := NewFileDeltas(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}

	authors := []AuthorEmail{zorin, alice, bob, claire, david, eleonore}

	compareSlices(t, deltas[0].ContentFilteredForUsers(Unique([]AuthorEmail{})), []string{"A", "C"})
	compareSlices(t, deltas[0].ContentFilteredForUsers(Unique(authors[:1])), []string{"A", "C"})
	compareSlices(t, deltas[0].ContentFilteredForUsers(Unique(authors[:2])), []string{"A", "C"})
	compareSlices(t, deltas[0].ContentFilteredForUsers(Unique(authors[:3])), []string{"A", "b1+", "b2+", "b3+", "C"})
	compareSlices(t, deltas[0].ContentFilteredForUsers(Unique(authors[:4])), []string{"A", "b1+", "b2+", "b3+", "C"})
	compareSlices(t, deltas[0].ContentFilteredForUsers(Unique(authors[:5])), []string{"A", "b1+", "b2+", "b3+", "C", "d1+"})

	compareSlices(t, deltas[1].ContentFilteredForUsers(Unique([]AuthorEmail{})), []string{"A", "D"})
	compareSlices(t, deltas[1].ContentFilteredForUsers(Unique(authors[:1])), []string{"x1+", "x2+", "A", "D"})
	compareSlices(t, deltas[1].ContentFilteredForUsers(Unique(authors[:2])), []string{"x1+", "x2+", "A", "D"})
	compareSlices(t, deltas[1].ContentFilteredForUsers(Unique(authors[:3])), []string{"x1+", "x2+", "A", "b1+", "D"})
	compareSlices(t, deltas[1].ContentFilteredForUsers(Unique(authors)), []string{"x1+", "x2+", "A", "b1+", "D"})

	compareSlices(t, deltas[2].ContentFilteredForUsers(Unique([]AuthorEmail{})), []string{"A", "B"})
	compareSlices(t, deltas[2].ContentFilteredForUsers(Unique(authors[:1])), []string{"A", "B", "c1+", "c2+"})
	compareSlices(t, deltas[2].ContentFilteredForUsers(Unique(authors[:2])), []string{"A", "B", "c1+", "c2+"})
	compareSlices(t, deltas[2].ContentFilteredForUsers(Unique(authors[:3])), []string{"A", "B", "c1+", "c2+"})
	compareSlices(t, deltas[2].ContentFilteredForUsers(Unique(authors[:4])), []string{"A", "B", "c1+", "c2+"})
	compareSlices(t, deltas[2].ContentFilteredForUsers(Unique(authors[:5])), []string{"A", "B", "c1+", "c2+"})
}
