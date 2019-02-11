package parserecipe

var herbMap = map[string]struct{}{
	"ajwain":               {},
	"alligator pepper":     {},
	"allspice":             {},
	"amchoor":              {},
	"angelica":             {},
	"anise":                {},
	"aonori":               {},
	"aromatic ginger":      {},
	"asafoetida":           {},
	"basil":                {},
	"bay leaf":             {},
	"black cardamom":       {},
	"black mustard":        {},
	"black peppercorn":     {},
	"boldo":                {},
	"bolivian coriander":   {},
	"borage":               {},
	"brazilian pepper":     {},
	"brown mustard":        {},
	"bunium persicum":      {},
	"camphor":              {},
	"caraway":              {},
	"cardamom":             {},
	"cassia":               {},
	"cayenne pepper":       {},
	"celery powder":        {},
	"celery seed":          {},
	"charoli":              {},
	"chenpi":               {},
	"chervil":              {},
	"chili pepper":         {},
	"chives":               {},
	"cicely":               {},
	"cilantro":             {},
	"cinnamon":             {},
	"clove":                {},
	"coriander leaf":       {},
	"coriander seed":       {},
	"cress":                {},
	"cubeb":                {},
	"culantro":             {},
	"cumin":                {},
	"curry leaf":           {},
	"dill":                 {},
	"dill seed":            {},
	"dried lime":           {},
	"east asian pepper":    {},
	"epazote":              {},
	"fennel":               {},
	"fenugreek":            {},
	"fingerroot":           {},
	"garlic":               {},
	"garlic chives":        {},
	"ginger":               {},
	"golpar":               {},
	"grains of paradise":   {},
	"grains of selim":      {},
	"greater galangal":     {},
	"green peppercorn":     {},
	"hemp":                 {},
	"hoja santa":           {},
	"holy basil":           {},
	"horseradish":          {},
	"houttuynia cordata":   {},
	"hyssop":               {},
	"indian bay leaf":      {},
	"jimbu":                {},
	"juniper berry":        {},
	"kinh gioi":            {},
	"kokum":                {},
	"korarima":             {},
	"lavender":             {},
	"lemon balm":           {},
	"lemon grass":          {},
	"lemon myrtle":         {},
	"lemon verbena":        {},
	"lesser galangal":      {},
	"limnophila aromatica": {},
	"liquorice":            {},
	"litsea cubeba":        {},
	"long pepper":          {},
	"lovage":               {},
	"mace":                 {},
	"mahlab":               {},
	"mango-ginger":         {},
	"marjoram":             {},
	"mastic":               {},
	"mint":                 {},
	"mitsuba":              {},
	"mugwort":              {},
	"nigella":              {},
	"nigella sativa":       {},
	"njangsa":              {},
	"nutmeg":               {},
	"oregano":              {},
	"paprika":              {},
	"parsley":              {},
	"perilla":              {},
	"peruvian pepper":      {},
	"pomegranate seed":     {},
	"poppy seed":           {},
	"radhuni":              {},
	"rose":                 {},
	"rosemary":             {},
	"rue":                  {},
	"saffron":              {},
	"sage":                 {},
	"salt":                 {},
	"sansho":               {},
	"sarsaparilla":         {},
	"sassafras":            {},
	"savory":               {},
	"sesame":               {},
	"shiso":                {},
	"sichuan pepper":       {},
	"sorrel":               {},
	"star anise":           {},
	"sumac":                {},
	"tamarind":             {},
	"tarragon":             {},
	"tasmanian pepper":     {},
	"thai basil":           {},
	"thyme":                {},
	"tonka bean":           {},
	"turmeric":             {},
	"uzazi":                {},
	"vanilla":              {},
	"vietnamese coriander": {},
	"voatsiperifery":       {},
	"wasabi":               {},
	"white mustard":        {},
	"white peppercorn":     {},
	"woodruff":             {},
	"yuzu":                 {},
	"zedoary":              {},
	"zereshk":              {},
	"zest":                 {},
}

var fruitMap = map[string]struct{}{
	"apple":             {},
	"apricot":           {},
	"avocado":           {},
	"banana":            {},
	"bell pepper":       {},
	"bilberry":          {},
	"blackberry":        {},
	"blackcurrant":      {},
	"blood orange":      {},
	"blueberry":         {},
	"boysenberry":       {},
	"breadfruit":        {},
	"canary melon":      {},
	"cantaloupe":        {},
	"cherimoya":         {},
	"cherry":            {},
	"chili pepper":      {},
	"clementine":        {},
	"cloudberry":        {},
	"coconut":           {},
	"cranberry":         {},
	"cucumber":          {},
	"currant":           {},
	"damson":            {},
	"date":              {},
	"dragonfruit":       {},
	"durian":            {},
	"eggplant":          {},
	"elderberry":        {},
	"feijoa":            {},
	"fig":               {},
	"goji berry":        {},
	"gooseberry":        {},
	"grape":             {},
	"grapefruit":        {},
	"guava":             {},
	"honeydew":          {},
	"huckleberry":       {},
	"jackfruit":         {},
	"jambul":            {},
	"jujube":            {},
	"kiwi fruit":        {},
	"kumquat":           {},
	"lemon":             {},
	"lime":              {},
	"loquat":            {},
	"lychee":            {},
	"mandarine":         {},
	"mango":             {},
	"mulberry":          {},
	"nectarine":         {},
	"nut":               {},
	"olive":             {},
	"orange":            {},
	"pamelo":            {},
	"papaya":            {},
	"passionfruit":      {},
	"peach":             {},
	"pear":              {},
	"persimmon":         {},
	"physalis":          {},
	"pineapple":         {},
	"plum":              {},
	"pomegranate":       {},
	"pomelo":            {},
	"purple mangosteen": {},
	"quince":            {},
	"raisin":            {},
	"rambutan":          {},
	"raspberry":         {},
	"redcurrant":        {},
	"rock melon":        {},
	"salal berry":       {},
	"satsuma":           {},
	"star fruit":        {},
	"strawberry":        {},
	"tamarillo":         {},
	"tangerine":         {},
	"ugli fruit":        {},
	"watermelon":        {},
}

var vegetableMap = map[string]struct{}{
	"acorn squash":        {},
	"alfalfa sprout":      {},
	"amaranth":            {},
	"anise":               {},
	"artichoke":           {},
	"arugula":             {},
	"asparagus":           {},
	"aubergine":           {},
	"azuki bean":          {},
	"banana squash":       {},
	"basil":               {},
	"bean sprout":         {},
	"beet":                {},
	"black bean":          {},
	"black-eyed pea":      {},
	"bok choy":            {},
	"borlotti bean":       {},
	"broad beans":         {},
	"broccoflower":        {},
	"broccoli":            {},
	"brussels sprout":     {},
	"butternut squash":    {},
	"cabbage":             {},
	"calabrese":           {},
	"caraway":             {},
	"carrot":              {},
	"cauliflower":         {},
	"cayenne pepper":      {},
	"celeriac":            {},
	"celery":              {},
	"chamomile":           {},
	"chard":               {},
	"chayote":             {},
	"chickpea":            {},
	"chives":              {},
	"cilantro":            {},
	"collard green":       {},
	"corn":                {},
	"corn salad":          {},
	"courgette":           {},
	"cucumber":            {},
	"daikon":              {},
	"delicata":            {},
	"dill":                {},
	"eggplant":            {},
	"endive":              {},
	"fennel":              {},
	"fiddlehead":          {},
	"frisee":              {},
	"garlic":              {},
	"gem squash":          {},
	"ginger":              {},
	"green bean":          {},
	"green pepper":        {},
	"habanero":            {},
	"herbs and spice":     {},
	"horseradish":         {},
	"hubbard squash":      {},
	"jalapeno":            {},
	"jerusalem artichoke": {},
	"jicama":              {},
	"kale":                {},
	"kidney bean":         {},
	"kohlrabi":            {},
	"lavender":            {},
	"leek ":               {},
	"legume":              {},
	"lemon grass":         {},
	"lentils":             {},
	"lettuce":             {},
	"lima bean":           {},
	"mamey":               {},
	"mangetout":           {},
	"marjoram":            {},
	"mung bean":           {},
	"mushroom":            {},
	"mustard green":       {},
	"navy bean":           {},
	"new zealand spinach": {},
	"nopale":              {},
	"okra":                {},
	"onion":               {},
	"oregano":             {},
	"paprika":             {},
	"parsley":             {},
	"parsnip":             {},
	"patty pan":           {},
	"pea":                 {},
	"pinto bean":          {},
	"potato":              {},
	"pumpkin":             {},
	"radicchio":           {},
	"radish":              {},
	"rhubarb":             {},
	"rosemary":            {},
	"runner bean":         {},
	"rutabaga":            {},
	"sage":                {},
	"scallion":            {},
	"shallot":             {},
	"skirret":             {},
	"snap pea":            {},
	"soy bean":            {},
	"spaghetti squash":    {},
	"spinach":             {},
	"squash ":             {},
	"sweet potato":        {},
	"tabasco pepper":      {},
	"taro":                {},
	"tat soi":             {},
	"thyme":               {},
	"topinambur":          {},
	"tubers":              {},
	"turnip":              {},
	"wasabi":              {},
	"water chestnut":      {},
	"watercress":          {},
	"white radish":        {},
	"yam":                 {},
	"zucchini":            {},
}

var corpusIngredients = []string{" can cream of chicken soup ",
	" cream of mushroom soup ",
	" graham cracker crumbs ",
	" salt and black pepper ",
	" heavy whipping cream ",
	" vegetable shortening ",
	" worcestershire sauce ",
	" apple cider vinegar ",
	" confectioners sugar ",
	" limnophila aromatica ",
	" vietnamese coriander ",
	" jerusalem artichoke ",
	" new zealand spinach ",
	" white wine vinegar ",
	" black peppercorns ",
	" bolivian coriander ",
	" cream of mushroom ",
	" grains of paradise ",
	" green bell pepper ",
	" houttuynia cordata ",
	" italian seasoning ",
	" mozzarella cheese ",
	" pumpkin pie spice ",
	" red pepper flakes ",
	" vanilla ice cream ",
	" artichoke hearts ",
	" balsamic vinegar ",
	" broccoli florets ",
	" butternut squash ",
	" east asian pepper ",
	" hot pepper sauce ",
	" jalapeno chilies ",
	" purple mangosteen ",
	" red wine vinegar ",
	" alligator pepper ",
	" american cheese ",
	" black peppercorn ",
	" brazilian pepper ",
	" butternut squash ",
	" cherry tomatoes ",
	" chicken breasts ",
	" cilantro leaves ",
	" cream of tartar ",
	" evaporated milk ",
	" flour tortillas ",
	" graham crackers ",
	" greater galangal ",
	" green peppercorn ",
	" lasagna noodles ",
	" mzarella cheese ",
	" parmesan cheese ",
	" pineapple juice ",
	" pomegranate seed ",
	" red bell pepper ",
	" salt and pepper ",
	" spaghetti squash ",
	" tasmanian pepper ",
	" vegetable broth ",
	" vegetable stock ",
	" water chestnuts ",
	" white peppercorn ",
	" yellow cake mix ",
	" yellow cornmeal ",
	" almond extract ",
	" aromatic ginger ",
	" barbecue sauce ",
	" brussels sprout ",
	" bunium persicum ",
	" cayenne pepper ",
	" cheddar cheese ",
	" chicken breast ",
	" chicken thighs ",
	" corn tortillas ",
	" cottage cheese ",
	" dark chocolate ",
	" dry white wine ",
	" grains of selim ",
	" herbs and spice ",
	" indian bay leaf ",
	" lesser galangal ",
	" parsley flakes ",
	" parsley leaves ",
	" peruvian pepper ",
	" powdered sugar ",
	" ricotta cheese ",
	" seasoning salt ",
	" smoked paprika ",
	" sweet potatoes ",
	" teriyaki sauce ",
	" whipping cream ",
	" alfalfa sprout ",
	" baking powder ",
	" black cardamom ",
	" black-eyed pea ",
	" cayenne pepper ",
	" chicken broth ",
	" chicken stock ",
	" cider vinegar ",
	" cooking spray ",
	" coriander leaf ",
	" coriander seed ",
	" dijon mustard ",
	" garlic powder ",
	" green chilies ",
	" green peppers ",
	" half and half ",
	" hubbard squash ",
	" nigella sativa ",
	" peanut butter ",
	" plum tomatoes ",
	" potato starch ",
	" sichuan pepper ",
	" tabasco pepper ",
	" tabasco sauce ",
	" vegetable oil ",
	" voatsiperifery ",
	" water chestnut ",
	" whipped cream ",
	" white vinegar ",
	" almond flour ",
	" baby spinach ",
	" banana squash ",
	" basil leaves ",
	" black mustard ",
	" black olives ",
	" black pepper ",
	" borlotti bean ",
	" bread crumbs ",
	" brown mustard ",
	" celery powder ",
	" chicken soup ",
	" chili powder ",
	" coconut milk ",
	" collard green ",
	" corn kernels ",
	" cream cheese ",
	" curry powder ",
	" dry red wine ",
	" garam masala ",
	" garlic chives ",
	" green onions ",
	" green pepper ",
	" juniper berry ",
	" kidney beans ",
	" lemon verbena ",
	" litsea cubeba ",
	" marshmallows ",
	" mustard green ",
	" onion powder ",
	" orange juice ",
	" purple onion ",
	" rice vinegar ",
	" sesame seeds ",
	" strawberries ",
	" swiss cheese ",
	" thyme leaves ",
	" tomato juice ",
	" tomato paste ",
	" tomato sauce ",
	" walnut flour ",
	" white mustard ",
	" white pepper ",
	" yellow onion ",
	" acorn squash ",
	" almond meal ",
	" almond milk ",
	" apple juice ",
	" baking soda ",
	" bell pepper ",
	" black beans ",
	" blackcurrant ",
	" blood orange ",
	" blue cheese ",
	" blueberries ",
	" bread flour ",
	" breadcrumbs ",
	" broccoflower ",
	" brown sugar ",
	" canary melon ",
	" cauliflower ",
	" celery seed ",
	" chili pepper ",
	" chili sauce ",
	" coconut oil ",
	" cooking oil ",
	" corn starch ",
	" cranberries ",
	" dry mustard ",
	" egg noodles ",
	" feta cheese ",
	" garlic salt ",
	" goat cheese ",
	" green beans ",
	" green onion ",
	" green pepper ",
	" heavy cream ",
	" horseradish ",
	" lemon juice ",
	" lemon myrtle ",
	" mango-ginger ",
	" maple syrup ",
	" mint leaves ",
	" orange peel ",
	" passionfruit ",
	" raspberries ",
	" red peppers ",
	" rice noodle ",
	" sarsaparilla ",
	" sweet onion ",
	" sweet potato ",
	" tomato soup ",
	" white onion ",
	" white radish ",
	" applesauce ",
	" bay leaves ",
	" bean sprout ",
	" beef broth ",
	" beef stock ",
	" bell pepper ",
	" boysenberry ",
	" broad beans ",
	" brown rice ",
	" buttermilk ",
	" canola oil ",
	" cauliflower ",
	" celery seed ",
	" corn flour ",
	" corn syrup ",
	" cornstarch ",
	" cumin seed ",
	" dragonfruit ",
	" dry sherry ",
	" fish sauce ",
	" horseradish ",
	" huckleberry ",
	" kidney bean ",
	" lemon grass ",
	" lemon peel ",
	" lemon rind ",
	" lime juice ",
	" long pepper ",
	" matzo meal ",
	" mayonnaise ",
	" peanut oil ",
	" peppermint ",
	" pomegranate ",
	" prosciutto ",
	" red pepper ",
	" rice flour ",
	" runner bean ",
	" salal berry ",
	" sesame oil ",
	" shortening ",
	" sour cream ",
	" white rice ",
	" white wine ",
	" asafoetida ",
	" asparagus ",
	" azuki bean ",
	" black bean ",
	" blackberry ",
	" breadfruit ",
	" cantaloupe ",
	" chickpeas ",
	" chocolate ",
	" clementine ",
	" cloudberry ",
	" cool whip ",
	" coriander ",
	" corn salad ",
	" curry leaf ",
	" dill weed ",
	" dried lime ",
	" elderberry ",
	" fiddlehead ",
	" fingerroot ",
	" gem squash ",
	" goji berry ",
	" gooseberry ",
	" grapefruit ",
	" green bean ",
	" hamburger ",
	" hazelnuts ",
	" hoja santa ",
	" holy basil ",
	" hot sauce ",
	" kiwi fruit ",
	" lemon balm ",
	" margarine ",
	" mushrooms ",
	" oat flour ",
	" olive oil ",
	" pie crust ",
	" pie shell ",
	" pine nuts ",
	" pineapple ",
	" pinto bean ",
	" poppy seed ",
	" red onion ",
	" redcurrant ",
	" rock melon ",
	" salad oil ",
	" scallions ",
	" skim milk ",
	" soy sauce ",
	" spaghetti ",
	" star anise ",
	" star fruit ",
	" strawberry ",
	" thai basil ",
	" tonka bean ",
	" topinambur ",
	" ugli fruit ",
	" watercress ",
	" watermelon ",
	" allspice ",
	" artichoke ",
	" asparagus ",
	" aubergine ",
	" bay leaf ",
	" blueberry ",
	" broccoli ",
	" calabrese ",
	" cardamom ",
	" chamomile ",
	" cherimoya ",
	" cherries ",
	" cilantro ",
	" cinnamon ",
	" cornmeal ",
	" courgette ",
	" crabmeat ",
	" cranberry ",
	" cucumber ",
	" dill seed ",
	" dressing ",
	" eggplant ",
	" fenugreek ",
	" jackfruit ",
	" kinh gioi ",
	" lima bean ",
	" liquorice ",
	" mandarine ",
	" mangetout ",
	" marjoram ",
	" molasses ",
	" mung bean ",
	" navy bean ",
	" nectarine ",
	" patty pan ",
	" persimmon ",
	" pineapple ",
	" potatoes ",
	" radicchio ",
	" raspberry ",
	" red wine ",
	" rosemary ",
	" sassafras ",
	" semolina ",
	" shallots ",
	" tamarillo ",
	" tangerine ",
	" tarragon ",
	" tomatoes ",
	" turmeric ",
	" zucchini ",
	" allspice ",
	" almonds ",
	" amaranth ",
	" angelica ",
	" arugula ",
	" avocado ",
	" bananas ",
	" bay leaf ",
	" bilberry ",
	" bok choy ",
	" broccoli ",
	" cabbage ",
	" cardamom ",
	" carrots ",
	" cashews ",
	" cayenne ",
	" celeriac ",
	" chicken ",
	" chickpea ",
	" cilantro ",
	" cinnamon ",
	" coconut ",
	" cucumber ",
	" culantro ",
	" delicata ",
	" eggplant ",
	" gelatin ",
	" habanero ",
	" honeydew ",
	" jalapeno ",
	" ketchup ",
	" kohlrabi ",
	" korarima ",
	" lavender ",
	" lettuce ",
	" marjoram ",
	" mulberry ",
	" mushroom ",
	" mustard ",
	" oatmeal ",
	" oranges ",
	" oregano ",
	" paprika ",
	" parsley ",
	" peaches ",
	" peanuts ",
	" physalis ",
	" pumpkin ",
	" raisins ",
	" rambutan ",
	" rosemary ",
	" rutabaga ",
	" sausage ",
	" scallion ",
	" shallot ",
	" snap pea ",
	" soy bean ",
	" spinach ",
	" tamarind ",
	" tarragon ",
	" turmeric ",
	" vanilla ",
	" vinegar ",
	" walnuts ",
	" woodruff ",
	" zucchini ",
	" amchoor ",
	" apples ",
	" apricot ",
	" arugula ",
	" avocado ",
	" banana ",
	" brandy ",
	" butter ",
	" cabbage ",
	" camphor ",
	" capers ",
	" caraway ",
	" carrot ",
	" catsup ",
	" celery ",
	" charoli ",
	" chayote ",
	" cheese ",
	" chervil ",
	" chives ",
	" coconut ",
	" coffee ",
	" currant ",
	" epazote ",
	" fennel ",
	" garlic ",
	" ginger ",
	" kumquat ",
	" lemons ",
	" lentils ",
	" lettuce ",
	" mitsuba ",
	" mugwort ",
	" nigella ",
	" njangsa ",
	" nutmeg ",
	" olives ",
	" onions ",
	" orange ",
	" oregano ",
	" paprika ",
	" parsley ",
	" parsnip ",
	" pecans ",
	" pepper ",
	" perilla ",
	" potato ",
	" pumpkin ",
	" quinoa ",
	" radhuni ",
	" rhubarb ",
	" saffron ",
	" salmon ",
	" satsuma ",
	" shallot ",
	" shrimp ",
	" skirret ",
	" spinach ",
	" squash ",
	" tat soi ",
	" tomato ",
	" turkey ",
	" vanilla ",
	" yogurt ",
	" zedoary ",
	" zereshk ",
	" ajwain ",
	" aonori ",
	" apple ",
	" bacon ",
	" banana ",
	" basil ",
	" beans ",
	" beets ",
	" borage ",
	" bread ",
	" carrot ",
	" cassia ",
	" celery ",
	" chenpi ",
	" cherry ",
	" chile ",
	" chives ",
	" cicely ",
	" clove ",
	" cocoa ",
	" cream ",
	" cumin ",
	" daikon ",
	" damson ",
	" dates ",
	" durian ",
	" endive ",
	" feijoa ",
	" fennel ",
	" flour ",
	" frisee ",
	" fruit ",
	" garlic ",
	" ginger ",
	" golpar ",
	" honey ",
	" hyssop ",
	" jambul ",
	" jicama ",
	" jujube ",
	" leeks ",
	" legume ",
	" lemon ",
	" loquat ",
	" lovage ",
	" lychee ",
	" mahlab ",
	" mango ",
	" mastic ",
	" nopale ",
	" nutmeg ",
	" onion ",
	" orange ",
	" pamelo ",
	" papaya ",
	" pasta ",
	" pears ",
	" pomelo ",
	" potato ",
	" quince ",
	" radish ",
	" raisin ",
	" salad ",
	" salsa ",
	" sansho ",
	" sauce ",
	" savory ",
	" sesame ",
	" sorrel ",
	" spray ",
	" steak ",
	" sugar ",
	" thyme ",
	" tubers ",
	" turnip ",
	" vodka ",
	" wasabi ",
	" water ",
	" yeast ",
	" anise ",
	" apple ",
	" basil ",
	" beef ",
	" beer ",
	" boldo ",
	" chard ",
	" clove ",
	" corn ",
	" cress ",
	" cubeb ",
	" cumin ",
	" dill ",
	" eggs ",
	" flax ",
	" grape ",
	" guava ",
	" jimbu ",
	" kale ",
	" kokum ",
	" lamb ",
	" leek ",
	" lemon ",
	" lime ",
	" mamey ",
	" mango ",
	" milk ",
	" mint ",
	" nuts ",
	" oats ",
	" olive ",
	" onion ",
	" peach ",
	" peas ",
	" pork ",
	" rice ",
	" sage ",
	" salt ",
	" shiso ",
	" soda ",
	" sumac ",
	" thru ",
	" thyme ",
	" tofu ",
	" tuna ",
	" uzazi ",
	" beet ",
	" corn ",
	" date ",
	" dill ",
	" ham ",
	" hemp ",
	" ice ",
	" kale ",
	" lime ",
	" mace ",
	" mint ",
	" oil ",
	" okra ",
	" pear ",
	" plum ",
	" rose ",
	" sage ",
	" salt ",
	" taro ",
	" yuzu ",
	" zest ",
	" fig ",
	" nut ",
	" pea ",
	" rue ",
	" yam "}
var corpusMeasures = []string{" tablespoons ",
	" milliliter ",
	" tablespoon ",
	" teaspoons ",
	" teaspoon ",
	" canned ",
	" ounces ",
	" pounds ",
	" quarts ",
	" grams ",
	" ounce ",
	" pints ",
	" pound ",
	" quart ",
	" tbsps ",
	" cans ",
	" cups ",
	" gram ",
	" pint ",
	" tbsp ",
	" tsps ",
	" can ",
	" cup ",
	" tbl ",
	" tsp ",
	" ml ",
	" oz ",
	" c ",
	" g "}
var corpusNumbers = []string{" 1/2 ",
	" 1/3 ",
	" 1/4 ",
	" 1/5 ",
	" 1/6 ",
	" 1/7 ",
	" 1/8 ",
	" 2/3 ",
	" 2/5 ",
	" 2/7 ",
	" 3/4 ",
	" 3/8 ",
	" 4/5 ",
	" 5/8 ",
	" 7/8 ",
	" 1 ",
	" 2 ",
	" 3 ",
	" 4 ",
	" 5 ",
	" 6 ",
	" 7 ",
	" 8 ",
	" 9 ",
	" 10 ",
	" 11 ",
	" 12 ",
	" 13 ",
	" 14 ",
	" 15 ",
	" 16 ",
	" 17 ",
	" 18 ",
	" 19 ",
	" 20 ",
	" ⅜ ",
	" ⅝ ",
	" ⅞ ",
	" ½ ",
	" ¼ ",
	" ⅛ ",
	" ¾ ",
	" ⅔ ",
	" ⅓ "}
var corpusDirections = []string{" 1 ",
	" 10 ",
	" 15 ",
	" 2 ",
	" 25 ",
	" 375 ",
	" 6 ",
	" 8 ",
	" 9x13 ",
	" a ",
	" about ",
	" achieve ",
	" add ",
	" additional ",
	" an ",
	" and ",
	" arrange ",
	" assemble ",
	" bake ",
	" baking ",
	" basil ",
	" beef ",
	" before ",
	" boil ",
	" boiling ",
	" bottom ",
	" bowl ",
	" bring ",
	" brown ",
	" browned ",
	" cheese ",
	" coconut ",
	" cold ",
	" color ",
	" combine ",
	" cook ",
	" cooking ",
	" cool ",
	" cover ",
	" covered ",
	" crushed ",
	" cup ",
	" cups ",
	" degrees ",
	" dish ",
	" distributed ",
	" does ",
	" drain ",
	" dutch ",
	" eggs ",
	" either ",
	" even ",
	" evenly ",
	" f ",
	" fennel ",
	" foil ",
	" for ",
	" from ",
	" garlic ",
	" ground ",
	" half ",
	" heat ",
	" hours ",
	" in ",
	" inch ",
	" into ",
	" italian ",
	" large ",
	" lasagna ",
	" layers ",
	" lengthwise ",
	" lightly ",
	" make ",
	" meat ",
	" medium ",
	" minutes ",
	" mix ",
	" mixing ",
	" mixture ",
	" mozzarella ",
	" noodles ",
	" not ",
	" nuts ",
	" oats ",
	" occasionally ",
	" of ",
	" onion ",
	" or ",
	" oven ",
	" over ",
	" pans ",
	" parmesan ",
	" parsley ",
	" paste ",
	" pepper ",
	" pot ",
	" preheat ",
	" preheated ",
	" prevent ",
	" raisins ",
	" remaining ",
	" remove ",
	" repeat ",
	" ricotta ",
	" rinse ",
	" salt ",
	" salted ",
	" sauce ",
	" sausage ",
	" season ",
	" seasoning ",
	" seeds ",
	" serving ",
	" sheet ",
	" simmer ",
	" slices ",
	" spoon ",
	" spray ",
	" spread ",
	" sprinkle ",
	" sticking ",
	" stir ",
	" stirring ",
	" sugar ",
	" sure ",
	" tablespoon ",
	" tablespoons ",
	" teaspoon ",
	" the ",
	" third ",
	" to ",
	" tomato ",
	" tomatoes ",
	" top ",
	" touch ",
	" transfer ",
	" until ",
	" water ",
	" well ",
	" with ",
	" ¼ ",
	" ½ "}
var corpusDirectionsNeg = []string{" a ",
	" all ",
	" also ",
	" amount ",
	" and ",
	" be ",
	" best ",
	" brand ",
	" butter ",
	" c ",
	" calcium ",
	" calories ",
	" carbohydrates ",
	" chip ",
	" chocolate ",
	" cholesterol ",
	" cookie ",
	" cookies. ",
	" costco ",
	" dessert ",
	" dinner ",
	" dough ",
	" ensure ",
	" equally ",
	" ever ",
	" excellent ",
	" fat ",
	" fiber ",
	" for ",
	" from ",
	" great! ",
	" have ",
	" i ",
	" images ",
	" iron ",
	" it's ",
	" just ",
	" kirkland ",
	" liking! ",
	" make ",
	" per ",
	" potassium ",
	" protein ",
	" recipe ",
	" recommend ",
	" results. ",
	" salted ",
	" saturated ",
	" serving ",
	" servings ",
	" sodium ",
	" tasting ",
	" text ",
	" that ",
	" the ",
	" then ",
	" these ",
	" tillamook ",
	" to ",
	" unsalted ",
	" use ",
	" used ",
	" vitamin ",
	" with ",
	" would ",
	" yield ",
	" your "}

type fractionNumber struct {
	fractionString string
	value          float64
}

var corpusFractionNumberMap = map[string]fractionNumber{
	"¼": {"1/4", 0.2500000000},
	"½": {"1/2", 0.5000000000},
	"¾": {"3/4", 0.7500000000},
	"⅓": {"1/3", 0.3333333333},
	"⅔": {"2/3", 0.6666666667},
	"⅛": {"1/8", 0.1250000000},
	"⅜": {"3/8", 0.3750000000},
	"⅝": {"5/8", 0.6250000000},
	"⅞": {"7/8", 0.8750000000},
}

var corpusMeasuresMap = map[string]string{
	"c":           "cup",
	"can":         "can",
	"canned":      "can",
	"cans":        "can",
	"cup":         "cup",
	"cups":        "cup",
	"g":           "gram",
	"gram":        "gram",
	"grams":       "gram",
	"milliliter":  "milliliter",
	"ml":          "milliliter",
	"ounce":       "ounce",
	"ounces":      "ounce",
	"oz":          "ounce",
	"pint":        "pint",
	"pints":       "pint",
	"pound":       "pound",
	"pounds":      "pound",
	"quart":       "quart",
	"quarts":      "quart",
	"tablespoon":  "tbl",
	"tablespoons": "tbl",
	"tbl":         "tbl",
	"tbsp":        "tbl",
	"tbsps":       "tbl",
	"teaspoon":    "tsp",
	"teaspoons":   "tsp",
	"tsp":         "tsp",
	"tsps":        "tsp",
}

var densities = map[string]float64{
	"almond milk":            245.5000000000,
	"almonds":                115.1000000000,
	"apple":                  144.6000000000,
	"apple juice":            243.5000000000,
	"apples":                 152.2000000000,
	"applesauce":             247.2000000000,
	"arugula":                20.0000000000,
	"asparagus":              204.3000000000,
	"avocado":                218.0000000000,
	"bacon":                  156.8000000000,
	"balsamic vinegar":       255.0000000000,
	"banana":                 159.1000000000,
	"bananas":                158.3000000000,
	"barbecue sauce":         279.2000000000,
	"basil":                  191.0000000000,
	"beans":                  189.5000000000,
	"beef":                   212.8000000000,
	"beef broth":             225.2000000000,
	"beef stock":             245.3000000000,
	"beer":                   237.0000000000,
	"beets":                  195.8000000000,
	"black beans":            194.5000000000,
	"blue cheese":            135.0000000000,
	"blueberries":            206.5000000000,
	"bread":                  145.8000000000,
	"bread crumbs":           114.0000000000,
	"broccoli":               154.6000000000,
	"brown rice":             174.2000000000,
	"brown sugar":            160.9000000000,
	"butter":                 227.0000000000,
	"buttermilk":             245.0000000000,
	"cabbage":                112.1000000000,
	"canola oil":             180.5000000000,
	"carrot":                 155.0000000000,
	"carrots":                178.6000000000,
	"catsup":                 240.0000000000,
	"cauliflower":            142.2000000000,
	"celery":                 182.8000000000,
	"cheddar cheese":         179.0000000000,
	"cheese":                 144.7000000000,
	"cherries":               208.2000000000,
	"chicken":                194.0000000000,
	"chicken broth":          203.2000000000,
	"chicken soup":           248.3000000000,
	"chicken stock":          240.0000000000,
	"chickpeas":              193.3000000000,
	"chile":                  37.0000000000,
	"chili sauce":            259.0000000000,
	"chives":                 3.2000000000,
	"chocolate":              171.2000000000,
	"cider vinegar":          239.0000000000,
	"cinnamon":               53.5000000000,
	"cocoa":                  51.4000000000,
	"coconut":                261.9000000000,
	"coconut milk":           236.5000000000,
	"coconut oil":            218.0000000000,
	"coffee":                 248.7000000000,
	"coriander":              16.0000000000,
	"corn":                   164.1000000000,
	"corn flour":             123.4000000000,
	"cornmeal":               144.7000000000,
	"cornstarch":             128.0000000000,
	"cottage cheese":         216.3000000000,
	"cranberries":            123.3000000000,
	"cream":                  187.8000000000,
	"cream cheese":           237.3000000000,
	"cream of mushroom soup": 249.9000000000,
	"cucumber":               146.3000000000,
	"dates":                  147.0000000000,
	"dill weed":              8.9000000000,
	"dressing":               245.4000000000,
	"egg noodles":            114.2000000000,
	"eggplant":               104.0000000000,
	"evaporated milk":        252.0000000000,
	"fennel":                 87.0000000000,
	"feta cheese":            150.0000000000,
	"flour":                  113.7000000000,
	"fruit":                  225.8000000000,
	"garlic":                 144.9000000000,
	"gelatin":                195.0000000000,
	"ginger":                 96.0000000000,
	"graham crackers":        84.0000000000,
	"green beans":            220.7000000000,
	"green chilies":          241.0000000000,
	"green onions":           71.0000000000,
	"green pepper":           130.0000000000,
	"green peppers":          238.5000000000,
	"ham":                    250.4000000000,
	"hamburger":              244.0000000000,
	"hazelnuts":              108.3000000000,
	"honey":                  59.1000000000,
	"ice":                    187.3000000000,
	"kale":                   117.0000000000,
	"kidney beans":           194.6000000000,
	"leeks":                  75.0000000000,
	"lemon":                  196.9000000000,
	"lemon juice":            244.0000000000,
	"lemons":                 212.0000000000,
	"lettuce":                49.2000000000,
	"lime":                   198.0000000000,
	"lime juice":             244.0000000000,
	"mango":                  251.0000000000,
	"maple syrup":            322.0000000000,
	"margarine":              225.7000000000,
	"marshmallows":           39.3000000000,
	"mayonnaise":             230.6000000000,
	"milk":                   220.3000000000,
	"molasses":               337.0000000000,
	"mushrooms":              116.8000000000,
	"mustard":                153.0000000000,
	"nuts":                   150.6000000000,
	"oatmeal":                60.5000000000,
	"oats":                   81.0000000000,
	"oil":                    207.3000000000,
	"olive oil":              216.0000000000,
	"onion":                  152.6000000000,
	"onions":                 176.5000000000,
	"orange":                 237.7000000000,
	"orange juice":           251.2000000000,
	"oranges":                176.0000000000,
	"parmesan":               126.7000000000,
	"parmesan cheese":        100.0000000000,
	"parsley":                32.8000000000,
	"pasta":                  116.0000000000,
	"peaches":                229.9000000000,
	"peanut butter":          154.4000000000,
	"peanut oil":             216.0000000000,
	"peanuts":                135.8000000000,
	"pears":                  190.6000000000,
	"peas":                   171.1000000000,
	"pecans":                 107.0000000000,
	"pepper":                 124.0000000000,
	"pie crust":              129.0000000000,
	"pine nuts":              135.0000000000,
	"pineapple":              235.3000000000,
	"pineapple juice":        250.0000000000,
	"pork":                   146.7000000000,
	"potato":                 191.7000000000,
	"potatoes":               180.4000000000,
	"pumpkin":                134.1000000000,
	"quinoa":                 177.5000000000,
	"raisins":                140.2000000000,
	"raspberries":            179.8000000000,
	"red pepper":             125.0000000000,
	"red wine vinegar":       239.0000000000,
	"rice":                   89.2000000000,
	"ricotta cheese":         247.0000000000,
	"salad":                  207.0000000000,
	"salad oil":              214.0000000000,
	"salmon":                 177.0000000000,
	"salsa":                  250.0000000000,
	"salt":                   292.0000000000,
	"sauce":                  251.4000000000,
	"sausage":                140.3000000000,
	"scallions":              100.0000000000,
	"semolina":               107.6000000000,
	"sesame oil":             218.0000000000,
	"sesame seeds":           144.0000000000,
	"shallots":               14.4000000000,
	"shortening":             207.3000000000,
	"shrimp":                 200.2000000000,
	"skim milk":              245.0000000000,
	"soda":                   57.5000000000,
	"sour cream":             237.2000000000,
	"soy sauce":              243.5000000000,
	"spaghetti":              151.4000000000,
	"spinach":                164.1000000000,
	"spray":                  247.0000000000,
	"sugar":                  205.8000000000,
	"sweet potatoes":         224.0000000000,
	"swiss cheese":           137.7000000000,
	"tofu":                   250.0000000000,
	"tomato":                 230.5000000000,
	"tomato juice":           241.2000000000,
	"tomato sauce":           217.1000000000,
	"tomato soup":            206.2000000000,
	"tomatoes":               177.6000000000,
	"tuna":                   146.0000000000,
	"turkey":                 206.6000000000,
	"vanilla":                175.9000000000,
	"vegetable broth":        225.3000000000,
	"vegetable oil":          211.5000000000,
	"vegetable shortening":   205.0000000000,
	"vinegar":                238.0000000000,
	"walnuts":                95.0000000000,
	"water":                  243.2000000000,
	"whipped cream":          70.0000000000,
	"white bread":            41.2000000000,
	"white rice":             172.9000000000,
	"worcestershire sauce":   275.0000000000,
	"yellow cornmeal":        143.3000000000,
	"yogurt":                 211.8000000000,
	"zucchini":               194.4000000000,
}
