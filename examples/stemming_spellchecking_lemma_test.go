package lemmingo_test

import (
	"fmt"
	"github.com/smileart/lemmingo"
	"github.com/zhexuany/wordGenerator"
	"sync"
)

func ExampleLemmaWithSpellchecking() {
	var l string
	var ok bool
	var err error

	lem, err := lemmingo.New("./en.lmm", "en", "freeling", true, true, false)
	if err != nil {
		panic(err)
	}

	// WARNING: if "teenager" with PoS "NONSENSE" with stemmer & speller turned on
	// it becomes "teenage" (which is ADJ, not NOUN)
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
	// teenage false <nil>
	// bubble false <nil>
	// love false <nil>
	// love false <nil>
	// abracadabra false <nil>
}

func ExampleLemmaWithSpellcheckingGoroutines() {
	var wg sync.WaitGroup

	words := wordGenerator.GetWords(50, 20)
	words[42] = "concurrency"
	lemmas := make([]string, len(words))

	l, err := lemmingo.New("./en.lmm", "en", "freeling", true, true, true)
	if err != nil {
		panic(err)
	}

	for i, word := range words {
		wg.Add(1)

		go func(i int, word string) {
			defer wg.Done()

			lmm, _, _ := l.Lemma(word, "NOUN")
			lemmas[i] = lmm
		}(i, word)
	}

	wg.Wait()
	fmt.Println(lemmas[42])

	// Output:
	// concurrency
}
