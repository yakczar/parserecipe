package parserecipe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	files := []string{
		"testing/sites/lasagna.html",
		"testing/sites/chocolatecake.html",
		"testing/sites/macandcheese.html",
		"testing/sites/granola-recipe-1939521",
		"testing/sites/1017060-doughnuts",
		"testing/sites/poutine.html",
		"testing/sites/waffles.html",
		"testing/sites/refriedbeans.html",
		"testing/sites/pecans.html",
	}
	for _, f := range files {
		log.Infof("working on %s", f)
		r, err := NewFromFile(f)
		assert.Nil(t, err)
		err = r.Parse()
		assert.Nil(t, err)
		ingredientList := r.IngredientList()
		if _, err := os.Stat(f + ".ingredients"); os.IsNotExist(err) {
			b, _ := json.MarshalIndent(ingredientList, "", " ")
			ioutil.WriteFile(f+".ingredients", b, 0644)
		} else {
			b, _ := ioutil.ReadFile(f + ".ingredients")
			var previousIngredientList IngredientList
			assert.Nil(t, json.Unmarshal(b, &previousIngredientList))
			assert.Equal(t, previousIngredientList, ingredientList)
		}
	}

}

func TestGetIngredientsInString(t *testing.T) {
	line := SanitizeLine("1/2 cup chilled oil (vegetable or canola oil)")
	wpi := GetIngredientsInString(line)
	assert.Equal(t, "oil", wpi[0].Word)
	assert.Equal(t, 1, len(wpi))
	fmt.Println(wpi)

	wp := GetNumbersInString(line)
	assert.Equal(t, 1, len(wp))
	assert.Equal(t, "½", wp[0].Word)

	wpm := GetMeasuresInString(line)
	assert.Equal(t, 1, len(wpm))
	assert.Equal(t, "cup", wpm[0].Word)

	fmt.Println(GetOtherInBetweenPositions(line, wpm[0], wpi[0]))
}

func TestTopHat(t *testing.T) {
	vector := []float64{0, 0, 0, 1, 0, 1, 1, 0, 0, 5, 4, 2, 6, 4, 1, 0, 0, 0, 4, 0, 0}
	s, e := GetBestTopHatPositions(vector)
	assert.Equal(t, 9, s)
	assert.Equal(t, 14, e)
}

func TestAmountToString(t *testing.T) {
	assert.Equal(t, "1 2/3", AmountToString(1.66666666))
	assert.Equal(t, "10", AmountToString(10))
	assert.Equal(t, "5 3/8", AmountToString(5.38))
	assert.Equal(t, "1/2", AmountToString(0.5))
}
