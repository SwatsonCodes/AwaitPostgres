package reachability

import (
	"net"
	"time"
)

// HostTimeout specifies a host to connect to within a given timeout
type HostTimeout struct {
	Host    string
	Timeout time.Duration
}

// IsHostReachable tests if a host is reachable
func (hostTimeout HostTimeout) IsHostReachable() error {
	conn, err := net.DialTimeout("tcp", hostTimeout.Host, hostTimeout.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
