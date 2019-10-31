package hakone

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Runner string
type Grade string
type Team string
type Note string

type Time int

type Record struct {
	Order             int    `json:"order"`
	Runner            Runner `json:"runner"`
	Grade             Grade  `json:"grade"`
	Team              Team   `json:"team"`
	TimeOf5km         Time   `json:"time_of_5_km"`
	TimeOf10km        Time   `json:"time_of_10_km"`
	TimeOf15km        Time   `json:"time_of_15_km"`
	TimeOf20km        Time   `json:"time_of_20_km"`
	FinishTime        Time   `json:"finish_time"`
	RapFrom5kmTo10km  Time   `json:"rap_5_to_10"`
	RapFrom10kmTo15km Time   `json:"rap_10_to_15"`
	RapFrom15kmTo20km Time   `json:"rap_15_to_20"`
	Note              Note
}

func (t Time) plus(d int) Time {
	result := int(t) + d
	return Time(result)
}

func NewTime(t string) (Time, error) {
	count := strings.Count(t, ":")
	if count == 0 || count > 2 {
		return 0, errors.New(fmt.Sprintf("invalid time char sequence: %s", t))
	} else if count == 1 {
		return Mins(t)
	} else { // count == 2
		sep := strings.Split(t, ":")
		hour, err := strconv.Atoi(sep[0])
		if err != nil {
			return 0, errors.New(fmt.Sprintf("invalid number at hour part: %s", t))
		}
		minTime, err := SplitString(t, sep[1], sep[2])
		if err != nil {
			return 0, err
		}
		return minTime.plus(hour * 60 * 60), nil
	}
}

func Mins(t string) (Time, error) {
	count := strings.Count(t, ":")
	if count != 1 {
		return 0, errors.New(fmt.Sprintf("invalid time char sequence: %s", t))
	}
	sep := strings.Split(t, ":")
	return SplitString(t, sep[0], sep[1])
}

func SplitString(original, m, s string) (Time, error) {
	min, err := strconv.Atoi(m)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("invalid number at minute part: %s", original))
	}
	sec, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("invalid number at second part: %s", original))
	}
	return Time(min*60 + sec), nil
}
