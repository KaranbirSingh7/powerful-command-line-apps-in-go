package scan

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single TCP port
type PortState struct {
	Port int
	Open state
}

type state bool

type Results struct {
	Host      string
	NotFound  bool //whether we can resolve into an IPv4 addresss
	PortState []PortState
}

// Run performs a port scan on the hosts list
func Run(hl *HostsList, ports []int) []Results {
	res := make([]Results, 0, len(hl.Hosts))
	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}

		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue // move onto next "host" in list
		}

		// scan ports and append results
		for _, p := range ports {
			r.PortState = append(r.PortState, scanPort(r.Host, p, 1*time.Second))
		}

		res = append(res, r)
	}

	return res
}

// String converts the boolean value of state to a human readable string
func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

// scanPort scans a port and returns its state
func scanPort(host string, port int, timeout time.Duration) PortState {
	p := PortState{
		Port: port,
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port)) // also coverts IPv6 addresses, you never know

	scanConn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil { // assumption, error means port is closed
		return p
	}
	scanConn.Close() // success here
	p.Open = true
	return p
}
