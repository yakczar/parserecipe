package parserecipe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.Nil(t, Parse("testing/sites/chocolatecake.html"))
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
}

func TestTopHat(t *testing.T) {
	vector := []int{0, 0, 0, 1, 0, 1, 1, 0, 0, 5, 4, 2, 6, 4, 1, 0, 0, 0, 4, 0, 0}
	s, e := GetBestTopHatPositions(vector)
	fmt.Println(vector[s:e])

}
