package lemmingo_test

import (
	"fmt"
	"github.com/smileart/lemmingo"
)

func ExampleLemmaWithStemming() {
	var l string
	var ok bool
	var err error

	lem, err := lemmingo.New("./en.lmm", "en-US", "freeling", true, false, false)
	if err != nil {
		panic(err)
	}

	l, ok, err = lem.Lemma("teenager", "NONSENSE")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("bubbling", "NONSENSE")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("loving", "NONSENSE")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("loveing", "ADJ")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("abracadabrated", "ADJ")
	fmt.Println(l, ok, err)

	// Output:
	// teenag false <nil>
	// bubbl false <nil>
	// love false <nil>
	// love false <nil>
	// abracadabr false <nil>
}
