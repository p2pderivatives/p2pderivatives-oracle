package iso8601

import (
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Time constant as time.Duration
const (
	Second = time.Second
	Minute = Second * 60
	Hour   = Minute * 60
	Day    = Hour * 24
	Month  = Day * 30
	Year   = Day * 365
)

// DurationRegex regular expression for duration iso 8601
var DurationRegex = regexp.MustCompile(`P([\d\.]+Y)?([\d\.]+M)?([\d\.]+D)?T([\d\.]+H)?([\d\.]+M)?([\d\.]+?S)?`)

// ParseDuration converts a ISO8601 duration into a time.Duration
func ParseDuration(str string) (time.Duration, error) {
	if !DurationRegex.MatchString(str) {
		return time.Duration(0), errors.New("duration string not matching ISO8601 format")
	}
	matches := DurationRegex.FindStringSubmatch(str)

	years := parseDurationPart(matches[1], Year)
	months := parseDurationPart(matches[2], Month)
	days := parseDurationPart(matches[3], Day)
	hours := parseDurationPart(matches[4], Hour)
	minutes := parseDurationPart(matches[5], Minute)
	seconds := parseDurationPart(matches[6], Second)

	return years + months + days + hours + minutes + seconds, nil
}

func parseDurationPart(value string, unit time.Duration) time.Duration {
	if len(value) != 0 {
		if parsed, err := strconv.ParseFloat(value[:len(value)-1], 64); err == nil {
			return time.Duration(float64(unit) * parsed)
		}
	}
	return 0
}

var formatArr = []struct {
	suffix string
	dur    time.Duration
}{
	{suffix: "Y", dur: Year},
	{suffix: "M", dur: Month},
	{suffix: "D", dur: Day},
	{suffix: "H", dur: Hour},
	{suffix: "M", dur: Minute},
	{suffix: "S", dur: Second},
}

// EncodeDuration converts a time.Duration to an ISO8601 duration string
// examples :
//	3h -> PT3H
//	7h22min -> PT7h22M
func EncodeDuration(d time.Duration) string {
	durSum := d.Nanoseconds()
	if durSum == 0 {
		return "PT0S"
	}
	res := "P"
	for i, format := range formatArr {
		if i == 3 {
			res += "T"
		}

		unit := durSum / format.dur.Nanoseconds()
		if unit != 0 {
			durSum -= unit * format.dur.Nanoseconds()
			res += strconv.FormatInt(unit, 10) + format.suffix
		}
	}

	return res
}
