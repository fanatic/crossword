package main

import "time"

type BoardLayout struct {
	Admin   bool `json:"admin"`
	Answers struct {
		Across []string `json:"across"`
		Down   []string `json:"down"`
	} `json:"answers"`
	Author   string      `json:"author"`
	Autowrap interface{} `json:"autowrap"`
	Bbars    interface{} `json:"bbars"`
	Circles  interface{} `json:"circles"`
	Clues    struct {
		Across []string `json:"across"`
		Down   []string `json:"down"`
	} `json:"clues"`
	Code            interface{} `json:"code"`
	Copyright       string      `json:"copyright"`
	Date            string      `json:"date"`
	Dow             string      `json:"dow"`
	Downmap         interface{} `json:"downmap"`
	Editor          string      `json:"editor"`
	Grid            []string    `json:"grid"`
	Gridnums        []int       `json:"gridnums"`
	Hastitle        bool        `json:"hastitle"`
	Hold            interface{} `json:"hold"`
	ID              interface{} `json:"id"`
	ID2             interface{} `json:"id2"`
	Interpretcolors interface{} `json:"interpretcolors"`
	Jnotes          interface{} `json:"jnotes"`
	Key             interface{} `json:"key"`
	Mini            interface{} `json:"mini"`
	Notepad         interface{} `json:"notepad"`
	Publisher       string      `json:"publisher"`
	Rbars           interface{} `json:"rbars"`
	Shadecircles    interface{} `json:"shadecircles"`
	Size            struct {
		Cols int `json:"cols"`
		Rows int `json:"rows"`
	} `json:"size"`
	Title   string      `json:"title"`
	Track   interface{} `json:"track"`
	Type    interface{} `json:"type"`
	Uniclue bool        `json:"uniclue"`
	Valid   bool        `json:"valid"`
}

var sampleBoard = BoardLayout{
	Answers: struct {
		Across []string `json:"across"`
		Down   []string `json:"down"`
	}{
		Across: []string{
			"RATPACK",
			"NBAJAM",
			"EURASIAN",
			"YESIDO",
			"ABUDHABI",
			"SEWNON",
			"PAIR",
			"OLGA",
			"NAGAT",
			"EDSEL",
			"OHMS",
			"NONA",
			"REM",
			"EXOTICA",
			"INN",
			"HAREM",
			"RUSSIA",
			"ONEDAYATATIME",
			"GROMIT",
			"RAPID",
			"ADS",
			"NEMESES",
			"PAC",
			"NELL",
			"DAFT",
			"MARSH",
			"GROAN",
			"SUEZ",
			"ROSA",
			"SOUPED",
			"EBENEZER",
			"TUCSON",
			"LUKEWARM",
			"ATHENA",
			"DETECTS",
		},
		Down: []string{
			"REAPER",
			"AUBADE",
			"TRUISM",
			"PADRE",
			"ASH",
			"CIAO",
			"KABLOOEY",
			"NYS",
			"BEEN",
			"ASWAN",
			"JINGOISM",
			"ADOANNIE",
			"MONTANA",
			"NIGHTMAREFUEL",
			"AMI",
			"LEADIN",
			"SCRAPE",
			"XRATED",
			"AUTISM",
			"HEM",
			"SID",
			"ORDEROUT",
			"NOSLOUCH",
			"TASTEBUD",
			"GANGSTA",
			"MAS",
			"PROZAC",
			"ASSERT",
			"CHARMS",
			"LAPSE",
			"AREWE",
			"NEON",
			"ZEKE",
			"DNA",
			"NET",
		},
	},
	Author: "Neville Fogarty and Doug Peterson",
	Clues: struct {
		Across []string `json:"across"`
		Down   []string `json:"down"`
	}{
		Across: []string{
			"1. Group in the original \"Ocean's 11\"",
			"8. Classic arcade game with lots of shooting",
			"14. Like Istanbul",
			"16. Emphatic admission",
			"17. First world capital, alphabetically",
			"18. Like clothes buttons, generally",
			"19. Unexciting poker holding",
			"20. 2008 Bond girl Kurylenko",
			"22. Bedevil",
			"23. Car once promoted with the line \"The thrill starts with the grille\"",
			"25. Speaker units",
			"27. Prefix with -gon",
			"28. Nocturnal acronym",
			"29. Strange things",
			"32. Super 8, e.g.",
			"33. Group of female seals",
			"34. Powerhouse in Olympic weightlifting",
			"36. Gradually",
			"39. Animated character who graduated from Dogwarts University",
			"40. The \"R\" of 28-Across",
			"41. Circular parts",
			"42. Formidable opponents",
			"44. Campaign aid",
			"47. \"The Old Curiosity Shop\" girl",
			"49. Touched",
			"50. Rail center?",
			"52. Express stress, in a way",
			"54. Gulf of ___",
			"56. Santa ___, Calif.",
			"57. Juiced (up)",
			"59. Jacob's partner in \"A Christmas Carol\"",
			"61. City nicknamed \"The Old Pueblo\"",
			"62. So-so, as support",
			"63. Acropolis figure",
			"64. Spots",
		},
		Down: []string{
			"1. One going against the grain?",
			"2. Poem greeting the dawn",
			"3. \"What's past is past,\" e.g.",
			"4. Giant competitor",
			"5. Last name of cosmetics giant Mary Kay",
			"6. \"See ya\"",
			"7. Bad way to go",
			"8. Buffalo's home: Abbr.",
			"9. Has-___",
			"10. Source of stone used to build the ancient Egyptian pyramids",
			"11. Flag-waving and such",
			"12. Musical \"girl who cain't say no\"",
			"13. Joe known as \"The Comeback Kid\"",
			"15. Cause of bad dreams, in modern lingo",
			"21. Follower of bon or mon",
			"24. Show immediately preceding another",
			"26. Scuffle",
			"30. For adults only",
			"31. Special-education challenge",
			"33. Bottom line?",
			"35. Tom Sawyer's half brother",
			"36. Request for food delivery",
			"37. Someone who's pretty darn good",
			"38. It could be on the tip of your tongue",
			"39. ___ rap",
			"43. More, in México",
			"44. O.C.D. fighter, maybe",
			"45. Put forth",
			"46. Enamors",
			"48. Small slip",
			"51. \"___ done now?\"",
			"53. Superbright",
			"55. \"The Wizard of Oz\" farmhand",
			"58. Helicases split it",
			"60. Court divider",
		},
	},
	Copyright: "2018, The New York Times",
	Date:      "3/9/2018",
	Dow:       "Friday",
	Editor:    "Will Shortz",
	Grid: []string{
		"R",
		"A",
		"T",
		"P",
		"A",
		"C",
		"K",
		".",
		".",
		"N",
		"B",
		"A",
		"J",
		"A",
		"M",
		"E",
		"U",
		"R",
		"A",
		"S",
		"I",
		"A",
		"N",
		".",
		"Y",
		"E",
		"S",
		"I",
		"D",
		"O",
		"A",
		"B",
		"U",
		"D",
		"H",
		"A",
		"B",
		"I",
		".",
		"S",
		"E",
		"W",
		"N",
		"O",
		"N",
		"P",
		"A",
		"I",
		"R",
		".",
		"O",
		"L",
		"G",
		"A",
		".",
		"N",
		"A",
		"G",
		"A",
		"T",
		"E",
		"D",
		"S",
		"E",
		"L",
		".",
		"O",
		"H",
		"M",
		"S",
		".",
		"N",
		"O",
		"N",
		"A",
		"R",
		"E",
		"M",
		".",
		"E",
		"X",
		"O",
		"T",
		"I",
		"C",
		"A",
		".",
		"I",
		"N",
		"N",
		".",
		".",
		".",
		"H",
		"A",
		"R",
		"E",
		"M",
		".",
		"R",
		"U",
		"S",
		"S",
		"I",
		"A",
		".",
		"O",
		"N",
		"E",
		"D",
		"A",
		"Y",
		"A",
		"T",
		"A",
		"T",
		"I",
		"M",
		"E",
		".",
		"G",
		"R",
		"O",
		"M",
		"I",
		"T",
		".",
		"R",
		"A",
		"P",
		"I",
		"D",
		".",
		".",
		".",
		"A",
		"D",
		"S",
		".",
		"N",
		"E",
		"M",
		"E",
		"S",
		"E",
		"S",
		".",
		"P",
		"A",
		"C",
		"N",
		"E",
		"L",
		"L",
		".",
		"D",
		"A",
		"F",
		"T",
		".",
		"M",
		"A",
		"R",
		"S",
		"H",
		"G",
		"R",
		"O",
		"A",
		"N",
		".",
		"S",
		"U",
		"E",
		"Z",
		".",
		"R",
		"O",
		"S",
		"A",
		"S",
		"O",
		"U",
		"P",
		"E",
		"D",
		".",
		"E",
		"B",
		"E",
		"N",
		"E",
		"Z",
		"E",
		"R",
		"T",
		"U",
		"C",
		"S",
		"O",
		"N",
		".",
		"L",
		"U",
		"K",
		"E",
		"W",
		"A",
		"R",
		"M",
		"A",
		"T",
		"H",
		"E",
		"N",
		"A",
		".",
		".",
		"D",
		"E",
		"T",
		"E",
		"C",
		"T",
		"S",
	},
	Gridnums: []int{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		0,
		0,
		8,
		9,
		10,
		11,
		12,
		13,
		14,
		0,
		0,
		0,
		0,
		0,
		0,
		15,
		0,
		16,
		0,
		0,
		0,
		0,
		0,
		17,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		18,
		0,
		0,
		0,
		0,
		0,
		19,
		0,
		0,
		0,
		0,
		20,
		0,
		0,
		21,
		0,
		22,
		0,
		0,
		0,
		0,
		23,
		0,
		0,
		0,
		24,
		0,
		25,
		0,
		0,
		26,
		0,
		27,
		0,
		0,
		0,
		28,
		0,
		0,
		0,
		29,
		30,
		0,
		0,
		0,
		0,
		31,
		0,
		32,
		0,
		0,
		0,
		0,
		0,
		33,
		0,
		0,
		0,
		0,
		0,
		34,
		0,
		35,
		0,
		0,
		0,
		0,
		36,
		37,
		0,
		0,
		0,
		0,
		0,
		38,
		0,
		0,
		0,
		0,
		0,
		0,
		39,
		0,
		0,
		0,
		0,
		0,
		0,
		40,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		41,
		0,
		0,
		0,
		42,
		0,
		43,
		0,
		0,
		0,
		0,
		0,
		44,
		45,
		46,
		47,
		0,
		0,
		48,
		0,
		49,
		0,
		0,
		0,
		0,
		50,
		51,
		0,
		0,
		0,
		52,
		0,
		0,
		0,
		53,
		0,
		54,
		0,
		0,
		55,
		0,
		56,
		0,
		0,
		0,
		57,
		0,
		0,
		0,
		0,
		58,
		0,
		59,
		0,
		0,
		60,
		0,
		0,
		0,
		0,
		61,
		0,
		0,
		0,
		0,
		0,
		0,
		62,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		63,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		64,
		0,
		0,
		0,
		0,
		0,
		0,
	},
	Publisher: "The New York Times",
	Size: struct {
		Cols int `json:"cols"`
		Rows int `json:"rows"`
	}{Cols: 15, Rows: 15},
	Title: "NY TIMES, FRI, MAR 09, 2018",
}

func (b BoardLayout) GetLastClue(number int, direction string) Clue {
	prevNum, prevDirection := b.Prev(number, direction)
	return b.GetClue(prevNum, prevDirection, nil)
}

func (b BoardLayout) GetClue(number int, direction string, expiresAt *time.Time) Clue {
	clues := b.Clues.Down
	answers := b.Answers.Down
	if direction == "across" {
		clues = b.Clues.Across
		answers = b.Answers.Across
	}
	return Clue{
		Number:      number,
		Direction:   direction,
		Description: clues[number],
		Answer:      answers[number],
		ExpiresAt:   expiresAt,
	}
}

func (b BoardLayout) Next(number int, direction string) (int, string) {
	if direction == "across" {
		if len(b.Clues.Across) == number+1 {
			return 0, "down"
		}
	} else {
		if len(b.Clues.Down) == number+1 {
			return 0, "across"
		}
	}
	return number + 1, direction
}

func (b BoardLayout) Prev(number int, direction string) (int, string) {
	if number == 0 {
		if direction == "across" {
			return len(b.Clues.Down) - 1, "down"
		} else {
			return len(b.Clues.Across) - 1, "across"
		}
	}
	return number - 1, direction
}
