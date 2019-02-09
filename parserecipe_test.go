package parserecipe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	files := []string{
		"testing/sites/lasagna.html",
		"testing/sites/chocolatecake.html",
		"testing/sites/macandcheese.html",
		"testing/sites/granola-recipe-1939521",
	}
	for _, f := range files {
		log.Infof("working on %s", f)
		_, err := Parse(f)
		assert.Nil(t, err)
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
