package main

import (
	"strings"
	"unicode"
	"fmt"
)

func main() {
	fmt.Println(GenerateStem("unabhangig"))
}

// Implements algorithm described here:
// http://snowball.tartarus.org/algorithms/german/stemmer.html
func GenerateStem(word string) string {
	// Replace s-zed with ss
	word = strings.ToLower(word)
	word = replaceSZed(word)

	letters := []rune(word)

	// Capitalize Y's and U's in between vowels
	letters = capitalizeYsUs(letters)

	// Get regions of the word
	r1, r2 := getRegions(letters)

	// Trim suffixes
	trimmedOnce := trimSuffixStep1(letters, r1)
	trimmedTwice := trimSuffixStep2(trimmedOnce, r1)
	stem := trimDSuffix(trimmedTwice, r1, r2)

	return string(stem)
}

// Step 1
func replaceSZed(word string) string {
	return strings.Replace(word, "ß", "ss", -1)
}

// Step 2
func capitalizeYsUs(letters []rune) []rune {
	for i, max := 1, len(letters) - 1; i < max; i++ {
		if (letters[i] == 'y' || letters[i] == 'u') &&
			isBetweenVowels(i, letters) {
			letters[i] = unicode.ToUpper(letters[i])
		}
	}
	return letters
}

// Step 3
func getRegions(letters []rune) (int, int) {
	r1 := calculateRegion(letters)
	r2 := r1 + calculateRegion(letters[r1:])
	r1 = adjustR1(letters, r1) // Special requirement for German R1
	return r1, r2
}

// R1 is the region after the first consonant following a vowel, or is the null
// region at the end of the word if there is no such non-vowel
func calculateRegion(letters []rune) int {
	for i := 0; i < len(letters) - 1; i++ {
		if isVowel(letters[i]) && !isVowel(letters[i+1]) {
			return i + 2
		}
	}
	return len(letters)
}

func adjustR1(letters []rune, r1 int) int {
	if r1 >= 3 { return r1 }
	if len(letters) < 4 { return len(letters) }
	return 3
}

// Step 4:1
func trimSuffixStep1(letters []rune, r1 int) []rune {
	groupA := []string{"ern", "em", "er"}
	groupB := []string{"en", "es", "e"}
	groupC := []string{"s"}

	positions := map[string] int {
		"a": getLongestSuffix(letters, groupA),
		"b": getLongestSuffix(letters, groupB),
		"c": getLongestSuffix(letters, groupC),
	}

	// Now, get the min position (meaning, longest suffix)
	group := findMinValue(positions, len(letters))
	if (group == "a" && positions["a"] >= r1) {
		return letters[:positions["a"]]
	}

	// If an ending in Group B should be trimmed, also delete the final
	// "s" if the remaining stem ends in "niss"
	if (group == "b" && positions["b"] >= r1) {
		trimmed := letters[:positions["b"]]
		if nissPos := getSuffixPosition(trimmed, []rune("niss"));
			nissPos != -1 {
			return trimmed[:len(trimmed) - 1]
		}
	}

	if group == "c" && positions["c"] >= r1 &&
		hasValidSEnding(letters[:len(letters)-1]) {
		return letters[:positions["c"]]
	}

	return letters
}

func findMinValue(positions map[string]int, min int) (group string) {
	group = ""
	for key, value := range positions {
		if (value < min) {
			min = value
			group = key
		}
	}
	return
}
// Step 4:2
func trimSuffixStep2(letters []rune, r1 int) []rune {
	groupA := []string{"en", "er", "est"}
	groupB := []string{"st"}

	positions := map[string] int {
		"a": getLongestSuffix(letters, groupA),
		"b": getLongestSuffix(letters, groupB),
	}

	group := findMinValue(positions, len(letters))

	if group == "a" && positions["a"] >= r1 {
		return letters[:positions["a"]]
	}

	bTrimmed := letters[:positions["b"]]

	if positions["b"] >= r1 && hasValidStEnding(bTrimmed) {
		return bTrimmed
	}

	return letters
}

func trimDSuffix(letters []rune, r1 int, r2 int) []rune {
	groupA := []string{"end", "ung"}
	groupB := []string{"ig", "ik", "isch"}
	groupC := []string{"lich", "heit"}
	groupD := []string{"keit"}

	positions := map[string] int {
		"a": getLongestSuffix(letters, groupA),
		"b": getLongestSuffix(letters, groupB),
		"c": getLongestSuffix(letters, groupC),
		"d": getLongestSuffix(letters, groupD),
	}

	group := findMinValue(positions, len(letters))

	if (group == "a" && positions["a"] >= r2) {
		trimmed := letters[:positions["a"]]
		if pos := getSuffixPosition(trimmed, []rune("ig"));
			pos >= r2 && getSuffixPosition(trimmed, []rune("e")) == -1 {
				return trimmed[:pos]
		}
	}

	if (group == "b" && positions["b"] >= r2) {
		if pos := getSuffixPosition(letters[:positions["b"]], []rune ("e"));
			pos == -1 {
			return letters[:positions["b"]]
		}
	}

	if (group == "c" && positions["c"] >= r2) {
		trimmed := letters[:positions["b"]]
		erPos := getSuffixPosition(trimmed, []rune("er"))
		enPos := getSuffixPosition(trimmed, []rune("en"))
		if (erPos != -1 && erPos >= r1) || (enPos != -1 && enPos >= r1) {
			return trimmed[:len(trimmed) - 2]
		}
		return trimmed
	}

	if group == "d" && positions["d"] >= r2 {
		trimmed := letters[:positions["d"]]
		lichPos := getSuffixPosition(trimmed, []rune("lich"))
		igPos := getSuffixPosition(trimmed, []rune("ig"))
		if lichPos != -1 && lichPos >= r2 {
			return trimmed[:lichPos]
		} else if igPos != -1 && igPos >= r2 {
			return trimmed[:igPos]
		}
		return trimmed
	}

	return letters
}

// Should check for all of them, and return the largest number
// from an array of ints
func getLongestSuffix(letters []rune, group []string) int {
	pos := -1
	for i := 0; pos != -1 && i < len(group); i++ {
		suffix := []rune(group[i])
		pos = getSuffixPosition(letters, suffix)
	}
	return pos
}

func getSuffixPosition(letters, suffix []rune) int {
	if len(letters) < len(suffix) { return -1 }

	end := len(letters) - 1
	for i := len(suffix) - 1; i >= 0; i-- {
		if suffix[i] != letters[end] {
			return -1
		}
		end--
	}
	return len(letters) - len(suffix)
}


// UTILS

func runeInSlice(r rune, list []rune) bool {
	for _, el := range list {
		if el == r {
			return true
		}
	}
	return false
}

func isVowel(letter rune) bool {
	vowels := []rune("aeiouyäöü")
	return runeInSlice(letter, vowels)
}

func isBetweenVowels(i int, letters []rune) bool {
	if i == 0 || i == len(letters) {
		return false
	}

	return isVowel(letters[i-1]) && isVowel(letters[i+1])
}

func hasValidSEnding(letters []rune) bool {
	ending := letters[len(letters)-1]
	sEndings := []rune("bdfghklmnrt")
	return runeInSlice(ending, sEndings)
}

func hasValidStEnding(letters []rune) bool {
	ending := letters[len(letters)-1]
	stEndings := []rune("bdfghklmnt")
	return runeInSlice(ending, stEndings)
}
