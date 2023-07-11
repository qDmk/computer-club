package entities

import (
	"math"
	"time"
)

type Table struct {
	clientRegisterTime time.Time
	stat               time.Duration
	clientHours        time.Duration
}

func (t *Table) Take(time time.Time) {
	if !t.clientRegisterTime.IsZero() {
		panic("table is busy")
	}

	t.clientRegisterTime = time
}

func (t *Table) Leave(leaveTime time.Time) {
	if t.clientRegisterTime.IsZero() {
		panic("table is not busy")
	}

	d := leaveTime.Sub(t.clientRegisterTime)
	t.stat += d
	t.clientHours += time.Hour * time.Duration(math.Ceil(d.Hours()))

	t.clientRegisterTime = time.Time{}

	return
}

func (t *Table) IsBusy() bool {
	return !t.clientRegisterTime.IsZero()
}

func (t *Table) TakenStat() time.Duration {
	return t.stat
}

func (t *Table) ClientHoursStat() time.Duration {
	return t.clientHours
}
