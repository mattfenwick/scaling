package loadgen

import (
	"context"
	"fmt"
)

var (
	FirstNames = []string{
		"Aziz",
		"Bob",
		"Corina",
		"Desmond",
		"Eric",
		"Frank",
		"Georgette",
		"Harvey",
		"Iris",
		"Julie",
		"Kelly",
		"Larry",
		"Monique",
		"Natalie",
		"Othello",
		"Pierre",
		"Quincy",
		"Raymondette",
		"Sarah",
		"Timmy",
		"Ulysses",
		"Virginia",
		"Will",
		"Xerxes",
		"Yannick",
		"Zebulon",
	}
	MiddleNames = []string{
		"Aardvark",
		"Crocodile",
		"Duck",
		"Groundhog",
		"Marmot",
		"Otter",
		"Squirrel",
	}
	LastNames = []string{
		"Aaronson",
		"Black",
		"White",
		"Xanthos",
		"Yastrzemski",
		"Zaborowski",
	}
)

type NameState struct {
	First     int
	Middle    int
	Last      int
	Iteration int
	Stamp     int
}

func (n *NameState) GetName() [2]string {
	first, middle, last := FirstNames[n.First], MiddleNames[n.Middle], LastNames[n.Last]
	return [2]string{
		fmt.Sprintf("%s %s %s %d (%d)", first, middle, last, n.Iteration, n.Stamp),
		fmt.Sprintf("%s.%s.%s.%d.%d@example.local", first, middle, last, n.Iteration, n.Stamp),
	}
}

func (n *NameState) Increment() {
	n.Last++
	if n.Last < len(LastNames) {
		return
	}

	n.Last = 0
	n.Middle++
	if n.Middle < len(MiddleNames) {
		return
	}

	n.Middle = 0
	n.First++
	if n.First < len(FirstNames) {
		return
	}

	n.First = 0
	n.Iteration++
}

func GenerateUsers(ctx context.Context, stamp int) <-chan [2]string {
	out := make(chan [2]string)
	go func() {
		state := &NameState{Stamp: stamp}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			out <- state.GetName()
			state.Increment()
		}
	}()
	return out
}
