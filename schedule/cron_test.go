package schedule

import (
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	now, _ := time.Parse("2006-01-02", "2015-01-01")
	testCases := []struct {
		cronLine string
		expected time.Time
	}{
		{
			cronLine: "30 * * * *",
			expected: now.Add(30 * time.Minute),
		},
		{
			cronLine: "* 2 * * *",
			expected: now.Add(2 * time.Hour),
		},
	}

	for _, tc := range testCases {
		c := Cron(tc.cronLine)
		actual := c.Next(now)
		if actual != tc.expected {
			t.Errorf("Expected: %s, but got: %s", tc.expected, actual)
		}
	}
}
