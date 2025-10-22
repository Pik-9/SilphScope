package repository

type SourceLine struct {
	Content    string
	Author     string
	CommitHash string
	NewlyAdded bool
}
