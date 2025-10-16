package strategy

import "testing"

var initializers []string
var expectedStrats []Strategy

func init() {
	initializers = []string{
		"author",
		"commit",
	}

	expectedStrats = []Strategy{
		Author,
		Commit,
	}
}

func TestNew(t *testing.T) {
	for ii, str := range initializers {
		sAuth, err := New(str)
		if err == nil {
			if sAuth != expectedStrats[ii] {
				t.Errorf("Wrong initialization: %s != %s", sAuth, expectedStrats[ii])
			}
		} else {
			t.Errorf("Initializing with _%s_ threw \"%s\"", str, err)
		}
	}

	_, err := New("nonexistent")
	if err == nil {
		t.Errorf("Initializing with invalid option _wrong_ did not throw an error")
	}
}

func TestString(t *testing.T) {
	for ii, strI := range expectedStrats {
		strat := strI.String()
		if strat != initializers[ii] {
			t.Errorf("Wrong Strategy name: %s != %s", strat, strI)
		}
	}
}
