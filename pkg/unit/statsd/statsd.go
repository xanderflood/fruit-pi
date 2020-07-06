package statsd

import (
	"fmt"
	"io"
	"net"
	"time"
)

type StatsD interface {
	Start()
	Stop()

	Count(metric string, value int)
	Time(metric string, took time.Duration)
	Gauge(metric string, value int)
}

type Client struct {
	host  string
	queue chan string
	done  chan struct{}
}

func New(host string) Client {
	return Client{
		host:  host,
		queue: make(chan string, 100),
		done:  make(chan struct{}),
	}
}

func (c Client) Start() {
	go c.start()
}

func (c Client) Stop() {
	close(c.queue)

	// wait for the client to flush
	for range c.done {
	}
}

func (c Client) Count(metric string, value int) {
	c.queue <- fmt.Sprintf("%s:%d|c", metric, value)
}

func (c Client) Time(metric string, took time.Duration) {
	c.queue <- fmt.Sprintf("%s:%d|ms", metric, took/1e6)
}

func (c Client) Gauge(metric string, value int) {
	c.queue <- fmt.Sprintf("%s:%d|g", metric, value)
}

func (c Client) start() {
	defer close(c.done)

	for s := range c.queue {
		if conn, err := net.Dial("udp", "127.0.0.1:8125"); err == nil {
			io.WriteString(conn, s)
			conn.Close()
		}
	}
}
