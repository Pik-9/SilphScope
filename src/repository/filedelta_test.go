package repository

import (
	"log"
	"os"
	"testing"

	"github.com/Pik-9/SilphScope/src/utils"
)

var repoPath string
var alice, bob, claire, david, eleonore, zorin AuthorEmail

func init() {
	repoPath = "../../testing_tmp/repo/"

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

func setup() {
	err := os.MkdirAll(repoPath, 0750)
	if err != nil {
		log.Fatal(err)
	}

	err = CreateSampleRepo(repoPath)
	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	os.RemoveAll(repoPath)
}

func TestFileDelta_GetAllOldAuthors(t *testing.T) {
	setup()
	t.Cleanup(tearDown)

	expectedAlpha := utils.Unique([]AuthorEmail{alice, zorin, claire})
	expectedBeta := utils.Unique([]AuthorEmail{zorin, alice, david})
	expectedGamma := utils.Unique([]AuthorEmail{alice, bob, zorin})

	patch, commit, parentCommit, _, err := ExtractPatch(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}
	deltas, err := NewFileDeltas(commit, parentCommit, patch)
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

func TestFileDelta_UnghostAuthors(t *testing.T) {
	setup()
	t.Cleanup(tearDown)

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

	patch, commit, parentCommit, _, err := ExtractPatch(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}
	gdeltas, err := NewFileDeltas(commit, parentCommit, patch)
	if err != nil {
		t.Error(err)
	}

	deltas := make([]TextFileDelta, len(gdeltas))
	for index, delta := range gdeltas[:3] {
		tdelta, ok := delta.(TextFileDelta)
		if !ok {
			t.Fatalf("File delta #%d is not a TextFileDelta", index)
		}
		deltas[index] = tdelta
	}

	newAuthorsAlpha := deltas[0].UnghostAuthors()
	newAuthorsBeta := deltas[1].UnghostAuthors()
	newAuthorsGamma := deltas[2].UnghostAuthors()

	utils.CompareSlices(t, expectedAlpha, newAuthorsAlpha)
	utils.CompareSlices(t, expectedBeta, newAuthorsBeta)
	utils.CompareSlices(t, expectedGamma, newAuthorsGamma)
}

func TestFileDelta_ContentFilteredForUsers(t *testing.T) {
	setup()
	t.Cleanup(tearDown)

	patch, commit, parentCommit, _, err := ExtractPatch(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}
	gdeltas, err := NewFileDeltas(commit, parentCommit, patch)
	if err != nil {
		t.Error(err)
	}

	deltas := make([]TextFileDelta, len(gdeltas))
	for index, delta := range gdeltas[:3] {
		tdelta, ok := delta.(TextFileDelta)
		if !ok {
			t.Fatalf("File delta #%d is not a TextFileDelta", index)
		}
		deltas[index] = tdelta
	}

	authors := []AuthorEmail{zorin, alice, bob, claire, david, eleonore}

	utils.CompareSlices(t, deltas[0].ContentFilteredForUsers(utils.Unique([]AuthorEmail{})), []string{"A", "C"})
	utils.CompareSlices(t, deltas[0].ContentFilteredForUsers(utils.Unique(authors[:1])), []string{"A", "C"})
	utils.CompareSlices(t, deltas[0].ContentFilteredForUsers(utils.Unique(authors[:2])), []string{"A", "C"})
	utils.CompareSlices(t, deltas[0].ContentFilteredForUsers(utils.Unique(authors[:3])), []string{"A", "b1+", "b2+", "b3+", "C"})
	utils.CompareSlices(t, deltas[0].ContentFilteredForUsers(utils.Unique(authors[:4])), []string{"A", "b1+", "b2+", "b3+", "C"})
	utils.CompareSlices(t, deltas[0].ContentFilteredForUsers(utils.Unique(authors[:5])), []string{"A", "b1+", "b2+", "b3+", "C", "d1+"})

	utils.CompareSlices(t, deltas[1].ContentFilteredForUsers(utils.Unique([]AuthorEmail{})), []string{"A", "D"})
	utils.CompareSlices(t, deltas[1].ContentFilteredForUsers(utils.Unique(authors[:1])), []string{"x1+", "x2+", "A", "D"})
	utils.CompareSlices(t, deltas[1].ContentFilteredForUsers(utils.Unique(authors[:2])), []string{"x1+", "x2+", "A", "D"})
	utils.CompareSlices(t, deltas[1].ContentFilteredForUsers(utils.Unique(authors[:3])), []string{"x1+", "x2+", "A", "b1+", "D"})
	utils.CompareSlices(t, deltas[1].ContentFilteredForUsers(utils.Unique(authors)), []string{"x1+", "x2+", "A", "b1+", "D"})

	utils.CompareSlices(t, deltas[2].ContentFilteredForUsers(utils.Unique([]AuthorEmail{})), []string{"A", "B"})
	utils.CompareSlices(t, deltas[2].ContentFilteredForUsers(utils.Unique(authors[:1])), []string{"A", "B", "c1+", "c2+"})
	utils.CompareSlices(t, deltas[2].ContentFilteredForUsers(utils.Unique(authors[:2])), []string{"A", "B", "c1+", "c2+"})
	utils.CompareSlices(t, deltas[2].ContentFilteredForUsers(utils.Unique(authors[:3])), []string{"A", "B", "c1+", "c2+"})
	utils.CompareSlices(t, deltas[2].ContentFilteredForUsers(utils.Unique(authors[:4])), []string{"A", "B", "c1+", "c2+"})
	utils.CompareSlices(t, deltas[2].ContentFilteredForUsers(utils.Unique(authors[:5])), []string{"A", "B", "c1+", "c2+"})
}

func TestNewFileDeltas(t *testing.T) {
	setup()
	t.Cleanup(tearDown)

	patch, commit, parent, _, err := ExtractPatch(repoPath, "HEAD")
	if err != nil {
		t.Error(err)
	}

	gdeltas, err := NewFileDeltas(commit, parent, patch)

	ctrText := 0
	ctrBin := 0
	for _, delta := range gdeltas {
		_, isText := delta.(TextFileDelta)
		if isText {
			ctrText++
		}

		bdelta, isBinary := delta.(BinaryFileDelta)
		if isBinary {
			ctrBin++

			length := len(bdelta.Content)
			if length != 512 {
				t.Errorf("Expected binary file to have a size of 512. Was %d", length)
			}
		}
	}

	if ctrText != 3 {
		t.Errorf("Expected %d text files. Got %d", 3, ctrText)
	}
	if ctrBin != 1 {
		t.Errorf("Expected %d binary files. Got %d", 1, ctrBin)
	}
}
