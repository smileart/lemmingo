package lemmingo_test

import (
	"fmt"
	"path/filepath"

	"github.com/smileart/lemmingo"
	"github.com/smileart/lemmingo/tagset"
)

func ExampleTagset() {
	dictAbsPath, err := filepath.Abs("../dicts/en.lmm")
	if err != nil {
		panic(err)
	}

	lem, err := lemmingo.New(dictAbsPath, "en-GB", "freeling", true, true, false)
	mapPos := tagset.MapPos("wordnet", "en-GB")

	pos, _ := mapPos("n")
	l, ok, e := lem.Lemma("words", pos)
	fmt.Println(l, ok, e)

	pos, _ = mapPos("v")
	l, ok, e = lem.Lemma("running", pos)
	fmt.Println(l, ok, e)

	// Output:
	// word true <nil>
	// run true <nil>
}
