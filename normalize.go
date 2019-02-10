package parserecipe

import (
	"errors"
	"fmt"
	"math"
)

var gramConversions = map[string]float64{
	"ounce": 28.3495,
	"gram":  1,
	"pound": 453.592,
}

var conversionToCup = map[string]float64{
	"tbl":        0.0625,
	"tsp":        0.020833,
	"cup":        1.0,
	"pint":       2.0,
	"quart":      4.0,
	"gallon":     16.0,
	"milliliter": 0.00423,
	"can":        1.75,
}
var ingredientToCups = map[string]float64{
	"eggs":    0.25,
	"garlic":  0.0280833,
	"chicken": 3,
	"celery":  0.5,
	"onion":   1,
	"carrot":  1,
}

// normalizeIngredient will try to normalize the ingredient to 1 cup
func normalizeIngredient(ingredient, measure string, amount float64) (cups float64, err error) {
	// convert measure to standard measure
	newMeasure, ok := corpusMeasuresMap[measure]
	if !ok && measure != "whole" {
		err = fmt.Errorf("could not find '%s'", measure)
		return
	}
	measure = newMeasure
	if _, ok := ingredientToCups[ingredient]; ok && measure == "" {
		// special ingredients
		cups = amount * ingredientToCups[ingredient]
	} else if _, ok := conversionToCup[measure]; ok {
		// check if it has a standard volume measurement
		cups = float64(amount) * conversionToCup[measure]
	} else if _, ok := gramConversions[measure]; ok {
		// check if it has a standard weight measurement
		var density float64
		density, ok = densities[ingredient]
		if !ok {
			density = 200 // grams / cup
		}
		cups = amount * gramConversions[measure] / density
	} else {
		if _, ok := fruitMap[ingredient]; ok {
			cups = 1 * amount
		} else if _, ok := vegetableMap[ingredient]; ok {
			cups = 1 * amount
		} else if _, ok := herbMap[ingredient]; ok {
			cups = 0.0208333 * amount
		} else {
			err = errors.New("could not convert weight or volume")
		}
	}
	return
}

func determineMeasurementsFromCups(cups float64) (amount float64, measure string, amountString string, err error) {
	if cups > 0.125 {
		amount = cups
		measure = "cup"
	} else if cups > 0.020833*3 {
		amount = cups * 16
		measure = "tablespoon"
	} else {
		amount = cups * 48
		measure = "teaspoon"
	}
	amountString = AmountToString(amount)
	if math.IsInf(amount, 0) {
		amount = 0
	}
	if math.IsInf(cups, 0) {
		cups = 0
	}
	return
}
