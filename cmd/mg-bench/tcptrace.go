package main

import (
	"net"
	"time"

	"sourcegraph.com/sourcegraph/appdash"
)

func init() {
	appdash.RegisterEvent(ConnEvent{})
	appdash.RegisterEvent(ConnReadEvent{})
	appdash.RegisterEvent(ConnWriteEvent{})
}

func NewConnEvent(c net.Conn) *ConnEvent {
	return &ConnEvent{Conn: connInfo(c)}
}

// RequestInfo describes an HTTP request.
type ConnInfo struct {
	RemoteAddr net.Addr
	LocalAddr  net.Addr
}

func connInfo(c net.Conn) ConnInfo {
	return ConnInfo{
		RemoteAddr: c.RemoteAddr(),
		LocalAddr:  c.LocalAddr(),
	}
}

// ConnEvent records an connection event.
type ConnEvent struct {
	Conn          ConnInfo
	ConnOpen      time.Time
	ConnConnected time.Time
	ConnClosed    time.Time
}

// Schema returns the constant "HTTPClient".
func (ConnEvent) Schema() string { return "ConnClient" }

// Important implements the appdash ImportantEvent.
func (ConnEvent) Important() []string {
	return []string{}
}

// Start implements the appdash TimespanEvent interface.
func (e ConnEvent) Start() time.Time { return e.ConnOpen }

// End implements the appdash TimespanEvent interface.
func (e ConnEvent) End() time.Time { return e.ConnClosed }

func NewConnReadEvent() *ConnReadEvent {
	return &ConnReadEvent{}
}

// ConnEvent records an connection event.
type ConnReadEvent struct {
	Count     uint64
	Error     string
	ReadStart time.Time
	ReadEnd   time.Time
}

// Schema returns the constant "HTTPClient".
func (ConnReadEvent) Schema() string { return "ConnRead" }

// Important implements the appdash ImportantEvent.
func (ConnReadEvent) Important() []string {
	return []string{}
}

// Start implements the appdash TimespanEvent interface.
func (e ConnReadEvent) Start() time.Time { return e.ReadStart }

// End implements the appdash TimespanEvent interface.
func (e ConnReadEvent) End() time.Time { return e.ReadEnd }

func NewConnWriteEvent() *ConnWriteEvent {
	return &ConnWriteEvent{}
}

// ConnEvent records an connection event.
type ConnWriteEvent struct {
	Count      uint64
	Error      string
	WriteStart time.Time
	WriteEnd   time.Time
}

// Schema returns the constant "HTTPClient".
func (ConnWriteEvent) Schema() string { return "ConnWrite" }

// Important implements the appdash ImportantEvent.
func (ConnWriteEvent) Important() []string {
	return []string{}
}

// Start implements the appdash TimespanEvent interface.
func (e ConnWriteEvent) Start() time.Time { return e.WriteStart }

// End implements the appdash TimespanEvent interface.
func (e ConnWriteEvent) End() time.Time { return e.WriteEnd }
