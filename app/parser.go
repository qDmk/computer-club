package app

import (
	"bufio"
	"computerClub/club"
	"computerClub/entities"
	"computerClub/utils"
	"errors"
	"strconv"
	"strings"
)

func ParseInput(s *bufio.Scanner) (c *club.Club, es []entities.Event, parseErr error) {
	nextLine := func(errMsg string) bool {
		ok := s.Scan()
		if !ok {
			parseErr = errors.New(errMsg)
			return false
		}

		return true
	}

	if !nextLine("table amount not found") {
		return
	}
	tablesAmount, err := strconv.Atoi(s.Text())
	if err != nil {
		parseErr = errors.New(s.Text())
		return
	}

	if !nextLine("open hours not found") {
		return
	}
	schedule := strings.Split(s.Text(), " ")
	if len(schedule) != 2 {
		parseErr = errors.New(s.Text())
		return
	}
	openingTime, err := utils.ParseTime(schedule[0])
	if err != nil {
		parseErr = errors.New(s.Text())
		return
	}
	closingTime, err := utils.ParseTime(schedule[1])
	if err != nil || closingTime.Before(openingTime) {
		parseErr = errors.New(s.Text())
		return
	}

	if !nextLine("price not found") {
		return
	}
	price, err := strconv.Atoi(s.Text())
	if err != nil || price <= 0 {
		parseErr = errors.New(s.Text())
		return
	}

	c = club.NewClub(openingTime, closingTime, price, tablesAmount)

	prevTime := utils.Midnight
	for s.Scan() {
		eventStrings := strings.Split(s.Text(), " ")
		if len(eventStrings) < 3 {
			parseErr = errors.New(s.Text())
			return
		}

		eventTime, err := utils.ParseTime(eventStrings[0])
		if err != nil || eventTime.Before(prevTime) {
			parseErr = errors.New(s.Text())
			return
		}

		rawCode, err := strconv.Atoi(eventStrings[1])
		code := entities.EventCode(rawCode)
		if err != nil || !code.IsValidIncoming() {
			parseErr = errors.New(s.Text())
			return
		}

		client, err := entities.NewClientName(eventStrings[2])
		if err != nil {
			parseErr = errors.New(s.Text())
			return
		}

		var event entities.Event
		switch code {
		case entities.IncomingClientSat:
			if len(eventStrings) != 4 {
				parseErr = errors.New(s.Text())
				return
			}
			tableI, err := strconv.Atoi(eventStrings[3])
			if err != nil || 1 > tableI || tableI > tablesAmount {
				parseErr = errors.New(s.Text())
				return
			}
			tableI--

			event = entities.NewTableEvent(eventTime, code, client, tableI)
		default:
			if len(eventStrings) != 3 {
				parseErr = errors.New(s.Text())
				return
			}

			event = entities.NewClientEvent(eventTime, code, client)
		}

		es = append(es, event)

		prevTime = eventTime
	}

	return
}
