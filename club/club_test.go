package club

import (
	"computerClub/entities"
	"computerClub/utils"
	"fmt"
	"runtime/debug"
	"strconv"
	"testing"
	"time"
)

func assertEqual(t *testing.T, expected, actual any) {
	if expected != actual {
		t.Log(string(debug.Stack()))
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Fatalf("no error expected, got: %v", err)
	}
}

func newClub() *Club {
	openingTime, _ := utils.ParseTime("09:00")
	closingTime, _ := utils.ParseTime("18:00")
	return NewClub(openingTime, closingTime, 10, 3)
}

func TestClub_ClientCame(t *testing.T) {
	club := newClub()

	clientCame := func(timeShift int, client entities.Client, expectedError error, expectIn bool) {
		cameTime := club.OpeningTime.Add(time.Duration(timeShift) * time.Minute)
		err := club.ClientCame(cameTime, client)
		assertEqual(t, expectedError, err)
		assertEqual(t, expectIn, club.IsInClub(client))
	}

	// too early
	clientCame(-10, "supersonic", NotOpenYetError, false)

	// ok
	clientCame(10, "dima", nil, true)

	// already in
	clientCame(11, "dima", AlreadyInClubError, true)

	// too late
	tooLateTime := int(club.ClosingTime.Sub(club.OpeningTime).Minutes()) + 10
	clientCame(tooLateTime, "superturtle", NotOpenYetError, false)
}

func TestClub_ClientAwait(t *testing.T) {
	club := newClub()

	checkAwait := func(client entities.Client, expectedError error, expectLeave bool) {
		clientLeft, err := club.ClientAwait(client)
		assertEqual(t, expectedError, err)
		if err == nil {
			assertEqual(t, expectLeave, clientLeft)
		}
	}

	expectCome := func(client entities.Client) {
		err := club.ClientCame(club.OpeningTime, client)
		assertNoError(t, err)
	}

	expectSit := func(client entities.Client, tableI int) {
		err := club.ClientSat(club.OpeningTime, client, tableI)
		assertNoError(t, err)
	}

	// dima is unknown yet
	checkAwait("dima", ClientUnknownError, false)

	// dima came
	expectCome("dima")

	// there are tables available
	checkAwait("dima", TablesAvailableError, false)

	// Fill tables
	for club.tablesUsed != len(club.tables) {
		client := entities.Client("client_table_" + strconv.Itoa(club.tablesUsed))

		expectCome(client)
		expectSit(client, club.tablesUsed)
	}

	// ok
	checkAwait("dima", nil, false)

	// Fill queue
	for len(club.queue) != club.maxQueueSize {
		inQueueClient := entities.Client("client_in_queue_" + strconv.Itoa(len(club.queue)))

		expectCome(inQueueClient)
		checkAwait(inQueueClient, nil, false)
	}

	// Queue is full
	expectCome("i_left_the_queue")
	checkAwait("i_left_the_queue", nil, true)
}

// This test found me 1 error
func TestClub_ClientSat(t *testing.T) {
	club := newClub()

	expectCome := func(client entities.Client) {
		err := club.ClientCame(club.OpeningTime, client)
		assertNoError(t, err)
	}

	checkSit := func(client entities.Client, tableI int, expectedError error) {
		err := club.ClientSat(club.OpeningTime, client, tableI)
		assertEqual(t, expectedError, err)
	}

	// dima have not come yet
	checkSit("dima", 0, ClientUnknownError)

	// ok sit
	expectCome("dima")
	checkSit("dima", 0, nil)

	// can not sit to the same table
	checkSit("dima", 0, TableIsBusyError)

	// ok change table
	checkSit("dima", 1, nil) // here was the error, "dima" did not leave previous table

	// ok sit
	expectCome("teo")
	checkSit("teo", 0, nil)

	// table is busy
	checkSit("teo", 1, TableIsBusyError)
}

func TestClub_ClientLeave(t *testing.T) {
	club := newClub()

	expectCome := func(client entities.Client) {
		err := club.ClientCame(club.OpeningTime, client)
		assertNoError(t, err)
	}

	expectAwait := func(client entities.Client) {
		_, err := club.ClientAwait(client)
		assertNoError(t, err)
	}

	expectSit := func(client entities.Client, tableI int) {
		err := club.ClientSat(club.OpeningTime, client, tableI)
		assertNoError(t, err)
	}

	checkLeave := func(client entities.Client, expectedError error, expectedStatus ClientStatus) {
		status, err := club.ClientLeave(club.OpeningTime, client)
		assertEqual(t, expectedError, err)
		if err == nil {
			assertEqual(t, expectedStatus, status)
		}
	}

	// dima have not come yet
	checkLeave("dima", ClientUnknownError, 0)

	// ok leave club
	expectCome("dima")
	checkLeave("dima", nil, InClub)

	// ok leave table
	expectCome("dima")
	expectSit("dima", 0)
	checkLeave("dima", nil, 0)

	// fill tables
	for i := 0; i < len(club.tables); i++ {
		clientName := fmt.Sprint("client_table_", i)
		client := entities.Client(clientName)

		expectCome(client)
		expectSit(client, i)
	}

	// ok leave queue
	expectCome("dima")
	expectAwait("dima")
	checkLeave("dima", nil, InQueue)
}
