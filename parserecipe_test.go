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
	}
	for _, f := range files {
		log.Infof("working on %s", f)
		_, err := Parse(f)
		assert.Nil(t, err)
	}

}

func TestGetIngredientsInString(t *testing.T) {
	line := "1/2 cup oil (vegetable or canola oil)"
	wp := GetIngredientsInString(line)
	assert.Equal(t, "oil", wp[0].Word)
	assert.Equal(t, 1, len(wp))

	wp = GetNumbersInString(line)
	assert.Equal(t, 1, len(wp))
	assert.Equal(t, "1/2", wp[0].Word)

	wp = GetMeasuresInString(line)
	assert.Equal(t, 1, len(wp))
	assert.Equal(t, "cup", wp[0].Word)

	line = "* 3/4 pound mozzarella cheese, sliced"
	fmt.Println(SanitizeLine(line))
	wp = GetIngredientsInString(SanitizeLine(line))
	fmt.Println(wp)
}

func TestTopHat(t *testing.T) {
	vector := []int{0, 0, 0, 1, 0, 1, 1, 0, 0, 5, 4, 2, 6, 4, 1, 0, 0, 0, 4, 0, 0}
	s, e := GetBestTopHatPositions(vector)
	fmt.Println(vector[s:e])

}
