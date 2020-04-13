package lemmingo_test

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/smileart/lemmingo"
	"github.com/zhexuany/wordGenerator"
)

type LemmaTestCase struct {
	word  string
	pos   string
	lemma string
}

type StemTestCase struct {
	word  string
	stem  string
	spell string
}

var lemmaCases = []LemmaTestCase{
	{
		word:  "am",
		pos:   "vbp",
		lemma: "be",
	},
	{
		word:  "Are",
		pos:   "vbp",
		lemma: "be",
	},
	{
		word:  "caresses",
		pos:   "NNS",
		lemma: "caress",
	},
	{
		word:  "operational",
		pos:   "JJ",
		lemma: "operational",
	},
	{
		word:  "marketing",
		pos:   "nn",
		lemma: "marketing",
	},
	{
		word:  "abandoning",
		pos:   "VBG",
		lemma: "abandon",
	},
	{
		word:  "abracadabra",
		pos:   "NN",
		lemma: "abracadabra",
	},
}

var stemCases = []StemTestCase{
	{
		word:  "laboratory",
		stem:  "laboratori",
		spell: "laboratory",
	},
	{
		word:  "teenager",
		stem:  "teenag",
		spell: "teenage",
	},
	{
		word:  "ababagalamaga",
		stem:  "ababagalamaga",
		spell: "ababagalamaga",
	},
}

var l = loadLemmingo(false, false, false)
var ls = loadLemmingo(true, false, false)
var lss = loadLemmingo(true, true, false)
var lssc = loadLemmingo(true, true, true)

func loadLemmingo(withStemmer bool, withSpeller bool, concurrent bool) *lemmingo.Lemmingo {
	abs, err := filepath.Abs("./dicts/en.lmm")
	if err != nil {
		panic(err)
	}

	lem, err := lemmingo.New(abs, "en-US", "", withStemmer, withSpeller, concurrent)
	if err != nil {
		panic(err)
	}

	return lem
}

// =============== Tests ===============

func TestNew(t *testing.T) {
	_, err := lemmingo.New("./en.lmm", "en-US", "penn", false, false, false)
	if err != nil {
		t.Error(err)
	}
}

func TestNewWrongTagset(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on getting wrong tagset!")
		}
	}()

	_, _ = lemmingo.New("./en.lmm", "en-US", "wrong", false, false, false)
}

func TestNewWrongTagsetLang(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			l.Close()
			t.Errorf("The code did not panic on getting wrong tagset language!")
		}
	}()

	_, _ = lemmingo.New("./en.lmm", "pt-BR", "freeling", false, false, false)
}

func TestNewWrongRelativePath(t *testing.T) {
	_, err := lemmingo.New("uk.lmm", "pt-BR", "freeling", false, false, false)

	if err == nil {
		//open /Users/user/.lemmingo/uk.lmm: no such file or directory
		t.Errorf("The method was supposed to return an error on wrong relative path to the dictionary!")
	}
}

func TestNewWrongAbsolutePath(t *testing.T) {
	_, err := lemmingo.New("/tmp/en.lmm", "pt-BR", "freeling", false, false, false)

	if err == nil {
		//open /tmp/en.lmm: no such file or directory
		t.Errorf("The method was supposed to return an error on wrong relative path to the dictionary!")
	}
}

// NOTICE: Currently similar test for the speller would panic with the SIGABRT, core dump and CGO stack trace
// SEE: https://github.com/trustmaster/go-aspell/issues/1
func TestBuildWrongStemmerLang(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			l.Close()
			t.Errorf("The code did not panic on getting wrong stemmer language!")
		}
	}()

	l, _ = lemmingo.Build("./en.lmm", true, "tokipona", false, "", "", "", false)
}

func TestLemma(t *testing.T) {
	for _, c := range lemmaCases {
		lmm, _, e := l.Lemma(c.word, c.pos)
		if lmm != c.lemma {
			t.Error(e)
			t.Errorf("For the word '%s' we've got: '%s' lemma, expected: '%s'.", c.word, lmm, c.lemma)
		}
	}
}

func TestStem(t *testing.T) {
	for _, c := range stemCases {
		st, _ := ls.Stem(c.word)
		if st != c.stem {
			t.Errorf("For the word '%s' we've got: '%s' stem, expected: '%s'.", c.word, st, c.stem)
		}
	}
}

func TestStemWithSpell(t *testing.T) {
	for _, c := range stemCases {
		st, _ := lss.Stem(c.word)
		if st != c.spell {
			t.Errorf("For the word '%s' we've got: '%s' stem after spellcheck, expected: '%s'.", c.word, st, c.spell)
		}
	}
}

func TestInGoroutines(t *testing.T) {
	defer lssc.Close()

	var wg sync.WaitGroup

	words := wordGenerator.GetWords(10000, 20)
	lemmas := make([]*string, len(words))

	for i, word := range words {
		wg.Add(1)

		go func(i int, word string) {
			defer wg.Done()

			lemma, _, _ := lssc.Lemma(word, "POS")
			lemmas[i] = &lemma
		}(i, word)
	}

	wg.Wait()
}

// ============ Benchmarks =============

func BenchmarkBareLemmingo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = loadLemmingo(false, false, false)
	}
}

func BenchmarkStemmerLemmingo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = loadLemmingo(true, false, false)
	}
}

func BenchmarkStemmerSpellerLemmingo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = loadLemmingo(true, true, false)
	}
}

func BenchmarkStemmerSpellerConcurrentLemmingo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = loadLemmingo(true, true, true)
	}
}

func BenchmarkLemma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l.Lemma("weirdest", "JJS")
	}
}

func BenchmarkStemmerLemma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ls.Lemma("weirdest", "NONE")
	}
}

func BenchmarkStemmerSpellerLemma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lss.Lemma("weirdest", "NONE")
	}
}

func BenchmarkStemmerSpellerConcurrentLemma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lssc.Lemma("weirdest", "NONE")
	}
}
