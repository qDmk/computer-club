package entities

import (
	"fmt"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format("15:04")
}

type EventCode int

const (
	IncomingClientCame EventCode = iota + 1
	IncomingClientSat
	IncomingClientAwait
	IncomingClientLeave
	OutgoingClientLeft = iota + 7
	OutgoingClientSat
	OutgoingError
)

func (code EventCode) IsValidIncoming() bool {
	return IncomingClientCame <= code && code <= IncomingClientLeave
}

// Event inspired by reflect's "panics if unsuitable"
type Event interface {
	Time() time.Time
	Code() EventCode
	Client() Client
	TableIdx() int
	ErrorMsg() string
	fmt.Stringer
}

type baseEvent struct {
	time time.Time
	code EventCode
}

func newBaseEvent(time time.Time, code EventCode) baseEvent {
	return baseEvent{time: time, code: code}
}

func (e baseEvent) Time() time.Time {
	return e.time
}

func (e baseEvent) Code() EventCode {
	return e.code
}

func (e baseEvent) Client() Client {
	panic("client undefined")
}

func (e baseEvent) TableIdx() int {
	panic("table index undefined")
}

func (e baseEvent) ErrorMsg() string {
	panic("error undefined")
}

func (e baseEvent) String() string {
	return fmt.Sprint(formatTime(e.Time()), " ", e.Code())
}

type errorEvent struct {
	baseEvent
	err string
}

func (e errorEvent) ErrorMsg() string {
	return e.err
}

func (e errorEvent) String() string {
	return fmt.Sprint(e.baseEvent.String(), " ", e.ErrorMsg())
}

type clientEvent struct {
	baseEvent
	client Client
}

func (e clientEvent) Client() Client {
	return e.client
}

func (e clientEvent) String() string {
	return fmt.Sprint(e.baseEvent.String(), " ", e.Client())
}

type tableEvent struct {
	clientEvent
	table int
}

func (e tableEvent) TableIdx() int {
	return e.table
}

func (e tableEvent) String() string {
	return fmt.Sprint(e.clientEvent.String(), " ", e.table+1)
}

func NewClientEvent(time time.Time, code EventCode, client Client) Event {
	return clientEvent{
		newBaseEvent(time, code),
		client,
	}
}

func NewTableEvent(time time.Time, code EventCode, client Client, tableI int) Event {
	return tableEvent{
		clientEvent{
			newBaseEvent(time, code),
			client,
		},
		tableI,
	}
}

func NewOutgoingClientLeft(e Event) Event {
	return NewClientEvent(e.Time(), OutgoingClientLeft, e.Client())
}

func NewOutgoingClientSat(e Event, client Client, tableI int) Event {
	return NewTableEvent(e.Time(), OutgoingClientSat, client, tableI)
}

func NewOutgoingError(e Event, msg string) Event {
	return errorEvent{
		newBaseEvent(e.Time(), OutgoingError),
		msg,
	}
}
