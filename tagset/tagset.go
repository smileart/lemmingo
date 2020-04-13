// Package tagset provides mappings from various tagsets to Universal Tagset
package tagset

import (
	"errors"

	"golang.org/x/text/language"
)

// tagsetMapping names, and maps supported tagsets and returns the one for the languageTag provided.
func tagsetMapping(tagsetName string, languageTag string) (map[string]string, error) {
	closureMap := map[string]map[string]string{
		"penn_en":     pennEn,     // Penn Treebank Project, EN
		"freeling_en": freelingEn, // FreeLing Project, EN
	}

	mappingKey := tagsetName + "_" + languageTag
	mapping, ok := closureMap[mappingKey]

	if !ok {
		return mapping, errors.New("Tagset mapping for `" + mappingKey + "` was not found!")
	}

	return mapping, nil
}

// MapPos gets the mapping for the tagsetName/languageTag provided and returns a function to access PoS mappings.
func MapPos(tagsetName string, languageTag string) func(posTag string) (string, bool) {
	lang := language.Make(languageTag)
	base, _ := lang.Base()

	languageTag = base.String()
	tagset, err := tagsetMapping(tagsetName, languageTag)

	// if we wanted a certain mapping, but can't load one there's no way to recover
	if err != nil {
		panic(err)
	}

	return func(posTag string) (string, bool) {
		val, ok := tagset[posTag]

		return val, ok
	}
}
