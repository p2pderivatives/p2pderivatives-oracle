package iso8601_test

import (
	"p2pderivatives-oracle/internal/utils/iso8601"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TestDurationVectors = []struct {
	iso8601 string
	dur     time.Duration
}{
	{iso8601: "PT1H30M", dur: iso8601.Hour + 30*iso8601.Minute},
	{iso8601: "PT20M50S", dur: 20*iso8601.Minute + 50*iso8601.Second},
	{iso8601: "P1DT4H", dur: iso8601.Day + 4*iso8601.Hour},
	{iso8601: "P3MT22M", dur: 3*iso8601.Month + 22*iso8601.Minute},
	{iso8601: "P1YT3M", dur: iso8601.Year + 3*iso8601.Minute},
}

var TestInvalidIso8601 = []string{
	"4378T891",
	"P2687642",
	"bdakhdo8478.24",
}

func TestParseDuration_ValidISO8601_ReturnsCorrectValue(t *testing.T) {
	for _, v := range TestDurationVectors {
		actual, err := iso8601.ParseDuration(v.iso8601)
		assert.NoError(t, err)
		assert.Equal(t, v.dur, actual)
	}
}

func TestParseDuration_InvalidISO8601_ReturnsError(t *testing.T) {
	for _, v := range TestInvalidIso8601 {
		_, err := iso8601.ParseDuration(v)
		assert.Error(t, err, "value: %v", v)
	}
}

func TestEncodeDuration_ReturnsCorrectValue(t *testing.T) {
	for _, v := range TestDurationVectors {
		actual := iso8601.EncodeDuration(v.dur)
		assert.Equal(t, v.iso8601, actual, "value %v", v.dur)
	}
}
