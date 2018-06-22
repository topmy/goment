package goment

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Thanks to https://github.com/go-shadow/moment/blob/master/moment_parser.go for help on formatReplacements and regex.
var formatReplacements = map[string]string{
	"M":    "1",
	"Mo":   "1<stdOrdinal>", // stdNumMonth 1st 2nd ... 11th 12th
	"MM":   "01",
	"MMM":  "Jan",
	"MMMM": "January",
	"D":    "2",
	"Do":   "2<stdOrdinal>",
	"DD":   "02",
	"DDD":  "<stdDayOfYear>",
	"DDDo": "<stdDayOfYear><stdOrdinal>",
	"DDDD": "<stdDayOfYearZero>",
	"d":    "<stdDayOfWeek>",
	"do":   "<stdDayOfWeek><stdOrdinal>",
	"dd":   "<stdShortDay>",
	"ddd":  "Mon",
	"dddd": "Monday",
	"e":    "<stdDayOfWeek>",
	"E":    "<stdDayOfWeekISO>",
	"w":    "<stdWeekOfYear>",
	"wo":   "<stdWeekOfYear><stdOrdinal>",
	"ww":   "<stdWeekOfYearZero>",
	"W":    "<stdIsoWeekOfYear>",
	"Wo":   "<stdIsoWeekOfYear><stdOrdinal>",
	"WW":   "<stdIsoWeekOfYearZero>",
	"YY":   "06",
	"YYYY": "2006",
	"Q":    "<stdQuarter>",
	"A":    "PM",
	"a":    "pm",
	"H":    "<stdHourNoZero>",
	"HH":   "15",
	"h":    "3",
	"hh":   "03",
	"m":    "4",
	"mm":   "04",
	"s":    "5",
	"ss":   "05",
	"z":    "MST",
	"zz":   "MST",
	"Z":    "Z07:00",
	"ZZ":   "-0700",
	"X":    "<stdUnix>",
	"LT":   "3:04 PM",
	"LTS":  "3:04:05 PM",
	"L":    "01/02/2006",
	"LL":   "January 2, 2006",
	"l":    "1/2/2006",
	"ll":   "Jan 2, 2006",
	"LLL":  "January 2, 2006 3:04 PM",
	"lll":  "Jan 2, 2006 3:04 PM",
	"LLLL": "Monday, January 2, 2006 3:04 PM",
	"llll": "Mon, Jan 2, 2006 3:04 PM",
}

// Format takes a string of tokens and replaces them with their corresponding values to display the Goment.
func (g *Goment) Format(args ...interface{}) string {
	layout := ""

	numArgs := len(args)
	if numArgs < 1 {
		layout = "YYYY-MM-DDTHH:mm:ssZ"
	} else {
		layout = args[0].(string)
	}

	format := convert(layout)
	formatted := g.ToTime().Format(format)

	return performReplacements(g, formatted)
}

func performReplacements(g *Goment, formatted string) string {
	if strings.Contains(formatted, "<std") {
		formatted = strings.Replace(formatted, "<stdDayOfYear>", fmt.Sprintf("%d", g.DayOfYear()), -1)
		formatted = strings.Replace(formatted, "<stdDayOfYearZero>", g.dayOfYearZero(), -1)
		formatted = strings.Replace(formatted, "<stdDayOfWeek>", fmt.Sprintf("%d", g.Day()), -1)
		formatted = strings.Replace(formatted, "<stdDayOfWeekISO>", fmt.Sprintf("%d", g.ISOWeekday()), -1)
		formatted = strings.Replace(formatted, "<stdUnix>", fmt.Sprintf("%d", g.ToUnix()), -1)
		formatted = strings.Replace(formatted, "<stdQuarter>", fmt.Sprintf("%d", g.Quarter()), -1)
		formatted = strings.Replace(formatted, "<stdIsoWeekOfYear>", fmt.Sprintf("%d", g.ISOWeek()), -1)
		formatted = strings.Replace(formatted, "<stdIsoWeekOfYearZero>", g.isoWeekOfYearZero(), -1)
		formatted = strings.Replace(formatted, "<stdWeekOfYear>", fmt.Sprintf("%d", g.Week()), -1)
		formatted = strings.Replace(formatted, "<stdWeekOfYearZero>", g.weekOfYearZero(), -1)
		formatted = strings.Replace(formatted, "<stdHourNoZero>", fmt.Sprintf("%d", g.Hour()), -1)
		formatted = strings.Replace(formatted, "<stdShortDay>", fmt.Sprintf("%v", g.ToTime().Weekday().String()[0:2]), -1)

		if strings.Contains(formatted, "<stdOrdinal>") {
			regex := regexp.MustCompile("([0-9]+)(?:<stdOrdinal>)")

			formatted = regex.ReplaceAllStringFunc(formatted, func(n string) string {
				o, _ := strconv.Atoi(strings.Replace(n, "<stdOrdinal>", "", 1))
				return ordinal(o)
			})
		}
	}

	return formatted
}

func (g *Goment) dayOfYearZero() string {
	day := g.ToTime().YearDay()

	if day < 10 {
		return fmt.Sprintf("00%d", day)
	}

	if day < 100 {
		return fmt.Sprintf("0%d", day)
	}

	return fmt.Sprintf("%d", day)
}

func (g *Goment) isoWeekOfYearZero() string {
	week := g.ISOWeek()

	if week < 10 {
		return fmt.Sprintf("0%d", week)
	}
	return fmt.Sprintf("%d", week)
}

func (g *Goment) weekOfYearZero() string {
	week := g.Week()

	if week < 10 {
		return fmt.Sprintf("0%d", week)
	}
	return fmt.Sprintf("%d", week)
}

func ordinal(x int) string {
	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return strconv.Itoa(x) + suffix
}

func convert(layout string) string {
	reBrackets := regexp.MustCompile(`\[([^\[\]]*)\]`)
	reFormats := regexp.MustCompile("(LT[S]?|LL?L?L?|l{1,4}|Mo|MM?M?M?|Do|DDDo|DD?D?D?|ddd?d?|do?|w[o|w]?|W[o|W]?|YYYYY|YYYY|YY|gg(ggg?)?|GG(GGG?)?|e|E|a|A|hh?|HH?|mm?|ss?|SS?S?|X|zz?|ZZ?|Q)")

	bracketMatch := reBrackets.FindAllStringSubmatch(layout, -1)
	bracketsFound := len(bracketMatch) > 0

	if bracketsFound {
		for i := range bracketMatch {
			layout = strings.Replace(layout, bracketMatch[i][0], makeToken(i+1), 1)
		}
	}

	var match [][]int
	if match = reFormats.FindAllStringSubmatchIndex(layout, -1); match == nil {
		return layout
	}

	for i := range match {
		start, end := match[i][0], match[i][1]
		matchText := layout[start:end]

		if replaceText, ok := formatReplacements[matchText]; ok {
			diff := len(replaceText) - len(matchText)
			layout = layout[0:start] + replaceText + layout[end:len(layout)]

			// If the replacement text is longer/shorter than the match, shift the remaining indexes.
			if diff != 0 {
				for j := i + 1; j < len(match); j++ {
					match[j][0] += diff
					match[j][1] += diff
				}
			}
		}
	}

	if bracketsFound {
		for i := range bracketMatch {
			layout = strings.Replace(layout, makeToken(i+1), bracketMatch[i][1], 1)
		}
	}

	return layout
}

func makeToken(num int) string {
	return fmt.Sprintf("$%v", num)
}
