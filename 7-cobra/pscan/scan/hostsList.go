package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists   = fmt.Errorf("host already in the list")
	ErrNotExist = fmt.Errorf("host not in the list")
)

// HostsList represents a list of hosts to run port scan on
type HostsList struct {
	Hosts []string
}

// Save saves hosts to a hosts file
func (hl *HostsList) Save(hostsFile string) error {
	output := "" // single string to carry all data seperated by \n
	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}
	return os.WriteFile(hostsFile, []byte(output), 0644) // write to file and return error if any
}

// Load - read and load hosts file content as slice of string
func (hl *HostsList) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close() // memory leak prevent

	// new scanner of file, read line by line and add to file
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}

	return nil
}

// Add - adds a host to the list
func (hl *HostsList) Add(host string) error {
	// check if exists in list already
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}

	// if not then add to list
	hl.Hosts = append(hl.Hosts, host)
	return nil
}

// Remove - deletes an existing host
func (hl *HostsList) Remove(host string) error {
	if found, i := hl.search(host); found {
		hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
		return nil
	}

	// if nothing matches, that means we are trying to remove element that was never part of original list
	return fmt.Errorf("%w: %s", ErrNotExist, host)

}

// search - searches for host in the list
func (hl *HostsList) search(host string) (bool, int) {
	// sort string first
	sort.Strings(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}
	return false, -1
}
