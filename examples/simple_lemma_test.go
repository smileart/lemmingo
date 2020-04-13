package lemmingo_test

import (
	"fmt"
	"path/filepath"

	"github.com/smileart/lemmingo"
)

func ExampleSimpleLemma() {
	var (
		l   string
		ok  bool
		err error
	)

	dictAbsPath, err := filepath.Abs("../dicts/en.lmm")
	if err != nil {
		panic(err)
	}

	lem, err := lemmingo.New(dictAbsPath, "", "", false, false, false)
	if err != nil {
		panic(err)
	}

	l, ok, err = lem.Lemma("i'dn't've", "PRP+MD+RB+VBP")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("stranger", "NN")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("stranger", "JJR")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("quadrillion", "NN")
	fmt.Println(l, ok, err)

	// Output:
	// i+would+not+have true <nil>
	// stranger true <nil>
	// strange true <nil>
	// quadrillion false Word's (quadrillion) lemma wasn't found!
}

func ExampleTagsetLemma() {
	var l string
	var ok bool
	var err error

	dictAbsPath, err := filepath.Abs("../dicts/en.lmm")
	if err != nil {
		panic(err)
	}

	lem, err := lemmingo.New(dictAbsPath, "en-GB", "freeling", false, false, false)
	if err != nil {
		panic(err)
	}

	l, ok, err = lem.Lemma("i'dn't've", "PRON")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("stranger", "NOUN")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("stranger", "ADJ")
	fmt.Println(l, ok, err)

	l, ok, err = lem.Lemma("quadrillion", "NOUN")
	fmt.Println(l, ok, err)

	// Output:
	// i+would+not+have true <nil>
	// stranger true <nil>
	// strange true <nil>
	// quadrillion false Word's (quadrillion) lemma wasn't found!
}
