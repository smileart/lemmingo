package lemmingo_test

import (
	"fmt"
	"github.com/smileart/lemmingo"
)

func ExampleLemmaWithSpellcheckingOnly() {
	var l string
	var ok bool
	var err error

	lem, err := lemmingo.New("./en.lmm", "en", "freeling", false, true, false)
	if err != nil {
		panic(err)
	}

	// NOTICE: ONLY in case when the spellchecker is turned on, but stemmer is not
	// we get the original word spellchecked instead of the stemmer results!
	l, ok, err = lem.Lemma("teeenager", "NOUN")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("bubling", "VERB")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("lovinh", "VERB")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("loveing", "VERB")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("typo", "NOUN")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("juse", "NOUN")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("abracadabrated", "ADJ")
	fmt.Println(l, ok, err)

	// Output:
	// teenager true <nil>
	// bubble true <nil>
	// love true <nil>
	// love true <nil>
	// typo true <nil>
	// Jude false Word's (Jude) lemma wasn't found!
	// abracadabra ted false Word's (abracadabra ted) lemma wasn't found!
}
