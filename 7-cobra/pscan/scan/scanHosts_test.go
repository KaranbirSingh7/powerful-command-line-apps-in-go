package scan_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/karanbirsingh7/pscan/scan"
)

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{
			"OpenPort", "open",
		},
		{
			"ClosedPort", "closed",
		},
	}

	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)
	ports := []int{} // ports to scan

	for _, tc := range testCases {

		// start a local TCP server (0 will pick any available port)
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close() // close server when exiting loop

		// grab TCP server port number
		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		// convert into int
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		// add TCP local server port to list of "ports to scan"
		ports = append(ports, port)

		// when running ClosedPort test case, close server before asserting checks
		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	// run scan on local TCP port (open or closed)
	res := scan.Run(hl, ports)

	// verify results for HostFound test

	// expecting single result because we running test on 1 host (localhost) only
	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}

	// host should be localhost
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}

	// host should be found since  "localhost" is a valid host
	if res[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d instead\n", ports[0], res[0].PortStates[i].Port)
		}
		if res[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectState)
		}
	}
}

func TestStateString(t *testing.T) {
	ps := scan.PortState{} // empty object meaning port should be closed by default

	// assert the closed port
	if ps.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q", "closed", ps.Open.String())
	}

	ps.Open = true // open the port

	// assert the open port
	if ps.Open.String() != "open" {
		t.Errorf("Expected %q, got %q", "open", ps.Open.String())
	}

}

func TestRunHostNotFound(t *testing.T) {
	host := "389.389.389.389" // non existing host in IPV4 range
	hl := &scan.HostsList{}

	hl.Add(host) // add to our list

	// since host is non-existent, there is no need to pass any ports because it would never go to that point ever
	res := scan.Run(hl, []int{})

	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}

	// host was non-existent so it should retur NOTFOUND
	if !res[0].NotFound {
		t.Errorf("Expected host %q NOT to be found\n", host)
	}

	// no ports passed so not expected in return
	if len(res[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port states, got %d instead\n", len(res[0].PortStates))
	}
}
