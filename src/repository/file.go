package repository

type SourceFile struct {
	FilePath string
	Lines    []SourceLine
}

func (f *SourceFile) Filtered(lineFilter func(*SourceLine) bool) SourceFile {
	ret := SourceFile{
		FilePath: f.FilePath,
		Lines:    make([]SourceLine, 0, len(f.Lines)),
	}

	for _, line := range f.Lines {
		ret.Lines = append(ret.Lines, SourceLine{
			Content:    line.Content,
			Author:     line.Author,
			CommitHash: line.CommitHash,
			NewlyAdded: line.NewlyAdded && lineFilter(&line),
		})
	}

	return ret
}

func (f *SourceFile) OldLines() []string {
	ret := make([]string, 0)
	for _, line := range f.Lines {
		if !line.NewlyAdded {
			ret = append(ret, line.Content)
		}
	}

	return ret
}

func (f *SourceFile) AllLines() []string {
	ret := make([]string, len(f.Lines))
	for idx, line := range f.Lines {
		ret[idx] = line.Content
	}
	return ret
}

func (f *SourceFile) GetAllAuthors() []string {
	authors := make(map[string]bool)
	for _, line := range f.Lines {
		authors[line.Author] = true
	}

	ret := make([]string, 0, len(authors))
	for k := range authors {
		ret = append(ret, k)
	}

	return ret
}

func (f *SourceFile) AuthorSeries(authors []string) []SourceFile {
	ret := make([]SourceFile, 0, len(authors))
	tmp := f
	for _, author := range authors {
		sf := tmp.Filtered(func(line *SourceLine) bool {
			return line.Author == author
		})
		ret = append(ret, sf)
		tmp = &sf
	}
	return ret
}
