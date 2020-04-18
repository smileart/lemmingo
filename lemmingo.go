// Package lemmingo implements a defensive lemmatiser/stemmer/spellchecker pipeline
package lemmingo

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Jeffail/tunny"
	"github.com/mitchellh/go-homedir"
	"github.com/otiai10/copy"
	"github.com/smileart/lemmingo/tagset"
	"github.com/tebeka/snowball"
	"github.com/trustmaster/go-aspell"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// Lemmingo is the lemmatiser and its configuration
type Lemmingo struct {
	dict            map[string]string
	stemmerLang     string
	spellerLang     string
	stemmerPool     *tunny.Pool
	spellerPool     *tunny.Pool
	stemmerFallback bool
	spellerFallback bool
	concurrent      bool
}

type fallbackPayload struct {
	Lang string
	Word string
}

// New creates a new instance of Lemmingo struct with dictionary from dictPath.
//
// It automatically generates language arguments for stemmer and speller fallbacks from langTag (BCP 47).
//
// It also instantiates one or more stemmer/speller instances depending on stemmerFallback, spellerFallback and concurrent flag.
//
// The only difference from Build method is that you don't need to figure out the language tags.
//
// It returns a pointer to Lemmingo struct.
func New(dictPath string, langTag string, tagsetName string, stemmerFallback bool, spellerFallback bool, concurrent bool) (*Lemmingo, error) {
	lang := language.Make(langTag)
	base, _ := lang.Base()

	tagsetLang := base.String()
	spellerLang := langTag
	stemmerLang := strings.ToLower(display.English.Languages().Name(base))

	return Build(dictPath, stemmerFallback, stemmerLang, spellerFallback, spellerLang, tagsetName, tagsetLang, concurrent)
}

// Build creates a new instance of Lemmingo struct according to provided option values
//
// If dictPath is a relative path the default dictionary gets first copied into the current user's home `.lemmingo` directory, and then loaded from there,
// otherwise it's loaded from absolute path provided
//
// If stemmer/speller fallbacks enabled, relevant language params will be passed to the binding libraries (NOTICE: the libs accept language options in different formats).
//
// When tagsetName provided, the Lemmingo will convert the PoS in the dictionary provided to Universal Tagset PoS and therefore all further lookups should use Universal Tagset PoS.
//
// It returns a pointer to Lemmingo struct.
func Build(dictPath string, stemmerFallback bool, stemmerLang string, spellerFallback bool, spellerLang string, tagsetName string, tagsetLang string, concurrent bool) (*Lemmingo, error) {

	var (
		l   Lemmingo
		err error
	)

	l.spellerFallback = spellerFallback
	l.stemmerFallback = stemmerFallback
	l.stemmerLang = stemmerLang
	l.spellerLang = spellerLang
	l.concurrent = concurrent

	l.dict, err = loadDict(dictPath, tagsetName, tagsetLang)
	if err != nil {
		return &l, err
	}

	// SEE: https://github.com/tebeka/snowball/issues/3
	// SEE: https://github.com/goodsign/snowball/blob/master/README.md#thread-safety
	if stemmerFallback {
		// Try to load stemmer with the given lang to fail early
		_, err = loadStemmer(l.stemmerLang)

		// If we wanted a stemmer fallback, but can't load it there's no way to recover
		if err != nil {
			panic(err)
		}

		l.stemmerPool = loadStemmerPool(1)

		if l.concurrent {
			l.stemmerPool.SetSize(runtime.NumCPU())
		} else {
			l.stemmerPool.SetSize(1)
		}
	}

	// SEE: http://aspell.net/ ("To Do" section on thread safety)
	if spellerFallback {
		// Try to load speller with the given lang to fail early
		_, err := loadSpeller(l.spellerLang)

		if err != nil {
			panic(err)
		}

		// If we wanted a speller fallback, but can't load it there's no way to recover
		l.spellerPool = loadSpellerPool(1)

		if l.concurrent {
			l.spellerPool.SetSize(runtime.NumCPU())
		} else {
			l.spellerPool.SetSize(1)
		}
	}

	return &l, nil
}

// Lemma passes the word through lemmatisation/(stemming)/(spelling) pipeline.
//
// If stemmerFallback was enabled it passes the original word to Snowball stemmer on dictionary lookup failure.
//
// If spellerFallback was enabled the stemming result gets passed to Aspell spell checker to correct the stemming issues if any.
//
// It returns lemmatised(/stemmed/spell-checked) word, and boolean flag marking if it was found in the dictionary
// (and the error if stemmer was disabled and the word wasn't found).
//
// It returns lemmatised word form from the dictionary.
func (l *Lemmingo) Lemma(word string, pos string) (string, bool, error) {
	var (
		lmm string
		err error
		ok  bool
	)

	word = strings.ToLower(word)

	// handle case when the speller is on, but stemmer is off
	if !l.stemmerFallback && l.spellerFallback {
		word = l.spellCheck(word)
	}

	// look for a word/PoS in the dict
	lmm, ok = l.dict[word+" "+strings.ToUpper(pos)]

	// if there's no word in the dict and no stemmer - fail
	if lmm == "" && !l.stemmerFallback {
		return word, false, errors.New("Word's (" + word + ") lemma wasn't found!")
	}

	// if there's no word in the dict but the stemmer is on - stem
	if lmm == "" && l.stemmerFallback {
		lmm, err = l.Stem(word)

		if err != nil {
			return word, false, err
		}
	}

	return lmm, ok, nil
}

// Stem allows to get stemma for a given word, without trying the lemmatisation dictionary first.
//
// Additionally if spellerFallback was enabled it'd pass the stemma to Aspell spell-checker for further correction.
//
// WARNING: in spelling the first suggested correction is chosen automatically, despite the fact that there could be better ones in the list.
//
// It returns stemma of the word (optionally a spell-checked one).
func (l *Lemmingo) Stem(word string) (string, error) {
	word = l.stem(word)

	if !l.spellerFallback {
		return word, nil
	}

	return l.spellCheck(word), nil
}

// Close ensures that stemmer/speller polls are cleared after the usage.
func (l *Lemmingo) Close() {
	defer l.spellerPool.Close()
	defer l.stemmerPool.Close()
}

// spellCheck gets one spell-checking goroutine from the pool, checks the spelling and if the word was misspelled, returns the first correction
func (l *Lemmingo) spellCheck(word string) string {
	return l.spellerPool.Process(fallbackPayload{
		Lang: l.spellerLang,
		Word: word,
	}).(string)
}

// stem gets one stemming goroutine from the pool, gets the stemma for a word and returns it
func (l *Lemmingo) stem(word string) string {
	return l.stemmerPool.Process(fallbackPayload{
		Lang: l.stemmerLang,
		Word: word,
	}).(string)
}

// loadDict loads the dictionary provided through dictPath in the following format: "inflected_word<space>canonical_form<space>PoS_tag"
//
// On receiving relative path: copies the files from ./dicts to $HOME/.lemmingo and treats the path as relative to that directory (suggests the existence of ./dicts)
// When tagsetName/tagsetLang provided: maps the dictionary PoS to Universal Tagset PoS
//
// It returns a map with keys in the following format: "<inflected_word>_<PoS>" with canonical form value
func loadDict(dictPath string, tagsetName string, tagsetLang string) (map[string]string, error) {
	dict := make(map[string]string)

	if !filepath.IsAbs(dictPath) {
		lemmingoHome, err := installDicts()
		if err != nil {
			return nil, err
		}

		dictPath = filepath.Join(lemmingoHome, dictPath)
	}

	file, err := os.Open(dictPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var (
		de     []string
		mapPos func(string) (string, bool)
	)

	if tagsetName != "" {
		mapPos = tagset.MapPos(tagsetName, tagsetLang)
	}

	for scanner.Scan() {
		de = strings.Split(scanner.Text(), " ")

		if mapPos != nil {
			uniPos, _ := mapPos(de[2])

			// inflected    PoS    lemma
			dict[de[0]+" "+uniPos] = de[1]
		} else {
			// inflected    PoS     lemma
			dict[de[0]+" "+de[2]] = de[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dict, nil
}

// loadStemmer creates a new Snowball stemmer for the language provided
func loadStemmer(stemmerLang string) (*snowball.Stemmer, error) {
	stemmer, err := snowball.New(stemmerLang)
	if err != nil {
		return stemmer, err
	}

	return stemmer, nil
}

// loadSpeller creates a new Aspell stemmer for the language provided
func loadSpeller(spellerLang string) (*aspell.Speller, error) {
	speller, err := aspell.NewSpeller(map[string]string{
		"lang": spellerLang,
	})
	if err != nil {
		return &speller, err
	}

	return &speller, nil
}

// loadStemmerPool creates a new goroutines pool with goLimit stemming goroutines
func loadStemmerPool(goLimit int) *tunny.Pool {
	pool := tunny.NewFunc(goLimit, func(payload interface{}) interface{} {
		args := payload.(fallbackPayload)
		stemmer, _ := loadStemmer(args.Lang)

		word := args.Word

		return stemmer.Stem(word)
	})

	return pool
}

// loadSpellerPool creates a new goroutines pool with goLimit spelling goroutines
func loadSpellerPool(goLimit int) *tunny.Pool {
	pool := tunny.NewFunc(goLimit, func(payload interface{}) interface{} {
		args := payload.(fallbackPayload)
		speller, _ := loadSpeller(args.Lang)

		word := args.Word

		if speller.Check(word) {
			return word
		}

		suggestions := speller.Suggest(word)

		if len(suggestions) == 0 {
			return word
		}

		return suggestions[0]
	})

	return pool
}

// installDicts copies the default ./dicts to current user's $HOME for DEVELOPMENT convenience.
// For production, containers and binary distributions use absolute path to the dictionary,
// so this method woudn't be called
func installDicts() (string, error) {
	_, dbg := os.LookupEnv("DEBUG")

	homePath, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	_, filename, _, _ := runtime.Caller(1)

	localDictsPath, err := filepath.Abs(filepath.Join(filepath.Dir(filename), "./dicts"))
	if err != nil {
		return "", err
	}

	lemmingoHome := filepath.Join(homePath, ".lemmingo")
	if _, err := os.Stat(lemmingoHome); !os.IsNotExist(err) {
		return lemmingoHome, nil
	}

	if dbg {
		log.Printf("Default Lemmingo dictionaries are going to be installed to: %s", lemmingoHome)
	}

	err = copy.Copy(localDictsPath, lemmingoHome)
	if err != nil {
		return "", err
	}

	return lemmingoHome, nil
}
