package port

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Network string
type NetworkState string

type ScanResult struct {
	Protocol Network
	Port     uint16
	Result   bool
}

const (
	TCP Network = "tcp"
	UDP Network = "udp"
)

const (
	OPEN   = "open"
	CLOSED = "closed"
)

// Scan a port to check if it is active for a given Network protocol
func Scan(network Network, hostname string, port uint16) (ScanResult, error) {

	var address = hostname + ":" + strconv.FormatUint(uint64(port), 10)
	var conn, err = net.DialTimeout(string(network), address, 60*time.Second)

	if err != nil {
		return ScanResult{
			Protocol: network,
			Port:     port,
			Result:   false,
		}, nil
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("port::Scan() Error closing the connection")
		}
	}(conn)

	return ScanResult{
		Protocol: network,
		Port:     port,
		Result:   true,
	}, nil
}

// ScanRange of ports to check if it is active for a given Network protocol
func ScanRange(network Network, hostname string, portBegin, portEnd uint16) []ScanResult {
	var result []ScanResult
	for i := portBegin; i < portEnd; i++ {
		var res, err = Scan(network, hostname, i)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, res)
	}
	return result
}
