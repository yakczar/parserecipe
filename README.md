# parserecipe

<img src="https://img.shields.io/badge/coverage-93%25-brightgreen.svg?style=flat-square" alt="Code coverage">&nbsp;<a href="https://travis-ci.org/schollz/parserecipe"><img src="https://img.shields.io/travis/schollz/parserecipe.svg?style=flat-square" alt="Build Status"></a>&nbsp;<a href="https://godoc.org/github.com/schollz/parserecipe"><img src="http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square" alt="Go Doc"></a> 

This is a Golang library for extracting culinary recipes from **any site on the internet**. This library compartmentalizes and improve aspects of recipe extraction that I did previously with [schollz/meanrecipe](https://github.com/schollz/meanrecipe) and [schollz/extract_recipe](https://github.com/schollz/extract_recipe). 

_Note:_ This is still a WIP.

## Install

```
go get github.com/schollz/parserecipe
```

## Usage

### Recipe extraction 

Firstly, you can parse the ingredients and directions from any website. For instance here's one I found on the [JoyFoodSunshine](https://joyfoodsunshine.com/the-most-amazing-chocolate-chip-cookies/):

```go
r, _ := NewFromURL("https://joyfoodsunshine.com/the-most-amazing-chocolate-chip-cookies/")
r.Parse()
fmt.Println(r.PrintIngredientList())
// 1 cup butter
// 1 cup sugar
// 1 cup brown sugar
// 2 tsp vanilla
// 2 whole eggs
// 3 cups flour
// 1 tsp baking soda
// 1/2 tsp baking powder
// 1 tsp salt
// 2 cups chocolate

fmt.Println(r.PrintDirections())
// 1) Preheat oven to 375 degrees F. Line a baking pan with parchment paper and set aside.
// 2) In a separate bowl mix flour, baking soda, salt, baking powder. Set aside.
// 3) Cream together butter and sugars until combined.
// 4) Beat in eggs and vanilla until fluffy.
// 5) Mix in the dry ingredients until combined.
// 6) Add 12 oz package of chocolate chips and mix well.
// 7) Roll 3 TBS of dough at a time into balls and place them evenly spaced on your prepared cookie sheets. (alternately, use a small cookie scoop to make your cookies)!
// 8) Bake in preheated oven for approximately 8-10 minutes. Take them out when they are just *BARELY* starting to turn brown.
// 9) Let them sit on the baking pan for 2 minutes before removing to cooling rack.
```

### Recipe averaging

If you extract multiple recipes, you can even average them together. Printing the ingredient list will show percentages, indicating the percentage of recipes that use those ingredients. Here's an example which parses three random pancake recipes and averages them together.

```go
recipes := []string{
    "https://www.graceandgoodeats.com/best-ever-pancake-recipe/",
    "https://cafedelites.com/best-fluffy-pancakes/",
    "https://www.allrecipes.com/recipe/21014/good-old-fashioned-pancakes/",
}
r := make([]*Recipe, len(recipes))
for i := 0; i < len(recipes); i++ {
    r[i], _ = NewFromURL(recipes[i])
    r[i].Parse()
    r[i].Analyze()
}
averageRecipe, _ := AverageRecipes(r)
fmt.Println(averageRecipe.PrintIngredientList())
// 1/8 cup baking powder (100.0%)
// 3/8 tsp baking soda (33.3%)
// 3/8 cup butter (100.0%)
// 2 whole eggs (100.0%)
// 3 1/6 cup flour (100.0%)
// 2 1/2 cup milk (100.0%)
// 1 1/8 tsp salt (100.0%)
// 1/4 cup sugar (100.0%)
// 1 tbl vanilla (33.3%)
```

## Develop

If you modify the `corpus/` information then you will need to run 

```
$ go generate
```

before using the library again.

## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## License

MIT
