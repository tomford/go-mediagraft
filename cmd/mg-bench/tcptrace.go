package main

import (
	"net"
	"time"

	"sync"

	"sourcegraph.com/sourcegraph/appdash"
)

var spanMap *connSpanMap

func init() {
	spanMap = &connSpanMap{
		lock: &sync.RWMutex{},
		smap: make(map[net.Conn]*appdash.Recorder),
	}

	appdash.RegisterEvent(ConnEvent{})
	appdash.RegisterEvent(ConnReadEvent{})
	appdash.RegisterEvent(ConnWriteEvent{})
}

func NewConnEvent(c net.Conn) *ConnEvent {
	return &ConnEvent{Connection: connInfo(c)}
}

// RequestInfo describes an HTTP request.
type ConnInfo struct {
	RemoteAddr string
	LocalAddr  string
}

func connInfo(c net.Conn) ConnInfo {
	return ConnInfo{
		RemoteAddr: c.RemoteAddr().String(),
		LocalAddr:  c.LocalAddr().String(),
	}
}

// ConnEvent records an connection event.
type ConnEvent struct {
	Connection ConnInfo
	Opened     time.Time
	Connected  time.Time
}

// Schema returns the constant "HTTPClient".
func (ConnEvent) Schema() string { return "ConnOpen" }

// Important implements the appdash ImportantEvent.
func (ConnEvent) Important() []string {
	return []string{"Opened", "Connected"}
}

// Start implements the appdash TimespanEvent interface.
func (e ConnEvent) Start() time.Time { return e.Opened }

// End implements the appdash TimespanEvent interface.
func (e ConnEvent) End() time.Time { return e.Connected }

func NewConnReadEvent() *ConnReadEvent {
	return &ConnReadEvent{}
}

// ConnEvent records an connection event.
type ConnReadEvent struct {
	ReadCount int
	Error     string
	ReadStart time.Time
	ReadEnd   time.Time
}

// Schema returns the constant "HTTPClient".
func (ConnReadEvent) Schema() string { return "ConnRead" }

// Important implements the appdash ImportantEvent.
func (ConnReadEvent) Important() []string {
	return []string{"ReadStart", "ReadEnd", "ReadCount"}
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
	WriteCount int
	Error      string
	WriteStart time.Time
	WriteEnd   time.Time
}

// Schema returns the constant "HTTPClient".
func (ConnWriteEvent) Schema() string { return "ConnWrite" }

// Important implements the appdash ImportantEvent.
func (ConnWriteEvent) Important() []string {
	return []string{"WriteStart", "WriteEnd", "WriteCount"}
}

// Start implements the appdash TimespanEvent interface.
func (e ConnWriteEvent) Start() time.Time { return e.WriteStart }

// End implements the appdash TimespanEvent interface.
func (e ConnWriteEvent) End() time.Time { return e.WriteEnd }

// a net.Conn implementation that tracks connections and send/recv
type traceConn struct {
	base net.Conn
	rec  *appdash.Recorder
}

func (c traceConn) Read(b []byte) (n int, err error) {
	ev := NewConnReadEvent()
	ev.ReadStart = time.Now()
	n, err = c.base.Read(b)
	ev.ReadEnd = time.Now()
	ev.ReadCount = n
	if err != nil {
		ev.Error = err.Error()
	}

	c.rec.Event(ev)

	return n, err
}

func (c traceConn) Write(b []byte) (n int, err error) {
	ev := NewConnWriteEvent()
	ev.WriteStart = time.Now()
	n, err = c.base.Write(b)
	ev.WriteEnd = time.Now()
	ev.WriteCount = n
	if err != nil {
		ev.Error = err.Error()
	}

	c.rec.Event(ev)

	return n, err
}

func (c traceConn) Close() error {
	return c.base.Close()
}

func (c traceConn) LocalAddr() net.Addr {
	return c.base.LocalAddr()
}

func (c traceConn) RemoteAddr() net.Addr {
	return c.base.RemoteAddr()
}

func (c traceConn) SetDeadline(t time.Time) error {
	return c.base.SetDeadline(t)
}

func (c traceConn) SetReadDeadline(t time.Time) error {
	return c.base.SetReadDeadline(t)
}

func (c traceConn) SetWriteDeadline(t time.Time) error {
	return c.base.SetWriteDeadline(t)
}

// Thread safe map for tracking connections and spans
type connSpanMap struct {
	lock *sync.RWMutex
	smap map[net.Conn]*appdash.Recorder
}

func (m *connSpanMap) Get(c net.Conn) (*appdash.Recorder, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r, ok := m.smap[c]
	return r, ok
}

func (m *connSpanMap) Set(c net.Conn, r *appdash.Recorder) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.smap[c] = r
	return
}

func MakeTraceDialer(r *appdash.Recorder, defaultDial func(network string, address string) (net.Conn, error)) func(network string, address string) (net.Conn, error) {
	return func(network string, address string) (net.Conn, error) {
		begin := time.Now()
		conn, err := defaultDial(network, address)
		conned := time.Now()

		cr, ok := spanMap.Get(conn)
		if !ok {
			cr = r.Child()
			cr.Name("net.Conn")
			spanMap.Set(conn, cr)

			ce := NewConnEvent(conn)
			ce.Opened = begin
			ce.Connected = conned
			cr.Event(ce)
		}

		// TODO(tcm) This span should really be rooted off of the connection
		tconn := traceConn{
			base: conn,
			rec:  cr,
		}
		return tconn, err
	}
}
