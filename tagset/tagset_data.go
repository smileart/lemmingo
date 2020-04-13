// Package tagset implements tagset mappings to Universal Part-of-Speech Tagset
// REF: https://arxiv.org/abs/1104.2086
// REF: http://petrovi.de/data/universal.pdf
package tagset

// Penn Treebank Project - EN, consolidated among multiple sources
// REF: https://www.ling.upenn.edu/courses/Fall_2003/ling001/penn_treebank_pos.html
// REF: https://www.clips.uantwerpen.be/pages/mbsp-tags
// REF: https://github.com/jdkato/prose
var pennEn = map[string]string{
	"CC":   "CONJ", // coordinating conjunction - and, or, but
	"CD":   "NUM",  // cardinal number -	1, third, five, three, 13%
	"DT":   "DET",  // determiner - the, a, these
	"EX":   "DET",  // existential there - there is, there were six girls
	"FW":   "X",    // foreign word - les, mais
	"IN":   "ADP",  // preposition, subordinating conjunction - in, like, of, on, before, unless
	"JJ":   "ADJ",  // adjective - green, nice, easy
	"JJR":  "ADJ",  // adjective, comparative - greener, nicer, easier
	"JJS":  "ADJ",  // adjective, superlative - greenest, nicest, easiest
	"LS":   "X",    // list marker - 1), 7., *
	"MD":   "VERB", // modal - could, will, may, should
	"NN":   "NOUN", // noun, singular or mass - apple, tiger, chair, laughter
	"NNS":  "NOUN", // noun plural - apples, tigers, chairs, insects
	"NNP":  "NOUN", // proper noun, singular - Germany, God, Alice, Serge
	"NNPS": "NOUN", // proper noun, plural - Christmases, Vikings
	"PDT":  "DET",  // predeterminer - both his children
	"POS":  "PRT",  // possessive ending - friend’s, Serge's
	"PRP":  "PRON", // personal pronoun - I, she, me, you, it
	"PRP$": "PRON", // pronoun, possessive - my, your, our
	"RB":   "ADV",  // adverb - extremely, loudly, hard
	"RBR":  "ADV",  // adverb, comparative - better
	"RBS":  "ADV",  // adverb, superlative - best
	"RP":   "PRT",  // adverb, particle - about, off, up
	"SYM":  "X",    // symbol - % © / [ = *
	"TO":   "PRT",  // infinitival to - to do, to go
	"UH":   "X",    // interjection - oh, oops, gosh
	"VB":   "VERB", // verb, base form - think
	"VBZ":  "VERB", // verb, 3rd person singular present - she thinks
	"VBP":  "VERB", // verb, non-3rd person singular present - I think
	"VBD":  "VERB", // verb, past tense - they thought
	"VBN":  "VERB", // verb, past participle - sunk
	"VBG":  "VERB", // verb, gerund or present participle - thinking is fun
	"WDT":  "DET",  // wh-determiner - which, whatever, whichever
	"WP":   "PRON", // wh-pronoun, personal - what, who, whom
	"WP$":  "PRON", // wh-pronoun, possessive - whose, whosever
	"WRB":  "ADV",  // wh-adverb - where, when
	"(":    ".",    // left round bracket
	")":    ".",    // right round bracket
	",":    ".",    // comma
	":":    ".",    // colon
	".":    ".",    // period
	"''":   ".",    // closing quotation mark
	"``":   ".",    // opening quotation mark
	"#":    ".",    // number sign
	"$":    ".",    // currency
}

// FreeLing, an open source language analysis tool suite - EN
// REF: http://nlp.lsi.upc.edu/freeling/
// REF: https://github.com/TALP-UPC/freeling
// REF: https://freeling-user-manual.readthedocs.io/en/latest/tagsets/tagset-en/
var freelingEn = map[string]string{
	"CC":            "CONJ",
	"DT":            "DET",
	"DT+MD":         "DET",
	"DT+VB":         "DET",
	"DT+VBD/MD":     "DET",
	"EX":            "DET",
	"IN":            "ADP",
	"JJ":            "ADJ",
	"JJR":           "ADJ",
	"JJS":           "ADJ",
	"MD":            "VERB",
	"MD+RB":         "VERB",
	"MD+RB+VBP":     "VERB",
	"MD+VB":         "VERB",
	"NN":            "NOUN",
	"NNS":           "NOUN",
	"POS":           "PRT",
	"PRP":           "PRON",
	"PRP$":          "PRON",
	"PRP+DT":        "PRON",
	"PRP+MD":        "PRON",
	"PRP+MD+RB":     "PRON",
	"PRP+MD+RB+VBP": "PRON",
	"PRP+MD+VBP":    "PRON",
	"PRP+VB":        "PRON",
	"PRP+VBD/MD":    "PRON",
	"PRP+VBP":       "PRON",
	"RB":            "ADV",
	"RB+MD":         "ADV",
	"RB+VBD/MD":     "ADV",
	"RB+VBZ":        "ADV",
	"RBR":           "ADV",
	"RBS":           "ADV",
	"RP":            "PRT",
	"TO":            "PRT",
	"UH":            "X",
	"VB":            "VERB",
	"VB+PRP":        "VERB",
	"VB+RB":         "VERB",
	"VBD":           "VERB",
	"VBD+RB":        "VERB",
	"VBG":           "VERB",
	"VBN":           "VERB",
	"VBP":           "VERB",
	"VBP+RB":        "VERB",
	"VBZ":           "VERB",
	"VBZ+RB":        "VERB",
	"WDT":           "DET",
	"WP":            "PRON",
	"WP$":           "PRON",
	"WP+MD":         "PRON",
	"WP+MD+VBP":     "PRON",
	"WP+VB":         "PRON",
	"WP+VBD/MD":     "PRON",
	"WP+VBP":        "PRON",
	"WRB":           "ADV",
	"WRB+VB":        "ADV",
}
