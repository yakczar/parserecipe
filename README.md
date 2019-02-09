# parserecipe

<img src="https://img.shields.io/badge/coverage-93%25-brightgreen.svg?style=flat-square" alt="Code coverage">&nbsp;<a href="https://travis-ci.org/schollz/parserecipe"><img src="https://img.shields.io/travis/schollz/parserecipe.svg?style=flat-square" alt="Build Status"></a>&nbsp;<a href="https://godoc.org/github.com/schollz/parserecipe"><img src="http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square" alt="Go Doc"></a> 

This is a Golang library for extracting culinary recipes from the internet. This library compartmentalizes and improve aspects of recipe extraction that I did previously with [schollz/meanrecipe](https://github.com/schollz/meanrecipe) and [schollz/extract_recipe](https://github.com/schollz/extract_recipe). Maybe this one will be the last one that I write.

_Note:_ This is still a WIP.

## Install

```
go get github.com/schollz/parserecipe
```

## Usage

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
```

## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## License

MIT
