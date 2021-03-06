package openinghours

import (
	"fmt"
	"strings"
	"time"
)

/*
Unser Ziel:
Mon=07:00-20:00;Tue=07:00-20:00;Wen=07:00-20:00;Thu=07:00-20:00;Fri=07:00-20:00;Sat=07:00-20:00;Sun=x
*/

type Weekday string

type OpeningHours map[Weekday]*OpeningDay

type ExtraOpeningHours []ExtraOpening

type ExtraOpening struct {
	Start   time.Time   `bson:"start"`
	End     time.Time   `bson:"end"`
	Opening *OpeningDay `bson:"opening"`
}

type OpeningDay struct {
	Opening string `bson:"opening"`
	Closing string `bson:"closing"`
}

const (
	Monday    Weekday = "Mon"
	Tuesday   Weekday = "Tue"
	Wednesday Weekday = "Wen"
	Thursday  Weekday = "Thu"
	Friday    Weekday = "Fri"
	Saturday  Weekday = "Sat"
	Sunday    Weekday = "Sun"
)

var possibleWeekdays = map[Weekday]bool{
	Monday:    true,
	Tuesday:   true,
	Wednesday: true,
	Thursday:  true,
	Friday:    true,
	Saturday:  true,
	Sunday:    true,
}

func ParseOpeningHours(openingHoursString string) (OpeningHours, ExtraOpeningHours, error) {
	parts := strings.Split(openingHoursString, ";")

	result := OpeningHours{}
	extra := ExtraOpeningHours{}
	for _, singleDay := range parts {
		singleDay = strings.ReplaceAll(singleDay, " ", "=")
		dayKv := strings.Split(singleDay, "=")

		if len(dayKv) != 2 {
			return nil, nil, fmt.Errorf("expected day statement to contain two parts, got %v", dayKv)
		}

		if weekday, ok := tryParseWeekday(dayKv[0]); ok {
			timeRange, err := tryParseTimeRange(dayKv[1])
			if err != nil {
				return nil, nil, err
			}
			result[weekday] = timeRange
		} else {
			extraOpening, err := tryParseExtraOpening(dayKv[0], dayKv[1])
			if err != nil {
				return nil, nil, err
			}
			extra = append(extra, extraOpening)
		}
	}
	if len(result) != 7 {
		return nil, nil, fmt.Errorf("incomplete opening hours string: %s", openingHoursString)
	}

	return result, extra, nil
}

func tryParseWeekday(stringToTest string) (Weekday, bool) {
	if _, ok := possibleWeekdays[Weekday(stringToTest)]; ok {
		return Weekday(stringToTest), true
	}
	return "", false
}

func tryParseTimeRange(timeRange string) (*OpeningDay, error) {
	if timeRange == "x" {
		return nil, nil
	}

	splitTimeRange := strings.Split(timeRange, "-")
	if len(splitTimeRange) != 2 {
		return nil, fmt.Errorf("expected two parts, got %d, input: %v", len(splitTimeRange), timeRange)
	}

	return &OpeningDay{
		Opening: splitTimeRange[0],
		Closing: splitTimeRange[1],
	}, nil
}

func tryParseExtraOpening(dateRange, timeRange string) (ExtraOpening, error) {
	var dateParts []string

	if strings.Contains(dateRange, "-") {
		dateParts = strings.Split(dateRange, "-")
	} else {
		dateParts = []string{dateRange, dateRange}
	}

	if len(dateParts) != 2 {
		return ExtraOpening{}, fmt.Errorf("error parsing date range, expected one or two parts, got %d, value: %s", len(dateParts), dateRange)
	}

	start, err := time.Parse("02.01.2006", dateParts[0])
	if err != nil {
		return ExtraOpening{}, err
	}

	end, err := time.Parse("02.01.2006", dateParts[1])
	if err != nil {
		return ExtraOpening{}, err
	}

	openingDay, err := tryParseTimeRange(timeRange)
	if err != nil {
		return ExtraOpening{}, err
	}

	return ExtraOpening{
		Start:   start,
		End:     end,
		Opening: openingDay,
	}, nil
}
