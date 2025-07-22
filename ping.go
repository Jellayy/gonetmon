package main

import (
	"log"
	"net"
	"os"
	"time"
	"strings"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ianaProtocolICMP = 1
)

func ping(host string, count int) (float32, error) {
	// resolve host + create icmp message
	dst, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return -2, err
	}
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return -2, err
	}
	defer conn.Close()
	message := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff,
			Seq: 1,
			Data: []byte("Hello, World!"),
		},
	}
	data, err := message.Marshal(nil)
	if err != nil {
		return -2, err
	}

	// send and listen for ping count
	pingTimes := make([]time.Duration, count)
	for i := 0; i < count; i++ {
		start := time.Now()
		_, err = conn.WriteTo(data, dst)
		if err != nil {
			return -2, err
		}

		reply := make([]byte, 1500)
		err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			return -2, err
		}

		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				log.Printf("Ping to %v timed out", host)
				return -1, nil
			}
			return -2, err
		}
		duration := time.Since(start)

		rm, err := icmp.ParseMessage(ianaProtocolICMP, reply[:n])
		if err != nil {
			return -2, err
		}

		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			log.Printf("Ping reply from %v: time=%v\n", peer, duration)
			pingTimes[i] = duration
		default:
			log.Printf("Got %+v from %v\n", rm, peer)
		}
	}

	// average times
	var pingTimeTotal time.Duration
	for _, d := range pingTimes {
		pingTimeTotal += d
	}
	avgPingTime := pingTimeTotal / time.Duration(len(pingTimes))
	log.Printf("Average reply time: %v\n", avgPingTime)

	return float32(avgPingTime.Milliseconds()), nil
}
