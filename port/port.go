package port

import (
	"errors"
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
	Result   NetworkState
}

const (
	TCP Network = "tcp"
	//UDP Network = "udp"
)

const (
	OPEN   NetworkState = "open"
	CLOSED NetworkState = "closed"
)

// TODO add support for other network protocols

// Scan a port to check if it is active for a given Network protocol
func Scan(network Network, hostname string, port uint16) (ScanResult, error) {

	var address = hostname + ":" + strconv.FormatUint(uint64(port), 10)
	var conn, err = net.DialTimeout(string(network), address, 60*time.Second)

	if err != nil {
		return ScanResult{
			Protocol: network,
			Port:     port,
			Result:   CLOSED,
		}, nil
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("port::scan() Error closing the connection")
		}
	}(conn)

	return ScanResult{
		Protocol: network,
		Port:     port,
		Result:   OPEN,
	}, nil
}

// ScanRange scans a range of ports to check if it is active for a given Network protocol
func ScanRange(network Network, hostname string, portBegin, portEnd uint16) []ScanResult {
	return scanRange(network, hostname, portBegin, portEnd, false)
}

// ScanRangeDetailed scans a range of ports to check if it is active for a given Network protocol
// and prints the progress of the scan if it scans more than 100 ports
// also prints every open port it finds
func ScanRangeDetailed(network Network, hostname string, portBegin, portEnd uint16) []ScanResult {
	return scanRange(network, hostname, portBegin, portEnd, true)
}

func scanRange(network Network, hostname string, portBegin, portEnd uint16, isDetailed bool) []ScanResult {
	var result []ScanResult
	var progress, err = progressMap(portEnd - portBegin)
	for port := portBegin; port < portEnd; port++ {

		if err == nil {
			var curr = progress[portEnd-port]
			if isDetailed && curr != "" {
				fmt.Println(curr)
			}
		}

		var res, err = Scan(network, hostname, port)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if isDetailed && res.Result == OPEN {
			fmt.Printf("port::scan() Found open port: %d\n", port)
		}
		result = append(result, res)
	}
	return result
}

func progressMap(n uint16) (map[uint16]string, error) {

	if n < 100 {
		return nil, errors.New("not enough to create progress")
	}

	var result = make(map[uint16]string)
	for i := uint16(10); i <= 90; i += 10 {
		var s = strconv.FormatUint(uint64(i), 10)
		var numForPercent = percent(i, n)
		result[numForPercent] = "Progress: " + s + "%"
	}
	return result, nil
}

func percent(percent uint16, all uint16) uint16 {
	return uint16((float64(all) * float64(percent)) / float64(100))
}
