package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"sync"

	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sourcegraph/appdash/httptrace"
	"sourcegraph.com/sourcegraph/appdash/traceapp"

	mg "github.com/we7/go-mediagraft/pkg/mediagraft"
	"github.com/we7/go-mediagraft/pkg/mediagraft/oauth"
)

// Execute execudes the commands with the given arguments and returns an error,
// if any.
func main() {
	spanMap := &connSpanMap{
		lock: &sync.RWMutex{},
		smap: make(map[net.Conn]appdash.SpanID),
	}

	// We create a new in-memory store. All information about traces will
	// eventually be stored here.
	store := appdash.NewMemoryStore()

	// Listen on any available TCP port locally.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		log.Fatal(err)
	}
	collectorPort := l.Addr().(*net.TCPAddr).Port
	log.Printf("Appdash collector listening on tcp:%d", collectorPort)

	// Start an Appdash collection server that will listen for spans and
	// annotations and add them to the local collector (stored in-memory).
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	cs.Debug = true
	cs.Trace = true
	go cs.Start()

	// Print the URL at which the web UI will be running.
	port := ":3000"
	appdashURLStr := "http://localhost" + port
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		log.Fatalf("Error parsing http://localhost:%s: %s", port, err)
	}
	log.Printf("Appdash web UI running at %s", appdashURL)

	// Start the web UI in a separate goroutine.
	tapp := traceapp.New(nil)
	tapp.Store = store
	tapp.Queryer = store
	go func() {
		log.Fatal(http.ListenAndServe(port, tapp))
	}()

	localCollector := appdash.NewRemoteCollector(fmt.Sprintf(":%d", collectorPort))

	doReq := func() {

		span := appdash.NewRootSpanID()

		dialerSpan := appdash.NewSpanID(span)
		defaultDial := (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial

		traceDial := func(network string, address string) (net.Conn, error) {
			conn, err := defaultDial(network, address)
			id, ok := spanMap.Get(conn)
			if !ok {
				id = appdash.NewSpanID(dialerSpan)
				spanMap.Set(conn, id)
			}

			// TODO(tcm) This span should really be rooted off of the connection
			tconn := traceConn{
				base: conn,
				id:   id,
			}
			return tconn, err
		}

		// A customized version of http.DefaultTransport
		netTraceTransport := &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                traceDial,
			TLSHandshakeTimeout: 10 * time.Second,
		}

		httpClient := &http.Client{
			Transport: &httptrace.Transport{Recorder: appdash.NewRecorder(span, localCollector), SetName: true, Transport: netTraceTransport},
		}

		testdomain := "api.stagingf.we7.com"

		creds := oauth.DefaultCredentials()
		creds.ClientID = os.Getenv("OAUTH_CLIENT_ID")
		creds.ClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
		creds.Username = os.Getenv("OAUTH_USERNAME")
		creds.Password = os.Getenv("OAUTH_PASSWORD")

		oc := oauth.New()
		oc.Option(oauth.HTTPClient(httpClient))
		oc.AddDomain(testdomain, creds)

		c := mg.New()
		c.ApiKey = "test"
		c.Host = testdomain
		c.Option(mg.OAuthClient(oc))

		r, _ := c.SimpleSearch("jimi hendrix purple haze", []string{"tracks"})

		t := r.Tracks[0]
		log.Println(t.Id)
		log.Println(t.Images)

		r, _ = c.SimpleSearch("jimi hendrix purple haze", []string{"tracks"})

		t = r.Tracks[0]
		log.Println(t.Id)
		log.Println(t.Images)

		s, err := c.StreamInfo(t.Id, "RADIO", 0, []string{"MP3"})
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(s)
		log.Println(s.Unique)
		log.Println(s.Location)

		d := (time.Second * 15)
		time.Sleep(d)

		err = c.StreamEnd(s.Unique, d, 0)
		log.Println(err)
	}

	go doReq()
	go doReq()

	time.Sleep(15 * time.Minute)

	return
}

type traceConn struct {
	base net.Conn
	id   appdash.SpanID
}

func (c traceConn) Read(b []byte) (n int, err error) {
	//rid := appdash.NewSpanID(c.id)
	return c.Read(b)
}

func (c traceConn) Write(b []byte) (n int, err error) {
	//wid := appdash.NewSpanID(c.id)
	return c.Write(b)
}

func (c traceConn) Close() error {
	return c.Close()
}

func (c traceConn) LocalAddr() net.Addr {
	return c.LocalAddr()
}

func (c traceConn) RemoteAddr() net.Addr {
	return c.RemoteAddr()
}

func (c traceConn) SetDeadline(t time.Time) error {
	return c.SetDeadline(t)
}

func (c traceConn) SetReadDeadline(t time.Time) error {
	return c.SetReadDeadline(t)
}

func (c traceConn) SetWriteDeadline(t time.Time) error {
	return c.SetWriteDeadline(t)
}

// Thread safe map for tracking connections and spans
type connSpanMap struct {
	lock *sync.RWMutex
	smap map[net.Conn]appdash.SpanID
}

func (m *connSpanMap) Get(c net.Conn) (appdash.SpanID, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	s, ok := m.smap[c]
	return s, ok
}

func (m *connSpanMap) Set(c net.Conn, s appdash.SpanID) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.smap[c] = s
	return
}
