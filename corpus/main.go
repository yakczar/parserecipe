package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Create("corpus.go")
	check(err)
	defer f.Close()

	f.WriteString("package parserecipe\n")

	// MAIN INGREDIENT LIST
	// sort the ingredient corpus by the length of each term
	// and then by alphabetizing
	bIngredients, err := ioutil.ReadFile("ingredients.txt")
	corpusIngredients := strings.Split(string(bIngredients), "\n")
	ingredientSizes := make(map[string]int)
	for _, ing := range corpusIngredients {
		if len(ing) == 0 {
			continue
		}
		ingredientSizes[ing] = len(ing)
	}

	pl := make(pairList, len(ingredientSizes))
	i := 0
	for k, v := range ingredientSizes {
		pl[i] = pair{k, v}
		i++
	}
	sort.Slice(pl, func(i, j int) bool {
		if pl[i].Value == pl[j].Value {
			return pl[i].Key < pl[j].Key
		}
		return pl[i].Value > pl[j].Value
	})

	corpusIngredients = make([]string, len(pl))
	for i, p := range pl {
		corpusIngredients[i] = " " + strings.TrimSpace(p.Key) + " "
	}
	f.WriteString(`var corpusIngredients = []string{"` + strings.Join(corpusIngredients, `"`+",\n"+`"`) + `"}` + "\n")
	f.Sync()

	// MAIN MEASURE LIST
	pl = make(pairList, len(corpusMeasuresMap))
	i = 0
	for k := range corpusMeasuresMap {
		pl[i] = pair{k, len(k)}
		i++
	}
	sort.Slice(pl, func(i, j int) bool {
		if pl[i].Value == pl[j].Value {
			return pl[i].Key < pl[j].Key
		}
		return pl[i].Value > pl[j].Value
	})
	corpusMeasures := make([]string, len(pl))
	for i, p := range pl {
		corpusMeasures[i] = " " + p.Key + " "
	}
	f.WriteString(`var corpusMeasures = []string{"` + strings.Join(corpusMeasures, `"`+",\n"+`"`) + `"}` + "\n")
	f.Sync()

	// MAKE NUMBERS
	b, err := ioutil.ReadFile("numbers.txt")
	corpusNumbers := strings.Split(string(b), "\n")
	for v := range corpusFractionNumberMap {
		corpusNumbers = append(corpusNumbers, v)
	}
	for i, c := range corpusNumbers {
		corpusNumbers[i] = " " + strings.TrimSpace(c) + " "
	}
	f.WriteString(`var corpusNumbers = []string{"` + strings.Join(corpusNumbers, `"`+",\n"+`"`) + `"}` + "\n")
	f.Sync()

	// MAKE DIRECTIONS CORPUS
	b, err = ioutil.ReadFile("directions_pos.txt")
	corpusDirections := strings.Fields(string(b))
	corpusDirectionsMap := make(map[string]struct{})
	for _, c := range corpusDirections {
		corpusDirectionsMap[strings.ToLower(c)] = struct{}{}
	}
	corpusDirections = make([]string, len(corpusDirectionsMap))
	i = 0
	for c := range corpusDirectionsMap {
		corpusDirections[i] = c
		i++
	}
	f.WriteString(`var corpusDirections = []string{"` + strings.Join(corpusDirections, `"`+",\n"+`"`) + `"}` + "\n")
	f.Sync()

	f.WriteString(`type fractionNumber struct {
		fractionString string
		value          float64
	}
	`)
	f.WriteString(`var corpusFractionNumberMap = map[string]fractionNumber{` + "\n")
	for k := range corpusFractionNumberMap {
		f.WriteString(fmt.Sprintf(`"%s": fractionNumber{"%s",%10.10f},`, k, corpusFractionNumberMap[k].fractionString, corpusFractionNumberMap[k].value) + "\n")
	}
	f.WriteString("}\n\n")

	f.WriteString(`var corpusMeasuresMap = map[string]string{` + "\n")
	for k := range corpusMeasuresMap {
		f.WriteString(fmt.Sprintf(`"%s": "%s",`, k, corpusMeasuresMap[k]) + "\n")
	}
	f.WriteString("}\n\n")

}

type pair struct {
	Key   string
	Value int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type fractionNumber struct {
	fractionString string
	value          float64
}

var corpusFractionNumberMap = map[string]fractionNumber{
	"½": fractionNumber{"1/2", 1.0 / 2},
	"¼": fractionNumber{"1/4", 1.0 / 4},
	"¾": fractionNumber{"3/4", 3.0 / 4},
	"⅛": fractionNumber{"1/8", 1.0 / 8},
	"⅜": fractionNumber{"3/8", 3.0 / 8},
	"⅝": fractionNumber{"5/8", 5.0 / 8},
	"⅞": fractionNumber{"7/8", 7.0 / 8},
	"⅔": fractionNumber{"2/3", 2.0 / 3},
	"⅓": fractionNumber{"1/3", 1.0 / 3},
}

var corpusMeasuresMap = map[string]string{
	"tablespoon":  "tbl",
	"tablespoons": "tbl",
	"tbl":         "tbl",
	"tbsp":        "tbl",
	"tbsps":       "tbl",
	"teaspoons":   "tsp",
	"teaspoon":    "tsp",
	"tsp":         "tsp",
	"tsps":        "tsp",
	"cups":        "cup",
	"cup":         "cup",
	"c":           "cup",
	"ounces":      "ounce",
	"ounce":       "ounce",
	"oz":          "ounce",
	"grams":       "gram",
	"g":           "gram",
	"gram":        "gram",
	"milliliter":  "milliliter",
	"ml":          "milliliter",
	"pint":        "pint",
	"pints":       "pint",
	"quart":       "quart",
	"quarts":      "quart",
	"pound":       "pound",
	"pounds":      "pound",
	"cans":        "can",
	"canned":      "can",
	"can":         "can",
}
