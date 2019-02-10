package parserecipe

//go:generate go run corpus/main.go
//go:generate gofmt -s -w corpus.go

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

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
	sort.Slice(wordPositions, func(i, j int) bool {
		return wordPositions[i].Position < wordPositions[j].Position
	})
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

type Recipe struct {
	FileName    string
	FileContent string
	Lines       []LineInfo
	Directions  []string
}

func NewFromFile(fname string) (r *Recipe, err error) {
	r = &Recipe{FileName: fname}
	_, err = os.Stat(fname)
	return
}

func NewFromURL(url string) (r *Recipe, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	r = &Recipe{FileName: url}
	r.FileContent, err = html2text.FromString(string(html), html2text.Options{PrettyTables: false, OmitLinks: true})

	return
}

type Ingredient struct {
	Name              string   `json:",omitempty"`
	Comment           string   `json:",omitempty"`
	MeasureOriginal   *Measure `json:",omitempty"`
	MeasureConverted  *Measure `json:",omitempty"`
	MeasureNormalized *Measure `json:",omitempty"`
}

type Measure struct {
	Amount float64
	Name   string
}

type IngredientList struct {
	Ingredients []Ingredient
}

func (r *Recipe) ParseDirections(lis []LineInfo) (rerr error) {
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
		for _, corpusDirection := range corpusDirectionsNeg {
			if strings.Contains(li.Line, corpusDirection) {
				score--
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
	r.Directions = []string{}
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
		r.Directions = append(r.Directions, direction)
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
func (r *Recipe) Parse() (rerr error) {
	if r.FileContent == "" && r.FileName != "" {
		var bFile []byte
		bFile, rerr = ioutil.ReadFile(r.FileName)
		if rerr != nil {
			return
		}
		r.FileContent, rerr = html2text.FromString(string(bFile), html2text.Options{PrettyTables: false, OmitLinks: true})
		if rerr != nil {
			return
		}
	}

	lines := strings.Split(r.FileContent, "\n")
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

	var wg sync.WaitGroup
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := r.ParseDirections(lineInfos[end:])
		if err != nil {
			log.Warn(err.Error())
		}
	}(&wg)

	r.Lines = []LineInfo{}
	for _, lineInfo := range lineInfos[start-3 : end+3] {
		if len(strings.TrimSpace(lineInfo.Line)) < 3 {
			continue
		}

		lineInfo.Ingredient.MeasureOriginal = &Measure{}

		// get amount, continue if there is an error
		err := lineInfo.getTotalAmount()
		if err != nil {
			log.WithFields(logrus.Fields{
				"line": strings.TrimSpace(lineInfo.LineOriginal),
			}).Debugf("%s", err.Error())
			continue
		}

		// get ingredient, continue if its not found
		err = lineInfo.getIngredient()
		if err != nil {
			log.WithFields(logrus.Fields{
				"line": strings.TrimSpace(lineInfo.LineOriginal),
			}).Debugf("%s", err.Error())
			continue
		}

		// get measure
		err = lineInfo.getMeasure()
		if err != nil {
			log.WithFields(logrus.Fields{
				"line": strings.TrimSpace(lineInfo.LineOriginal),
			}).Debugf("%s", err.Error())
		}

		// get comment
		if len(lineInfo.MeasureInString) > 0 && len(lineInfo.IngredientsInString) > 0 {
			lineInfo.Ingredient.Comment = GetOtherInBetweenPositions(lineInfo.Line, lineInfo.MeasureInString[0], lineInfo.IngredientsInString[0])
		}

		log.WithFields(logrus.Fields{
			"line": strings.TrimSpace(lineInfo.LineOriginal),
		}).Debugf("%s (%s): %+v", lineInfo.Ingredient.Name, lineInfo.Ingredient.Comment, lineInfo.Ingredient.MeasureOriginal)

		r.Lines = append(r.Lines, lineInfo)
	}
	wg.Done()
	wg.Wait()

	return
}

func (r *Recipe) PrintIngredientList() string {
	s := ""
	for _, li := range r.Lines {
		s += fmt.Sprintf("%s %s %s\n", AmountToString(li.Ingredient.MeasureOriginal.Amount), li.Ingredient.MeasureOriginal.Name, li.Ingredient.Name)
	}
	return s
}

// IngredientList will return a string containing the ingredient list
func (r *Recipe) IngredientList() (ingredientList IngredientList) {
	ingredientList = IngredientList{make([]Ingredient, len(r.Lines))}
	for i, li := range r.Lines {
		ingredientList.Ingredients[i] = li.Ingredient
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

func AmountToString(amount float64) string {
	r, _ := ParseDecimal(fmt.Sprintf("%2.10f", amount))
	rationalFraction := float64(r.n) / float64(r.d)
	if rationalFraction > 0 {
		bestFractionDiff := 1e9
		bestFraction := 0.0
		var fractions = map[float64]string{
			1.0 / 2: "1/2",
			1.0 / 3: "1/3",
			2.0 / 3: "2/3",
			1.0 / 8: "1/8",
			3.0 / 8: "3/8",
			5.0 / 8: "5/8",
			7.0 / 8: "7/8",
			1.0 / 4: "1/4",
			3.0 / 4: "3/4",
		}
		for f := range fractions {
			currentDiff := math.Abs(f - rationalFraction)
			if currentDiff < bestFractionDiff {
				bestFraction = f
				bestFractionDiff = currentDiff
			}
		}
		if r.i > 0 {
			return strconv.FormatInt(r.i, 10) + " " + fractions[bestFraction]
		} else {
			return fractions[bestFraction]
		}
	}
	return strconv.FormatInt(r.i, 10)
}

// A rational number r is expressed as the fraction p/q of two integers:
// r = p/q = (d*i+n)/d.
type Rational struct {
	i int64 // integer
	n int64 // fraction numerator
	d int64 // fraction denominator
}

func gcd(x, y int64) int64 {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

func ParseDecimal(s string) (r Rational, err error) {
	sign := int64(1)
	if strings.HasPrefix(s, "-") {
		sign = -1
	}
	p := strings.IndexByte(s, '.')
	if p < 0 {
		p = len(s)
	}
	if i := s[:p]; len(i) > 0 {
		if i != "+" && i != "-" {
			r.i, err = strconv.ParseInt(i, 10, 64)
			if err != nil {
				return Rational{}, err
			}
		}
	}
	if p >= len(s) {
		p = len(s) - 1
	}
	if f := s[p+1:]; len(f) > 0 {
		n, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			return Rational{}, err
		}
		d := math.Pow10(len(f))
		if math.Log2(d) > 63 {
			err = fmt.Errorf(
				"ParseDecimal: parsing %q: value out of range", f,
			)
			return Rational{}, err
		}
		r.n = int64(n)
		r.d = int64(d)
		if g := gcd(r.n, r.d); g != 0 {
			r.n /= g
			r.d /= g
		}
		r.n *= sign
	}
	return r, nil
}
