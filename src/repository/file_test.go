package repository

import (
	"fmt"
	"testing"
)

var fileA SourceFile
var fileB SourceFile

func assumeListsEqual[K comparable](t *testing.T, a []K, b []K) {
	if len(a) != len(b) {
		t.Errorf("Sizes don't match: %d != %d", len(a), len(b))
	}

	for idx, elem := range a {
		if elem != b[idx] {
			t.Error("Position", idx, "|", elem, "!=", b[idx])
		}
	}
}

func init() {
	fileA = SourceFile{
		FilePath: "A.txt",
		Lines: []SourceLine{
			{
				Content:    "Bulbasaur",
				Author:     "Alice",
				CommitHash: "aaa",
				NewlyAdded: true,
			},
			{
				Content:    "Charmander",
				Author:     "David",
				CommitHash: "ddd",
				NewlyAdded: true,
			},
			{
				Content:    "Squirtle",
				Author:     "Alice",
				CommitHash: "aaa",
				NewlyAdded: true,
			},
		},
	}

	fileB = SourceFile{
		FilePath: "B.txt",
		Lines: []SourceLine{
			{
				Content:    "Sun",
				Author:     "Gaia",
				CommitHash: "000",
				NewlyAdded: false,
			},
			{
				Content:    "Mercury",
				Author:     "Alice",
				CommitHash: "aaa",
				NewlyAdded: true,
			},
			{
				Content:    "Venus",
				Author:     "Bob",
				CommitHash: "bbb",
				NewlyAdded: true,
			},
			{
				Content:    "Earth",
				Author:     "Charles",
				CommitHash: "ccc",
				NewlyAdded: true,
			},
		},
	}
}

func TestOldlines(t *testing.T) {
	if len(fileA.OldLines()) != 0 {
		t.Error("Expected fileA to be empty, instead got", fileA.OldLines())
	}

	if len(fileB.OldLines()) != 1 {
		t.Error("Expected fileB to contain one entry, instead got", fileB.OldLines())
	}

	if fileB.OldLines()[0] != "Sun" {
		t.Error("Expected fileB's only old line to be \"Sun\", instead got", fileB.OldLines())
	}
}

func TestAllLines(t *testing.T) {
	expA := []string{
		"Bulbasaur",
		"Charmander",
		"Squirtle",
	}

	expB := []string{
		"Sun",
		"Mercury",
		"Venus",
		"Earth",
	}

	assumeListsEqual(t, expA, fileA.AllLines())
	assumeListsEqual(t, expB, fileB.AllLines())
}

func TestGetAllAuthors(t *testing.T) {
	expA := []string{
		"Alice",
		"David",
	}

	expB := []string{
		"Gaia",
		"Alice",
		"Bob",
		"Charles",
	}

	authA := fileA.GetAllAuthors()
	authB := fileB.GetAllAuthors()

	assumeListsEqual(t, expA, authA)
	assumeListsEqual(t, expB, authB)
}

func TestAuthorSeries(t *testing.T) {
	authors1 := []string{"David", "Alice"}
	seriesA1 := fileA.AuthorSeries(authors1)
	for _, elem := range seriesA1 {
		fmt.Println(elem)
	}
	assumeListsEqual(t, []string{"Charmander"}, seriesA1[0].OldLines())
	assumeListsEqual(t, []string{"Bulbasaur", "Charmander", "Squirtle"}, seriesA1[1].OldLines())
}
