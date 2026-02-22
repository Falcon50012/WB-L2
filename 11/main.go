package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

func anagramGroups(words []string) map[string][]string {
	if len(words) == 0 {
		return map[string][]string{}
	}

	type group struct {
		first string
		words map[string]struct{}
	}

	groups := make(map[string]*group)

	for _, word := range words {
		if word == "" {
			continue
		}

		lw := strings.ToLower(word)
		sign := signature(lw)

		g, ok := groups[sign]
		if !ok {
			g = &group{
				first: lw,
				words: make(map[string]struct{}),
			}
			groups[sign] = g
		}
		g.words[lw] = struct{}{}
	}

	anagramsMap := make(map[string][]string)

	for _, g := range groups {
		if len(g.words) < 2 {
			continue
		}

		anagrams := make([]string, 0, len(g.words))

		for anagram := range g.words {
			anagrams = append(anagrams, anagram)
		}

		sort.Strings(anagrams)
		anagramsMap[g.first] = anagrams
	}

	return anagramsMap
}

func signature(s string) string {
	runes := []rune(s)
	slices.Sort(runes)

	return string(runes)
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	groups := anagramGroups(words)

	for k, v := range groups {
		fmt.Printf("%q: %q\n", k, v)
	}
}
