package strategy

import (
	"fmt"
)

type Strategy int

const (
	Author = iota
	Commit
)

var str2strat map[string]Strategy

var strat2str = map[Strategy]string{
	Author: "author",
	Commit: "commit",
}

func init() {
	str2strat = make(map[string]Strategy)

	for strat, str := range strat2str {
		str2strat[str] = strat
	}
}

func (s Strategy) String() string {
	return strat2str[s]
}

func New(str string) (Strategy, error) {
	ret, ok := str2strat[str]
	if !ok {
		return -1, fmt.Errorf("%s is not a valid strategy.", str)
	}

	return ret, nil
}
