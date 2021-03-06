package openinghours

import (
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseOpeningHours_ValidStringSundayClosedNoExtraOpenings(t *testing.T) {
	// given
	testOpeningHours := "Mon=07:00-10:00;Tue=07:00-11:00;Wen=07:00-12:00;Thu=07:00-13:00;Fri=07:00-14:00;Sat=07:00-15:00;Sun=x"

	// when
	result, extra, err := ParseOpeningHours(testOpeningHours)

	// then
	require.NoError(t, err)
	require.Nil(t, deep.Equal(result, OpeningHours{
		Monday: {
			Opening: "07:00",
			Closing: "10:00",
		},
		Tuesday: {
			Opening: "07:00",
			Closing: "11:00",
		},
		Wednesday: {
			Opening: "07:00",
			Closing: "12:00",
		},
		Thursday: {
			Opening: "07:00",
			Closing: "13:00",
		},
		Friday: {
			Opening: "07:00",
			Closing: "14:00",
		},
		Saturday: {
			Opening: "07:00",
			Closing: "15:00",
		},
		Sunday: nil,
	}))
	require.Empty(t, extra)
}

func TestParseOpeningHours_ValidStringExtraOpeningHours(t *testing.T) {
	testOpeningHours := "Mon=07:00-21:00;Tue=07:00-21:00;Wen=07:00-21:00;Thu=07:00-21:00;Fri=07:00-21:00;Sat=07:00-21:00;Sun=x;20.02.2021 07:00-14:00;21.02.2021-28.02.2021 x;01.03.2021 07:00-21:00"

	normal, extra, err := ParseOpeningHours(testOpeningHours)

	require.NoError(t, err)
	require.NotNil(t, normal)

	require.Nil(t, deep.Equal(extra, ExtraOpeningHours{
		{
			Start: time.Date(2021, 02, 20, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 02, 20, 0, 0, 0, 0, time.UTC),
			Opening: &OpeningDay{
				Opening: "07:00",
				Closing: "14:00",
			},
		}, {
			Start:   time.Date(2021, 2, 21, 0, 0, 0, 0, time.UTC),
			End:     time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			Opening: nil,
		}, {
			Start: time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
			Opening: &OpeningDay{
				Opening: "07:00",
				Closing: "21:00",
			},
		},
	}))
}

func TestParseOpeningHours_InvalidStringInvalidWeekday(t *testing.T) {
	// given
	testOpeningHours := "Invalid=07:00-10:00;Tue=07:00-11:00;Wen=07:00-12:00;Thu=07:00-13:00;Fri=07:00-14:00;Sat=07:00-15:00;Sun=x"

	// when
	result, extra, err := ParseOpeningHours(testOpeningHours)

	require.Error(t, err)
	require.Nil(t, result)
	require.Nil(t, extra)
}

func TestParseOpeningHours_InvalidStringIncompleteWeek(t *testing.T) {
	// given
	testOpeningHours := "Mon=07:00-10:00;Tue=07:00-11:00;Wen=07:00-12:00;Thu=07:00-13:00;Fri=07:00-14:00;Sat=07:00-15:00"

	// when
	result, extra, err := ParseOpeningHours(testOpeningHours)

	require.Error(t, err)
	require.Nil(t, result)
	require.Nil(t, extra)
}
