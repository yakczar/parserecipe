package parserecipe

//go:generate go run corpus/main.go
//go:generate gofmt -s -w corpus.go

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bradfitz/slice"
	"github.com/jaytaylor/html2text"
	colorable "github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log = logrus.New()

// initialize the logger
func init() {
	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
	log.SetLevel(logrus.DebugLevel)
}

// WordPosition shows a word and its position
// Note: the position is memory-dependent as it will
// be the position after the last deleted word
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

// GetOtherInBetweenPositions returns the word positions comment string in the ingredients
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

// GetIngredientsInString returns the word positions of the ingredients
func GetIngredientsInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusIngredients)
}

// GetNumbersInString returns the word positions of the numbers in the ingredient string
func GetNumbersInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusNumbers)
}

// GetMeasuresInString returns the word positions of the measures in a ingredient string
func GetMeasuresInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusMeasures)
}

// SanitizeLine removes parentheses, trims the line, converts to lower case,
// replaces fractions with unicode and then does special conversion for ingredients (like eggs).
func SanitizeLine(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, "⁄", "/", -1)
	s = strings.Replace(s, " / ", "/", -1)

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

// LineInfo has all the information for the parsing of a given line
type LineInfo struct {
	LineOriginal        string
	Line                string         `json:",omitempty"`
	IngredientsInString []WordPosition `json:",omitempty"`
	AmountInString      []WordPosition `json:",omitempty"`
	MeasureInString     []WordPosition `json:",omitempty"`
	Ingredient          Ingredient     `json:",omitempty"`
}

// Recipe contains the info for the file and the lines
type Recipe struct {
	FileName    string
	FileContent string
	Lines       []LineInfo
	Directions  []string
	Ingredients []Ingredient
	Ratios      map[string]map[string]float64
}

// NewFromFile generates a new parser from a file
func NewFromFile(fname string) (r *Recipe, err error) {
	r = &Recipe{FileName: fname}
	_, err = os.Stat(fname)
	return
}

// NewFromURL generates a new parser from a url
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

// Ingredient is the basic struct for ingredients
type Ingredient struct {
	Name      string  `json:",omitempty"`
	Comment   string  `json:",omitempty"`
	Measure   Measure `json:",omitempty"`
	Frequency float64 `json:",omitempty`
}

// Measure includes the amount, name and the cups for conversions
type Measure struct {
	Amount float64
	Name   string
	Cups   float64
}

// IngredientList is a list of ingredients
type IngredientList struct {
	Ingredients []Ingredient
}

func (r *Recipe) parseDirections(lis []LineInfo) (rerr error) {
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

	start, end := getBestTopHatPositions(scores)
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

// Parse is the main parser for a given recipe.
// It looks for the following
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
		// does it contain an ingredient?
		if len(lineInfos[i].IngredientsInString) > 0 {
			score++
		}
		// does it contain an amount?
		if len(lineInfos[i].AmountInString) > 0 {
			score++
		}
		// does it contain a measure (cups, tsps)?
		if len(lineInfos[i].MeasureInString) > 0 {
			score++
		}
		// does the ingredient come after the measure?
		if len(lineInfos[i].IngredientsInString) > 0 && len(lineInfos[i].MeasureInString) > 0 && lineInfos[i].IngredientsInString[0].Position > lineInfos[i].MeasureInString[0].Position {
			score++
		}
		// does the ingredient come after the amount?
		if len(lineInfos[i].IngredientsInString) > 0 && len(lineInfos[i].AmountInString) > 0 && lineInfos[i].IngredientsInString[0].Position > lineInfos[i].AmountInString[0].Position {
			score++
		}
		// does the measure come after the amount?
		if len(lineInfos[i].MeasureInString) > 0 && len(lineInfos[i].AmountInString) > 0 && lineInfos[i].MeasureInString[0].Position > lineInfos[i].AmountInString[0].Position {
			score++
		}
		// is the line really long? (ingredient lines are short)
		if score > 0 && len(lineInfos[i].LineOriginal) > 100 {
			score--
		}
		// does it start with a list indicator (* or -)?
		fields := strings.Fields(line)
		if len(fields) > 0 && (fields[0] == "*" || fields[0] == "-") {
			score++
		}
		// if only one thing is right, its wrong
		if score == 1 {
			score = 0.0
		}
		// log.Debugf("'%s' (%d)", line, score)
		scores[i] = score
	}
	scores = scores[:i+1]
	lineInfos = lineInfos[:i+1]

	// debugging purposes
	// lines = make([]string, len(lineInfos))
	// for i, li := range lineInfos {
	// 	lines[i] = li.Line
	// }
	// ioutil.WriteFile("out", []byte(strings.Join(lines, "\n")), 0644)

	// get the most likely location
	start, end := getBestTopHatPositions(scores)

	var wg sync.WaitGroup
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := r.parseDirections(lineInfos[end:])
		if err != nil {
			log.Warn(err.Error())
		}
	}(&wg)

	r.Lines = []LineInfo{}
	for _, lineInfo := range lineInfos[start-3 : end+3] {
		if len(strings.TrimSpace(lineInfo.Line)) < 3 {
			continue
		}

		lineInfo.Ingredient.Measure = Measure{}

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

		// normalize into cups
		lineInfo.Ingredient.Measure.Cups, err = normalizeIngredient(
			lineInfo.Ingredient.Name,
			lineInfo.Ingredient.Measure.Name,
			lineInfo.Ingredient.Measure.Amount,
		)
		if err != nil {
			log.WithFields(logrus.Fields{
				"line": strings.TrimSpace(lineInfo.LineOriginal),
			}).Debugf("can't convert to cups: %s", err.Error())
		}

		log.WithFields(logrus.Fields{
			"sanitize": strings.TrimSpace(lineInfo.Line),
			"original": strings.TrimSpace(lineInfo.LineOriginal),
		}).Debugf("%s (%s): %+v", lineInfo.Ingredient.Name, lineInfo.Ingredient.Comment, lineInfo.Ingredient.Measure)

		r.Lines = append(r.Lines, lineInfo)
	}
	rerr = r.ConvertIngredients()
	if rerr != nil {
		return
	}

	// consolidate ingredients
	ingredients := make(map[string]Ingredient)
	ingredientList := []string{}
	for _, line := range r.Lines {
		if _, ok := ingredients[line.Ingredient.Name]; ok {
			if ingredients[line.Ingredient.Name].Measure.Name == line.Ingredient.Measure.Name {
				ingredients[line.Ingredient.Name] = Ingredient{
					Name:    line.Ingredient.Name,
					Comment: ingredients[line.Ingredient.Name].Comment,
					Measure: Measure{
						Name:   ingredients[line.Ingredient.Name].Measure.Name,
						Amount: ingredients[line.Ingredient.Name].Measure.Amount + line.Ingredient.Measure.Amount,
						Cups:   ingredients[line.Ingredient.Name].Measure.Cups + line.Ingredient.Measure.Cups,
					},
				}
			} else {
				ingredients[line.Ingredient.Name] = Ingredient{
					Name:    line.Ingredient.Name,
					Comment: ingredients[line.Ingredient.Name].Comment,
					Measure: Measure{
						Name:   ingredients[line.Ingredient.Name].Measure.Name,
						Amount: ingredients[line.Ingredient.Name].Measure.Amount,
						Cups:   ingredients[line.Ingredient.Name].Measure.Cups + line.Ingredient.Measure.Cups,
					},
				}
				log.Debugf("different measure!")
			}
		} else {
			ingredientList = append(ingredientList, line.Ingredient.Name)
			ingredients[line.Ingredient.Name] = Ingredient{
				Name:    line.Ingredient.Name,
				Comment: line.Ingredient.Comment,
				Measure: Measure{
					Name:   line.Ingredient.Measure.Name,
					Amount: line.Ingredient.Measure.Amount,
					Cups:   line.Ingredient.Measure.Cups + line.Ingredient.Measure.Cups,
				},
			}
		}
	}
	r.Ingredients = make([]Ingredient, len(ingredients))
	for i, ing := range ingredientList {
		r.Ingredients[i] = ingredients[ing]
	}

	wg.Done()
	wg.Wait()

	return
}

func (r *Recipe) ConvertIngredients() (err error) {

	return
}

func (r *Recipe) PrintIngredientList() string {
	s := ""
	for _, li := range r.Lines {
		s += fmt.Sprintf("%s %s %s\n", AmountToString(li.Ingredient.Measure.Amount), li.Ingredient.Measure.Name, li.Ingredient.Name)
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
		lineInfo.Ingredient.Measure.Amount = totalAmount
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
		lineInfo.Ingredient.Measure.Name = "whole"
		return
	}
	lineInfo.Ingredient.Measure.Name = lineInfo.MeasureInString[0].Word
	return
}

func getBestTopHatPositions(vectorFloat []float64) (start, end int) {
	bestTopHatResidual := 1e9
	for i, v := range vectorFloat {
		if v < 2 {
			continue
		}
		for j, w := range vectorFloat {
			if j <= i || w < 1 {
				continue
			}
			hat := generateHat(len(vectorFloat), i, j, averageFloats(vectorFloat[i:j]))
			res := calculateResidual(vectorFloat, hat) / float64(len(vectorFloat))
			if res < bestTopHatResidual {
				bestTopHatResidual = res
				start = i
				end = j
			}
		}
	}
	return
}

func calculateResidual(fs1, fs2 []float64) float64 {
	res := 0.0
	if len(fs1) != len(fs2) {
		return -1
	}
	for i := range fs1 {
		res += math.Pow(fs1[i]-fs2[i], 2)
	}
	return res
}

func averageFloats(fs []float64) float64 {
	f := 0.0
	for _, v := range fs {
		f += v
	}
	return f / float64(len(fs))
}

func generateHat(length, start, stop int, value float64) []float64 {
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
	r, _ := parseDecimal(fmt.Sprintf("%2.10f", amount))
	rationalFraction := float64(r.n) / float64(r.d)
	if rationalFraction > 0 {
		bestFractionDiff := 1e9
		bestFraction := 0.0
		var fractions = map[float64]string{
			0:       "",
			1:       "",
			1.0 / 2: "1/2",
			1.0 / 3: "1/3",
			2.0 / 3: "2/3",
			1.0 / 6: "1/6",
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
		if fractions[bestFraction] == "" {
			return strconv.FormatInt(int64(math.Round(amount)), 10)
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
type rational struct {
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

func parseDecimal(s string) (r rational, err error) {
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
				return rational{}, err
			}
		}
	}
	if p >= len(s) {
		p = len(s) - 1
	}
	if f := s[p+1:]; len(f) > 0 {
		n, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			return rational{}, err
		}
		d := math.Pow10(len(f))
		if math.Log2(d) > 63 {
			err = fmt.Errorf(
				"ParseDecimal: parsing %q: value out of range", f,
			)
			return rational{}, err
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

// Analyze will determine the ratios between all the normalized ingredients
func (r *Recipe) Analyze() (err error) {
	r.Ratios = make(map[string]map[string]float64)
	for _, ing1 := range r.Ingredients {
		r.Ratios[ing1.Name] = make(map[string]float64)
		for _, ing2 := range r.Ingredients {
			if ing1.Name >= ing2.Name {
				continue
			}
			r.Ratios[ing1.Name][ing2.Name] = ing1.Measure.Cups / ing2.Measure.Cups
		}
	}
	return
}

func compareRatios(r1, r2 map[string]map[string]float64, debug ...bool) (sumsq float64) {
	for ing1 := range r1 {
		if _, ok := r2[ing1]; !ok {
			continue
		}
		for ing2 := range r1[ing1] {
			if _, ok := r2[ing1][ing2]; !ok {
				continue
			}
			if len(debug) > 0 && debug[0] {
				log.Debugf("%s/%s %2.3f %2.3f", ing1, ing2, r1[ing1][ing2], r2[ing1][ing2])
			}
			sumsq += math.Pow(r1[ing1][ing2]-r2[ing1][ing2], 2)
		}
	}
	return sumsq
}

// AverageRecipes returns the distance between the recipes
func AverageRecipes(rs []*Recipe) (averagedRecipe *Recipe, err error) {
	if len(rs) < 2 {
		averagedRecipe = rs[0]
		return
	}
	totalFull := make([]float64, len(rs))
	// totalPartial := make([]float64, len(rs))
	rIngredients := make([]map[string]Ingredient, len(rs))
	ingredientFrequencies := make(map[string]float64)
	allIngredients := []string{}

	for i := 0; i < len(rs); i++ {
		rIngredients[i] = make(map[string]Ingredient)
		for _, ing := range rs[i].Ingredients {
			if _, ok := ingredientFrequencies[ing.Name]; !ok {
				ingredientFrequencies[ing.Name] = 0
				allIngredients = append(allIngredients, ing.Name)
			}
			ingredientFrequencies[ing.Name] += 1 / float64(len(rs))
			totalFull[i] += ing.Measure.Cups
			rIngredients[i][ing.Name] = ing
		}
	}
	sort.Sort(sort.Float64Slice(totalFull))
	log.Debug(totalFull)
	medianTotal := totalFull[len(totalFull)/2]
	log.Debugf("ingredientFrequencies: %+v", ingredientFrequencies)

	ingredientRatios := make(map[string]map[string][]float64)
	for i := 0; i < len(rs); i++ {
		for ing := range rs[i].Ratios {
			if _, ok := ingredientRatios[ing]; !ok {
				ingredientRatios[ing] = make(map[string][]float64)
			}
			for ing2 := range rs[i].Ratios[ing] {
				if _, ok := ingredientRatios[ing][ing2]; !ok {
					ingredientRatios[ing][ing2] = []float64{}
				}
				ingredientRatios[ing][ing2] = append(ingredientRatios[ing][ing2], rs[i].Ratios[ing][ing2])
			}
		}
	}
	log.Debugf("ingredientRatios: %+v", ingredientRatios)

	averageRatios := make(map[string]map[string]float64)
	for ing1 := range ingredientRatios {
		averageRatios[ing1] = make(map[string]float64)
		for ing2 := range ingredientRatios[ing1] {
			averageRatios[ing1][ing2] = averageFloats(ingredientRatios[ing1][ing2])
		}
	}

	// its not gauranteed that the average ratios are normalized, so many recipes should be created
	// based on the normalizations and the closest one should be taken
	s := rand.NewSource(time.Now().Unix())
	ran := rand.New(s) // initialize local pseudorandom generator
	averagedRecipe = new(Recipe)
	bestRecipeSumSQ := 1e9
	for iterations := 0; iterations < 100; iterations++ {
		randPerm := ran.Perm(len(allIngredients))
		aIngredients := make(map[string]Ingredient)
		aIngredients[allIngredients[randPerm[0]]] = Ingredient{
			Name:      allIngredients[randPerm[0]],
			Measure:   Measure{Amount: 1, Cups: 1},
			Frequency: ingredientFrequencies[allIngredients[randPerm[0]]],
		}
		// log.Debugf("%s determined to be 1", allIngredients[randPerm[0]])
		for {
			if len(aIngredients) == len(randPerm) {
				break
			}
			for i := 0; i < len(randPerm); i++ {
				ing := allIngredients[randPerm[i]]
				if _, ok := aIngredients[ing]; ok {
					continue
				}
				for ingDone := range aIngredients {
					if _, ok := aIngredients[ing]; ok {
						break
					}
					var ing1, ing2 string
					if ingDone == ing {
						continue
					}
					ing1 = ing
					ing2 = ingDone
					if ing > ingDone {
						ing1 = ingDone
						ing2 = ing
					}
					if _, ok := averageRatios[ing1]; ok {
						if _, ok := averageRatios[ing1][ing2]; ok {
							amount := averageRatios[ing1][ing2]
							if ingDone > ing {
								amount = amount * aIngredients[ingDone].Measure.Cups
							} else {
								amount = 1 / amount * aIngredients[ingDone].Measure.Cups
							}
							aIngredients[ing] = Ingredient{
								Name:      ing,
								Frequency: ingredientFrequencies[ing],
								Measure:   Measure{Cups: amount, Amount: amount, Name: "cups"},
							}
							// log.Debugf("%s determined from %s to be %2.5f", ing, ingDone, amount)
						}
					}
				}
			}
		}

		rNew := new(Recipe)
		rNew.Ingredients = make([]Ingredient, 0, len(aIngredients))
		for ing := range aIngredients {
			rNew.Ingredients = append(rNew.Ingredients, aIngredients[ing])
		}
		rNew.Analyze()
		ratioComparison := compareRatios(averageRatios, rNew.Ratios)
		if ratioComparison < bestRecipeSumSQ {
			log.Debugf("ratio comparison: %2.3f", ratioComparison)
			bestRecipeSumSQ = ratioComparison
			averagedRecipe.Ingredients = make([]Ingredient, 0, len(aIngredients))
			for ing := range aIngredients {
				averagedRecipe.Ingredients = append(averagedRecipe.Ingredients, aIngredients[ing])
			}
		}
	}
	averagedRecipe.Analyze()
	log.Debugf("comparison: %2.3f", compareRatios(averageRatios, averagedRecipe.Ratios, true))

	newRecipeTotal := 0.0
	for i := range averagedRecipe.Ingredients {
		newRecipeTotal += averagedRecipe.Ingredients[i].Measure.Cups
	}
	slice.Sort(averagedRecipe.Ingredients[:], func(i, j int) bool {
		return averagedRecipe.Ingredients[i].Name < averagedRecipe.Ingredients[j].Name
	})
	scalingFactor := medianTotal / newRecipeTotal
	for i := range averagedRecipe.Ingredients {
		averagedRecipe.Ingredients[i].Measure.Cups = averagedRecipe.Ingredients[i].Measure.Cups * scalingFactor
		averagedRecipe.Ingredients[i].Measure.Amount, averagedRecipe.Ingredients[i].Measure.Name = cupsToOther(
			averagedRecipe.Ingredients[i].Measure.Cups,
			averagedRecipe.Ingredients[i].Name,
		)
		log.Debugf("%s (%2.2f) %s %s (%2.1f%%)",
			AmountToString(averagedRecipe.Ingredients[i].Measure.Amount),
			averagedRecipe.Ingredients[i].Measure.Cups,
			averagedRecipe.Ingredients[i].Measure.Name,
			averagedRecipe.Ingredients[i].Name,
			averagedRecipe.Ingredients[i].Frequency*100,
		)
	}
	log.Debugf("best recipe scaled to %2.2f cups", medianTotal)
	log.Debugf("averagedRecipe: %+v", averagedRecipe)
	// log.Debug(averagedRecipe.PrintIngredientList())

	// ingredientsInBoth := []string{}
	// ingredientsInBothMap := make(map[string]struct{})
	// for ing := range r0Ingredients {
	// 	if _, ok := r1Ingredients[ing]; ok {
	// 		ingredientsInBoth = append(ingredientsInBoth, ing)
	// 		ingredientsInBothMap[ing] = struct{}{}
	// 	} else {
	// 		log.Debugf("r0 has %s %2.3f", ing, r0Ingredients[ing].Measure.Cups/totalFull[0])
	// 	}
	// }
	// for ing := range r1Ingredients {
	// 	if _, ok := r0Ingredients[ing]; !ok {
	// 		log.Debugf("r1 has %s %2.3f", ing, r1Ingredients[ing].Measure.Cups/totalFull[1])
	// 	}
	// }

	// missing := [2]int{0, 0}
	// missing[0] = len(r0.Ratios) - len(ingredientsInBoth)
	// missing[1] = len(r1.Ratios) - len(ingredientsInBoth)
	// log.Debugf("missing: %+v", missing)

	// for _, ing := range ingredientsInBoth {
	// 	totalPartial[0] += r0Ingredients[ing].Measure.Cups
	// 	totalPartial[1] += r1Ingredients[ing].Measure.Cups
	// }

	// for _, ing := range ingredientsInBoth {
	// 	log.Debugf("%s %2.1f%% %2.1f%%", ing,
	// 		r0Ingredients[ing].Measure.Cups/totalPartial[0]*100,
	// 		r1Ingredients[ing].Measure.Cups/totalPartial[1]*100,
	// 	)
	// }

	// for _, ing1 := range ingredientsInBoth {
	// 	for _, ing2 := range ingredientsInBoth {
	// 		if ing1 >= ing2 {
	// 			continue
	// 		}
	// 		log.Debugf("%s/%s %2.2f %2.2f",
	// 			ing1, ing2,
	// 			r0.Ratios[ing1][ing2],
	// 			r1.Ratios[ing1][ing2],
	// 		)
	// 	}
	// }

	return
}
