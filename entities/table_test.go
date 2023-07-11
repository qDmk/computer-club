package entities

import (
	"computerClub/utils"
	"math/rand"
	"testing"
	"time"
)

func TestTable_ClientHoursStat(t *testing.T) {
	table := Table{}

	N := 1_000_000
	for i := 0; i < N; i++ {
		day := 24 * int64(time.Hour)
		// pick random time in [00:00, 24:00)
		from := utils.Midnight.Add(time.Duration(rand.Int63n(day))).Truncate(time.Minute)

		// pick random minutes amount in [1, 60)
		minutesTaken := rand.Int31n(59) + 1
		to := from.Add(time.Duration(minutesTaken) * time.Minute)

		table.Take(from)
		table.Leave(to)

		hours := int(table.ClientHoursStat().Hours())
		expected := i + 1
		if hours != expected {
			t.Fatalf(`Expected %v got %v`, expected, hours)
		}
	}
}
