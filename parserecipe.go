package parserecipe

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strings"

	"github.com/jaytaylor/html2text"
)

type WordPosition struct {
	Word     string
	Position int
}

func getWordPositions(s string, corpus []string) (wordPositions []WordPosition) {
	s = " " + s + " "
	// remove parentheses
	re := regexp.MustCompile(`(?s)\((.*)\)`)
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		s = strings.Replace(s, m[0], " ", 1)
	}

	wordPositions = []WordPosition{}
	for _, ing := range corpus {
		pos := strings.Index(s, ing)
		if pos > -1 {
			s = strings.Replace(s, ing, " ", 1)
			ing = strings.TrimSpace(ing)
			wordPositions = append(wordPositions, WordPosition{ing, pos})
		}
	}
	return
}

func GetIngredientsInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusIngredients)
}

func GetNumbersInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusNumbers)
}

func GetMeasuresInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusMeasures)
}

func SanitizeLine(s string) string {
	s = strings.ToLower(s)
	s = " " + strings.TrimSpace(s) + " "
	s = strings.Replace(s, " one ", " 1 ", -1)
	return s
}

// Parse looks for the following
// - Contains number
// - Contains mass/volume
// - Contains ingredient
// - Number occurs before ingredient
// - Number occurs before mass/volume
// - Number of ingredients is 1
// - Percent of other words is less than 50%
// - Part of list (contains - or *)
func Parse(txtFile string) (err error) {
	bFile, err := ioutil.ReadFile(txtFile)
	if err != nil {
		return
	}
	txtFileData, err := html2text.FromString(string(bFile), html2text.Options{PrettyTables: false, OmitLinks: true})
	if err != nil {
		return
	}
	txtFileData = strings.Replace(txtFileData, "1/2", "½", -1)
	txtFileData = strings.Replace(txtFileData, "1/4", "¼", -1)
	txtFileData = strings.Replace(txtFileData, "3/4", "¾", -1)
	txtFileData = strings.Replace(txtFileData, "1/8", "⅛", -1)
	txtFileData = strings.Replace(txtFileData, "3/8", "⅜", -1)
	txtFileData = strings.Replace(txtFileData, "5/8", "⅝", -1)
	txtFileData = strings.Replace(txtFileData, "7/8", "⅞", -1)
	lines := strings.Split(strings.ToLower(txtFileData), "\n")
	scores := make([]int, len(lines))

	for i, line := range lines {
		lines[i] = SanitizeLine(line)
	}

	for i, line := range lines {
		ingInString := GetIngredientsInString(line)
		numInString := GetNumbersInString(line)
		measureInString := GetMeasuresInString(line)

		score := 0
		if len(ingInString) > 0 {
			score++
		}
		if len(numInString) > 0 {
			score++
		}
		if len(measureInString) > 0 {
			score++
		}
		if len(ingInString) > 0 && len(measureInString) > 0 && ingInString[0].Position > measureInString[0].Position {
			score++
		}
		if len(ingInString) > 0 && len(numInString) > 0 && ingInString[0].Position > numInString[0].Position {
			score++
		}
		if len(measureInString) > 0 && len(numInString) > 0 && measureInString[0].Position > numInString[0].Position {
			score++
		}
		fields := strings.Fields(line)
		if len(fields) > 0 && (fields[0] == "*" || fields[0] == "-") {
			score++
		}
		scores[i] = score
	}

	start, end := GetBestTopHatPositions(scores)
	for _, line := range lines[start:end] {
		ings := GetIngredientsInString(line)
		fmt.Println(ings)
	}
	return
}

func GetBestTopHatPositions(vector []int) (start, end int) {
	vectorFloat := make([]float64, len(vector))
	for i, v := range vector {
		vectorFloat[i] = float64(v)
	}

	bestTopHatResidual := 1e9
	for i, v := range vectorFloat {
		if v == 0 {
			continue
		}
		for j, w := range vectorFloat {
			if j <= i || w == 0 {
				continue
			}
			hat := GenerateHat(len(vectorFloat), i, j, AverageFloats(vectorFloat[i:j]))
			res := CalculateResidual(vectorFloat, hat)
			if res < bestTopHatResidual {
				bestTopHatResidual = res
				start = i
				end = j
			}
		}
	}
	return
}

func CalculateResidual(fs1, fs2 []float64) float64 {
	res := 0.0
	if len(fs1) != len(fs2) {
		return -1
	}
	for i := range fs1 {
		res += math.Pow(fs1[i]-fs2[i], 2)
	}
	return res
}

func AverageFloats(fs []float64) float64 {
	f := 0.0
	for _, v := range fs {
		f += v
	}
	return f / float64(len(fs))
}

func GenerateHat(length, start, stop int, value float64) []float64 {
	f := make([]float64, length)
	for i := start; i < stop; i++ {
		f[i] = value
	}
	return f
}
