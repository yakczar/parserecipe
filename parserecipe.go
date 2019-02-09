package parserecipe

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/jaytaylor/html2text"
	colorable "github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
	log.SetLevel(logrus.DebugLevel)
}

type WordPosition struct {
	Word     string
	Position int
}

func getWordPositions(s string, corpus []string) (wordPositions []WordPosition) {
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

func GetOtherInBetweenPositions(s string, pos1, pos2 WordPosition) (other string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()
	other = s[pos1.Position+len(pos1.Word)+1 : pos2.Position]
	other = strings.TrimSpace(other)
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

	// remove parentheses
	re := regexp.MustCompile(`(?s)\((.*)\)`)
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		s = strings.Replace(s, m[0], " ", 1)
	}

	s = " " + strings.TrimSpace(s) + " "

	// replace unicode fractions with fractions
	for v := range corpusFractionNumberMap {
		s = strings.Replace(s, v, corpusFractionNumberMap[v].fractionString, -1)
	}

	// remove non-alphanumeric
	reg, _ := regexp.Compile("[^a-zA-Z0-9/]+")
	s = reg.ReplaceAllString(s, " ")

	// replace fractions with unicode fractions
	for v := range corpusFractionNumberMap {
		s = strings.Replace(s, corpusFractionNumberMap[v].fractionString, v, -1)
	}

	s = strings.Replace(s, " one ", " 1 ", -1)
	s = strings.Replace(s, " egg ", " eggs ", -1)

	return s
}

type LineInfo struct {
	LineOriginal        string
	Line                string         `json:",omitempty"`
	IngredientsInString []WordPosition `json:",omitempty"`
	AmountInString      []WordPosition `json:",omitempty"`
	MeasureInString     []WordPosition `json:",omitempty"`
	Ingredient          Ingredient     `json:",omitempty"`
}

type Parsed struct {
	Lines []LineInfo
}

type Ingredient struct {
	Name              string  `json:",omitempty"`
	Comment           string  `json:",omitempty"`
	MeasureOriginal   Measure `json:",omitempty"`
	MeasureConverted  Measure `json:",omitempty"`
	MeasureNormalized Measure `json:",omitempty"`
}

type Measure struct {
	Amount float64
	Name   string
}

func ParseDirections(lis []LineInfo) (rerr error) {
	log.Debug(len(lis))
	scores := make([]float64, len(lis))
	for i, li := range lis {
		if i > 30 {
			break
		}
		if len(strings.TrimSpace(li.Line)) < 3 {
			continue
		}
		score := 0.0
		for _, corpusDirection := range corpusDirections {
			if strings.Contains(li.Line, corpusDirection) {
				score++
			}
		}
		if len(li.Line) < 5 {
			score = 0
		}
		scores[i] = score
	}

	start, end := GetBestTopHatPositions(scores)
	log.Debugf("direction are from line %d to %d", start, end)
	directionI := 1
	for i := start; i <= end; i++ {
		if len(strings.TrimSpace(lis[i].Line)) == 0 {
			continue
		}
		direction := strings.TrimSpace(lis[i].LineOriginal)
		if string(direction[0]) == string("*") {
			direction = strings.TrimSpace(direction[1:])
		}

		if len(strings.Fields(direction)) < 5 {
			continue
		}
		log.Debugf("%d) %s", directionI, direction)
		directionI++
	}
	return
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
func Parse(txtFile string) (parsed Parsed, rerr error) {
	bFile, rerr := ioutil.ReadFile(txtFile)
	if rerr != nil {
		return
	}
	txtFileData, rerr := html2text.FromString(string(bFile), html2text.Options{PrettyTables: false, OmitLinks: true})
	if rerr != nil {
		return
	}

	lines := strings.Split(txtFileData, "\n")
	scores := make([]float64, len(lines))
	lineInfos := make([]LineInfo, len(lines))
	i := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		i++
		lineInfos[i].LineOriginal = line
		line = SanitizeLine(line)
		lineInfos[i].Line = line
		lineInfos[i].IngredientsInString = GetIngredientsInString(line)
		lineInfos[i].AmountInString = GetNumbersInString(line)
		lineInfos[i].MeasureInString = GetMeasuresInString(line)

		score := 0.0
		if len(lineInfos[i].IngredientsInString) > 0 {
			score++
		}
		if len(lineInfos[i].AmountInString) > 0 {
			score++
		}
		if len(lineInfos[i].MeasureInString) > 0 {
			score++
		}
		if len(lineInfos[i].IngredientsInString) > 0 && len(lineInfos[i].MeasureInString) > 0 && lineInfos[i].IngredientsInString[0].Position > lineInfos[i].MeasureInString[0].Position {
			score++
		}
		if len(lineInfos[i].IngredientsInString) > 0 && len(lineInfos[i].AmountInString) > 0 && lineInfos[i].IngredientsInString[0].Position > lineInfos[i].AmountInString[0].Position {
			score++
		}
		if len(lineInfos[i].MeasureInString) > 0 && len(lineInfos[i].AmountInString) > 0 && lineInfos[i].MeasureInString[0].Position > lineInfos[i].AmountInString[0].Position {
			score++
		}
		if score > 0 && len(lineInfos[i].LineOriginal) > 100 {
			score--
		}
		fields := strings.Fields(line)
		if len(fields) > 0 && (fields[0] == "*" || fields[0] == "-") {
			score++
		}
		if score == 1 {
			score = 0.0
		}
		// log.Debugf("'%s' (%d)", line, score)
		scores[i] = score
	}
	scores = scores[:i+1]
	lineInfos = lineInfos[:i+1]

	lines = make([]string, len(lineInfos))
	for i, li := range lineInfos {
		lines[i] = li.Line
	}
	ioutil.WriteFile("out", []byte(strings.Join(lines, "\n")), 0644)

	start, end := GetBestTopHatPositions(scores)

	ParseDirections(lineInfos[end:])

	parsed = Parsed{[]LineInfo{}}
	for _, lineInfo := range lineInfos[start-3 : end+3] {
		if len(strings.TrimSpace(lineInfo.Line)) < 3 {
			continue
		}

		// get amount, continue if there is an error
		err := lineInfo.getTotalAmount()
		if err != nil {
			log.WithFields(logrus.Fields{
				"line": strings.TrimSpace(lineInfo.LineOriginal),
			}).Errorf("%s", err.Error())
			continue
		}

		// get ingredient, continue if its not found
		err = lineInfo.getIngredient()
		if err != nil {
			log.WithFields(logrus.Fields{
				"line": strings.TrimSpace(lineInfo.LineOriginal),
			}).Errorf("%s", err.Error())
			continue
		}

		// get measure
		err = lineInfo.getMeasure()
		if err != nil {
			log.WithFields(logrus.Fields{
				"line": strings.TrimSpace(lineInfo.LineOriginal),
			}).Errorf("%s", err.Error())
		}

		// get comment
		if len(lineInfo.MeasureInString) > 0 && len(lineInfo.IngredientsInString) > 0 {
			lineInfo.Ingredient.Comment = GetOtherInBetweenPositions(lineInfo.Line, lineInfo.MeasureInString[0], lineInfo.IngredientsInString[0])
		}

		log.WithFields(logrus.Fields{
			"line": strings.TrimSpace(lineInfo.LineOriginal),
		}).Debugf("%s (%s): %+v", lineInfo.Ingredient.Name, lineInfo.Ingredient.Comment, lineInfo.Ingredient.MeasureOriginal)

		parsed.Lines = append(parsed.Lines, lineInfo)
	}

	return
}

func (lineInfo *LineInfo) getTotalAmount() (err error) {
	lastPosition := -1
	totalAmount := 0.0
	wps := lineInfo.AmountInString
	for i := range wps {
		wps[i].Word = strings.TrimSpace(wps[i].Word)
		if lastPosition == -1 {
			totalAmount = convertStringToNumber(wps[i].Word)
		} else if math.Abs(float64(wps[i].Position-lastPosition)) < 2 {
			totalAmount += convertStringToNumber(wps[i].Word)
		}
		lastPosition = wps[i].Position
	}
	if totalAmount == 0 && strings.Contains(lineInfo.Line, "whole") {
		totalAmount = 1
	}
	if totalAmount == 0 {
		err = fmt.Errorf("no amount found")
	} else {
		lineInfo.Ingredient.MeasureOriginal.Amount = totalAmount
	}
	return
}

func (lineInfo *LineInfo) getIngredient() (err error) {
	if len(lineInfo.IngredientsInString) == 0 {
		err = fmt.Errorf("no ingredient found")
		return
	}
	lineInfo.Ingredient.Name = lineInfo.IngredientsInString[0].Word
	return
}

func (lineInfo *LineInfo) getMeasure() (err error) {
	if len(lineInfo.MeasureInString) == 0 {
		lineInfo.Ingredient.MeasureOriginal.Name = "whole"
		return
	}
	lineInfo.Ingredient.MeasureOriginal.Name = lineInfo.MeasureInString[0].Word
	return
}

func GetBestTopHatPositions(vectorFloat []float64) (start, end int) {

	bestTopHatResidual := 1e9
	for i, v := range vectorFloat {
		if v < 2 {
			continue
		}
		for j, w := range vectorFloat {
			if j <= i || w < 1 {
				continue
			}
			hat := GenerateHat(len(vectorFloat), i, j, AverageFloats(vectorFloat[i:j]))
			res := CalculateResidual(vectorFloat, hat) / float64(len(vectorFloat))
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

func convertStringToNumber(s string) float64 {
	switch s {
	case "½":
		return 0.5
	case "¼":
		return 0.25
	case "¾":
		return 0.75
	case "⅛":
		return 1.0 / 8
	case "⅜":
		return 3.0 / 8
	case "⅝":
		return 5.0 / 8
	case "⅞":
		return 7.0 / 8
	case "⅔":
		return 2.0 / 3
	case "⅓":
		return 1.0 / 3
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
