package club

import (
	"computerClub/entities"
	"errors"
	"time"
)

type ClientStatus int

const (
	InClub  ClientStatus = -2
	InQueue ClientStatus = -1
)

var (
	AlreadyInClubError   = errors.New("YouShallNotPass")
	ClientUnknownError   = errors.New("ClientUnknown")
	NotOpenYetError      = errors.New("NotOpenYet")
	TableIsBusyError     = errors.New("PlaceIsBusy")
	TablesAvailableError = errors.New("ICanWaitNoLonger!")
)

type Club struct {
	OpeningTime  time.Time
	ClosingTime  time.Time
	PricePerHour int
	clients      map[entities.Client]ClientStatus // -2 in club, -1 in queue, >=0 = tables idx
	// TODO: queue can be implemented with linked list for O(1) removals. If so, clients would store references to its nodes
	queue        []entities.Client
	maxQueueSize int
	tables       []entities.Table
	tablesUsed   int
}

func NewClub(openingTime, closingTime time.Time, price int, tablesAmount int) *Club {
	return &Club{
		OpeningTime:  openingTime,
		ClosingTime:  closingTime,
		PricePerHour: price,
		clients:      make(map[entities.Client]ClientStatus),
		queue:        nil,
		maxQueueSize: tablesAmount,
		tables:       make([]entities.Table, tablesAmount),
		tablesUsed:   0,
	}
}

func (c *Club) IsInClub(client entities.Client) bool {
	_, inClub := c.clients[client]
	return inClub
}

func (c *Club) TakeTable(time time.Time, client entities.Client, tableI int) {
	c.clients[client] = ClientStatus(tableI)
	c.tables[tableI].Take(time)
	c.tablesUsed++
}

func (c *Club) LeaveTable(time time.Time, client entities.Client, tableI int) {
	c.clients[client] = InClub
	c.tables[tableI].Leave(time)
	c.tablesUsed--
}

func (c *Club) ClientCame(time time.Time, name entities.Client) error {
	if c.IsInClub(name) {
		return AlreadyInClubError
	}

	if time.Before(c.OpeningTime) || time.After(c.ClosingTime) {
		return NotOpenYetError
	}

	c.clients[name] = InClub

	return nil
}

func (c *Club) ClientSat(time time.Time, client entities.Client, tableI int) error {
	if !c.IsInClub(client) {
		return ClientUnknownError
	}

	if c.tables[tableI].IsBusy() {
		return TableIsBusyError
	}

	clientStatus := c.clients[client]
	if clientStatus >= 0 {
		prevTable := int(clientStatus)
		c.LeaveTable(time, client, prevTable)
	}

	c.TakeTable(time, client, tableI)

	return nil
}

func (c *Club) ClientAwait(name entities.Client) (bool, error) {
	if !c.IsInClub(name) {
		return false, ClientUnknownError
	}

	if c.tablesUsed < len(c.tables) {
		return false, TablesAvailableError
	}

	if len(c.queue) == c.maxQueueSize {
		return true, nil
	}

	c.queue = append(c.queue, name)
	c.clients[name] = InQueue

	return false, nil
}

func (c *Club) ClientLeave(time time.Time, client entities.Client) (ClientStatus, error) {
	if !c.IsInClub(client) {
		return 0, ClientUnknownError
	}

	clientStatus := c.clients[client]
	if clientStatus == InQueue {
		for i, q := range c.queue {
			if q == client {
				c.queue = append(c.queue[:i], c.queue[i+1:]...)
				break
			}
		}
	} else if clientStatus >= 0 {
		tableI := int(clientStatus)
		c.LeaveTable(time, client, tableI)
	}

	delete(c.clients, client)

	return clientStatus, nil
}

func (c *Club) CheckQueue() (entities.Client, bool) {
	if len(c.queue) > 0 {
		client := c.queue[0]
		c.queue = c.queue[1:]

		return client, true
	}

	return "", false
}

func (c *Club) Close() []entities.Client {
	clients := make([]entities.Client, len(c.clients))

	i := 0
	for client, status := range c.clients {
		clients[i] = client

		if status >= 0 {
			tableI := int(status)
			c.LeaveTable(c.ClosingTime, client, tableI)
		}
	}

	c.queue = nil
	c.clients = make(map[entities.Client]ClientStatus)

	return clients
}

type TableStat struct {
	Revenue   int
	TimeTaken time.Duration
}

func (c *Club) ResetTables() []TableStat {
	stats := make([]TableStat, len(c.tables))
	for i, t := range c.tables {
		clientHours := int(t.ClientHoursStat().Hours())
		stats[i].Revenue = clientHours * c.PricePerHour
		stats[i].TimeTaken = t.TakenStat()

		// reset table stats
		c.tables[i] = entities.Table{}
	}

	return stats
}
