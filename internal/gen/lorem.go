package gen

import (
	"math/rand"
	"strings"
)

var wordPool = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum",
}

func GenerateWords(n int) string {
	if n <= 0 {
		return ""
	}
	words := make([]string, n)
	for i := range words {
		words[i] = wordPool[rand.Intn(len(wordPool))]
	}
	words[0] = capitalize(words[0])
	return strings.Join(words, " ")
}

func GenerateSentences(n int) string {
	if n <= 0 {
		return ""
	}
	sentences := make([]string, n)
	for i := range sentences {
		wordCount := 8 + rand.Intn(8) // 8-15 words
		words := make([]string, wordCount)
		for j := range words {
			words[j] = wordPool[rand.Intn(len(wordPool))]
		}
		words[0] = capitalize(words[0])
		sentences[i] = strings.Join(words, " ") + "."
	}
	return strings.Join(sentences, " ")
}

func GenerateParagraphs(n int) string {
	if n <= 0 {
		return ""
	}
	paragraphs := make([]string, n)
	for i := range paragraphs {
		sentenceCount := 3 + rand.Intn(4) // 3-6 sentences
		paragraphs[i] = GenerateSentences(sentenceCount)
	}
	return strings.Join(paragraphs, "\n\n")
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
