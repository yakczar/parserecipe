package parserecipe

import (
	"sort"
	"strings"
)

func init() {
	// sort the ingredient corpus by the length of each term
	ingredientSizes := make(map[string]int)
	for _, ing := range corpusIngredients {
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
		corpusIngredients[i] = p.Key
	}

	// sort corpus measures
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
	corpusMeasures = make([]string, len(pl))
	for i, p := range pl {
		corpusMeasures[i] = p.Key
	}

	for v := range corpusFractionNumberMap {
		corpusNumbers = append(corpusNumbers, v)
	}

	corpusDirectionsMap := make(map[string]struct{})
	for _, c := range corpusDirections {
		corpusDirectionsMap[c] = struct{}{}
	}
	corpusDirections = make([]string, len(corpusDirectionsMap))
	i = 0
	for c := range corpusDirectionsMap {
		corpusDirections[i] = c
		i++
	}

	// make sure each is flanked by space
	for i, c := range corpusMeasures {
		corpusMeasures[i] = " " + strings.TrimSpace(c) + " "
	}
	for i, c := range corpusIngredients {
		corpusIngredients[i] = " " + strings.TrimSpace(c) + " "
	}
	for i, c := range corpusNumbers {
		corpusNumbers[i] = " " + strings.TrimSpace(c) + " "
	}
}

type pair struct {
	Key   string
	Value int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var corpusNumbers = strings.Split(`1/2
1/3
1/4
1/5
1/6
1/7
1/8
2/3
2/5
2/7
3/4
3/8
4/5
5/8
7/8
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20`, "\n")

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

var corpusIngredients = strings.Split(`salt
sugar
butter
garlic
water
olive oil
flour
powdered sugar
milk
flour
onion
pepper
onions
brown sugar
eggs
cinnamon
baking powder
lemon juice
matzo meal
tomatoes
vanilla
parsley
baking soda
sour cream
vegetable oil
celery
ginger
lemon
cream cheese
carrots
cheddar cheese
beef
potatoes
oil
honey
nutmeg
cheese
soy sauce
mayonnaise
chicken broth
oregano
cumin
thyme
garlic powder
mushrooms
cilantro
basil
pecans
bacon
heavy cream
chicken breasts
worcestershire sauce
paprika
chocolate
chicken
flax
walnuts
walnut flour
dark chocolate
chili powder
almonds
lime juice
parmesan cheese
pineapple
rice
orange juice
green pepper
raisins
coconut
cayenne pepper
nuts
dijon mustard
cornstarch
mozzarella cheese
buttermilk
vinegar
apples
red pepper
tomato sauce
bread crumbs
steak
oats
spinach
shortening
red pepper flakes
shallots
tomato paste
red bell pepper
lime
shrimp
semolina
zucchini
strawberries
rosemary
canola oil
green onions
bananas
scallions
mustard
chicken stock
chives
whipping cream
maple syrup
orange
corn flour
corn starch
balsamic vinegar
dry white wine
coriander
bay leaf
ketchup
yogurt
red wine vinegar
avocado
sesame oil
cabbage
bay leaves
broccoli
chicken breast
cocoa
carrot
basil leaves
onion powder
cucumber
peanut butter
allspice
dry mustard
cranberries
mint
ham
green bell pepper
blueberries
soda
peas
curry powder
corn
coconut milk
lettuce
white pepper
sesame seeds
pork
turmeric
pasta
dill
yellow onion
white wine
red onion
jalapeno chilies
cream of mushroom
beans
almond flour
almond extract
black beans
garlic salt
peanuts
cider vinegar
white vinegar
margarine
green beans
cream
molasses
pumpkin
coconut oil
rice noodle
rice flour
turkey
yeast
olives
corn syrup
sage
rice vinegar
raspberries
beef broth
ricotta cheese
salsa
tomato
spray
cilantro leaves
parsley leaves
apple cider vinegar
capers
bell pepper
gelatin
green chilies
black olives
feta cheese
swiss cheese
cherry tomatoes
potato
potato starch
oranges
cool whip
cream of tartar
cornmeal
pineapple juice
italian seasoning
cherries
cauliflower
white wine vinegar
whipped cream
applesauce
asparagus
thyme leaves
salmon
cooking oil
cayenne
flour tortillas
dates
leeks
purple onion
green onion
mint leaves
dressing
skim milk
mango
graham cracker crumbs
fish sauce
peanut oil
red wine
cottage cheese
salad oil
heavy whipping cream
tuna
apple
sausage
vanilla ice cream
cooking spray
eggplant
plum tomatoes
tarragon
thru
peaches
goat cheese
kidney beans
tofu
corn tortillas
chickpeas
vegetable broth
celery seed
shallot
clove
chicken soup
spaghetti
lemon peel
black peppercorns
peppermint
banana
hamburger
cardamom
catsup
brandy
salad
horseradish
vodka
sweet potatoes
beer
coffee
butternut squash
white onion
smoked paprika
apple juice
chile
pie shell
pumpkin pie spice
lemons
vegetable stock
egg noodles
broccoli florets
pine nuts
sweet onion
pears
brown rice
parsley flakes
red peppers
quinoa
hot pepper sauce
tomato soup
dry sherry
blue cheese
arugula
dry red wine
corn kernels
hot sauce
green peppers
cumin seed
barbecue sauce
artichoke hearts
water chestnuts
lemon rind
chili sauce
tabasco sauce
beef stock
orange peel
marshmallows
kale
bread flour
vegetable shortening
american cheese
dill weed
fruit
white rice
hazelnuts
crabmeat
pie crust
beets
almond milk
almond meal
oat flour
marjoram
baby spinach
graham crackers
prosciutto
fennel
tomato juice
evaporated milk
parmesan cheese
yellow cornmeal
seasoning salt
garam masala
lamb
evaporated milk
salt
sugar
butter
garlic
water
olive oil
milk
flour
onion
pepper
onions
black pepper
brown sugar
eggs
cinnamon
baking powder
lemon juice
tomatoes
vanilla
parsley
baking soda
sour cream
vegetable oil
celery
ginger
lemon
cream cheese
carrots
cheddar cheese
beef
potatoes
oil
honey
nutmeg
cheese
soy sauce
mayonnaise
chicken broth
oregano
cumin
thyme
garlic powder
salt and pepper
mushrooms
cilantro
basil
pecans
bacon
heavy cream
chicken breasts
worcestershire sauce
paprika
chocolate
chicken
walnuts
chili powder
almonds
lime juice
parmesan cheese
pineapple
rice
orange juice
green pepper
raisins
coconut
cayenne pepper
nuts
dijon mustard
cornstarch
mzarella cheese
buttermilk
vinegar
apples
red pepper
tomato sauce
bread crumbs
oats
spinach
shortening
red pepper flakes
shallots
tomato paste
red bell pepper
lime
shrimp
semolina
zucchini
strawberries
rosemary
canola oil
green onions
bananas
scallions
mustard
chicken stock
chives
whipping cream
bread
maple syrup
orange
corn starch
balsamic vinegar
dry white wine
coriander
bay leaf
ketchup
yogurt
red wine vinegar
avocado
sesame oil
cabbage
bay leaves
broccoli
salt and black pepper
chicken breast
cocoa
carrot
basil leaves
onion powder
cucumber
peanut butter
allspice
dry mustard
cranberries
mint
ham
green bell pepper
blueberries
soda
peas
curry powder
corn
coconut milk
lettuce
white pepper
sesame seeds
pork
turmeric
pasta
dill
yellow onion
white wine
red onion
jalapeno chilies
cream of mushroom soup
beans
almond extract
black beans
garlic salt
peanuts
cider vinegar
white vinegar
margarine
green beans
cream
molasses
confectioners sugar
pumpkin
coconut oil
sauce
turkey
yeast
olives
corn syrup
sage
rice vinegar
raspberries
beef broth
salt and pepper
ricotta cheese
salsa
tomato
breadcrumbs
spray
cilantro leaves
parsley leaves
apple cider vinegar
capers
bell pepper
gelatin
green chilies
black olives
feta cheese
swiss cheese
cherry tomatoes
potato
oranges
cool whip
cream of tartar
cornmeal
pineapple juice
italian seasoning
cherries
cauliflower
white wine vinegar
whipped cream
applesauce
asparagus
thyme leaves
salmon
cooking oil
cayenne
flour tortillas
dates
leeks
purple onion
green onion
mint leaves
dressing
skim milk
oatmeal
mango
graham cracker crumbs
fish sauce
peanut oil
red wine
cottage cheese
salad oil
heavy whipping cream
tuna
apple
sausage
vanilla ice cream
cooking spray
eggplant
plum tomatoes
tarragon
thru
peaches
goat cheese
ice
kidney beans
mozzarella cheese
can cream of chicken soup
chicken thighs
tofu
corn tortillas
chickpeas
vegetable broth
lasagna noodles
celery seed
shallot
teriyaki sauce
clove
chicken soup
spaghetti
lemon peel
black peppercorns
half and half
cashews
raisins
onion
yellow cake mix
banana`, "\n")

var corpusMeasures []string

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

var corpusDirections = strings.Fields(` in a dutch oven cook sausage ground beef onion and garlic over medium heat until well browned stir in crushed tomatoes tomato paste tomato sauce and water season with sugar basil fennel seeds italian seasoning 1 tablespoon salt pepper and 2 tablespoons parsley simmer covered for about 1 ½ hours stirring occasionally in a large bowl combine the oats nuts coconut and brown sugar 
bring a large pot of lightly salted water to a boil cook lasagna noodles in boiling water for 8 to 10 minutes drain noodles and rinse with cold water in a mixing bowl combine ricotta cheese with eggs remaining parsley and ½ teaspoon salt 
preheat oven to 375 degrees f combine sheet pans stirring achieve an even color large bowl transfer remove oven mix
to assemble spread 1 ½ cups of meat sauce in the bottom of a 9x13 inch baking dish arrange 6 noodles lengthwise over meat sauce spread with 1 half of the ricotta cheese mixture top with a third of mozzarella cheese slices spoon 1 ½ cups meat sauce over mozzarella and sprinkle with ¼ cup parmesan cheese repeat layers and top with remaining mozzarella and parmesan cheese cover with foil to prevent sticking either spray foil with cooking spray or make sure the foil does not touch the cheese 
bake in preheated oven for 25 minutes remove foil and bake an additional 25 minutes cool for 15 minutes before serving remove from oven and transfer into a large bowl add raisins and mix until evenly distributed
`)
