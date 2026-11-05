package core

import (
	"strings"
	"unicode"
)

// NaturalLess compara dos strings ordenando números de forma natural
// (ej. "item_2.txt" precede a "item_10.txt").
func NaturalLess(a, b string) bool {
	x, y := []rune(a), []rune(b)
	i, j, zeroTie := 0, 0, 0
	digit := func(r rune) bool { return r >= '0' && r <= '9' }

	for i < len(x) && j < len(y) {
		if digit(x[i]) && digit(y[j]) {
			startI, startJ := i, j
			for i < len(x) && digit(x[i]) {
				i++
			}
			for j < len(y) && digit(y[j]) {
				j++
			}
			valueI, valueJ := startI, startJ
			for valueI < i && x[valueI] == '0' {
				valueI++
			}
			for valueJ < j && y[valueJ] == '0' {
				valueJ++
			}
			if i-valueI != j-valueJ {
				return i-valueI < j-valueJ
			}
			if cmp := strings.Compare(string(x[valueI:i]), string(y[valueJ:j])); cmp != 0 {
				return cmp < 0
			}
			if zeroTie == 0 {
				zeroTie = (i - startI) - (j - startJ)
			}
		} else {
			left, right := unicode.ToLower(x[i]), unicode.ToLower(y[j])
			if left != right {
				return left < right
			}
			i++
			j++
		}
	}
	if i != len(x) || j != len(y) {
		return len(x)-i < len(y)-j
	}
	if zeroTie != 0 {
		return zeroTie < 0
	}
	return a < b
}
