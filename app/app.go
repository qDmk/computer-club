package app

import (
	"computerClub/club"
	"computerClub/entities"
	"computerClub/utils"
	"log"
	"sort"
)

type App struct {
	Club *club.Club
	Log  *log.Logger
}

func (a *App) processEvent(event entities.Event) {
	a.Log.Println(event)
	c := a.Club

	outputError := func(msg string) {
		a.Log.Println(entities.NewOutgoingError(event, msg))
	}

	switch event.Code() {
	case entities.IncomingClientCame:
		err := c.ClientCame(event.Time(), event.Client())
		if err != nil {
			outputError(err.Error())
			return
		}

	case entities.IncomingClientSat:
		err := c.ClientSat(event.Time(), event.Client(), event.TableIdx())
		if err != nil {
			outputError(err.Error())
			return
		}

	case entities.IncomingClientAwait:
		clientLeft, err := c.ClientAwait(event.Client())
		if err != nil {
			outputError(err.Error())
			return
		}

		if clientLeft {
			a.Log.Println(entities.NewOutgoingClientLeft(event))
		}

	case entities.IncomingClientLeave:
		status, err := c.ClientLeave(event.Time(), event.Client())
		if err != nil {
			outputError(err.Error())
			return
		}

		if status >= 0 {
			tableI := int(status)
			newClient, ok := c.CheckQueue()
			if ok {
				c.TakeTable(event.Time(), newClient, tableI)
				a.Log.Println(entities.NewOutgoingClientSat(event, newClient, tableI))
			}
		}
	}
}

func (a *App) Run(es []entities.Event) {
	a.Log.Println(utils.FormatTime(a.Club.OpeningTime))

	for _, e := range es {
		a.processEvent(e)
	}

	clients := a.Club.Close()
	sort.Slice(clients, func(i, j int) bool { return clients[i] < clients[j] })
	for _, client := range clients {
		a.Log.Println(entities.NewClientEvent(a.Club.ClosingTime, entities.OutgoingClientLeft, client))
	}

	a.Log.Println(utils.FormatTime(a.Club.ClosingTime))

	for i, s := range a.Club.ResetTables() {
		a.Log.Println(i+1, s.Revenue, utils.FormatDuration(s.TimeTaken))
	}
}
