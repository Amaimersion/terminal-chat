package network

import (
	"errors"
	"net"
)

// LookupSystemNetwork performs lookup of all available host
// IP addresses in the current system network. It returns a
// list of unicast IPv4. These IP's can be used to dial this host
// from another host in the same system network.
func LookupSystemNetwork() ([]net.IP, error) {
	ifis, err := net.Interfaces()

	if err != nil {
		return nil, err
	}

	result := make([]net.IP, 0)

	for _, ifi := range ifis {
		addrs, err := ifi.Addrs()

		if err != nil {
			continue
		}

		for _, addr := range addrs {
			s := addr.String()
			ip, _, err := net.ParseCIDR(s)

			if err != nil {
				continue
			}

			ip = ip.To4()

			if ip == nil {
				continue
			}

			isUnicast :=
				ip.IsGlobalUnicast() ||
					ip.IsLinkLocalUnicast()

			if !isUnicast {
				continue
			}

			result = append(result, ip)
		}
	}

	return result, nil
}

// LookupOutbound performs lookup of host IP address that will be
// used by default for outbound requests to external networks.
func LookupOutbound() (net.IP, error) {
	// See this for how it works -
	// https://stackoverflow.com/a/37382208/8445442

	conn, err := net.Dial("udp", "1.1.1.1:80")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)

	if !ok {
		return nil, errors.New("invalid local addr")
	}

	return localAddr.IP, nil
}
