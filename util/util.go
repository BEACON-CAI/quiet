package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"quiet/vars"
	"strconv"
	"strings"
	"time"

	"github.com/malfunkt/iprange"
)

// String ips -> []net.IP
// include "1.1.1.1,1.1.1.1-255,1.1.1.*,1.1.1.1/24"
func GetIpList(ips string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}

	list := addressList.Expand()
	return list, err
}

// []net.IP -> []string
func Net2String(netip []net.IP) []string {
	strip := []string{}
	for _, i := range netip {
		s := i.String()
		strip = append(strip, s)
	}
	return strip
}

// String ports -> []int
// include "1,2,3,4-5"
func GetPortList(strPost string) ([]int, error) {
	intPosts := []int{}
	if strPost == "" {
		return intPosts, nil
	}

	// ","
	ranges := strings.Split(strPost, ",")
	for _, r := range ranges {
		// "space"
		r = strings.TrimSpace(r)
		// "-"
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("[!]Invalid port selection segment: '%s'", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("[!]Invalid port number: '%s'", parts[0])
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("[!]Invalid port number: '%s'", parts[1])
			}

			if p1 > p2 {
				return nil, fmt.Errorf("[!]Invalid port range: %d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				intPosts = append(intPosts, i)
			}

		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("[!]Invalid port number: '%s'", r)
			} else {
				intPosts = append(intPosts, port)
			}
		}
	}
	return intPosts, nil
}

// LocalAddr
func LocalIPPort(dstip net.IP) (net.IP, int, error) {
	// func ResolveUDPAddr(network, address string) (*UDPAddr, error)
	// ResolveUDPAddr returns an address of UDP end point.
	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":30274")
	if err != nil {
		return nil, 0, err
	}

	// func DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error)
	// DialUDP acts like Dial for UDP networks.
	// If the Port field of gaddr is 0, a port number is automatically chosen.
	// So we will get a random local port.
	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		// func (c *IPConn) LocalAddr() Addr
		// LocalAddr returns the local network address.
		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			return udpaddr.IP, udpaddr.Port, nil
		}
	}
	return nil, -1, err
}

// // Local address (TCP) unusefulness
// func LocalIPPortTCP(dstip string) (net.IP, int, error) {
// 	// func ResolveTCPAddr(net, addr string) (*TCPAddr, error)
// 	// ip:port must open
// 	serverAddr, err := net.ResolveTCPAddr("tcp", dstip+":53")
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	// func DialTCP(net string, laddr, raddr *TCPAddr) (*TCPConn, error)
// 	if con, err := net.DialTCP("tcp", nil, serverAddr); err == nil {
// 		// local address
// 		if tcpaddr, ok := con.LocalAddr().(*net.TCPAddr); ok {
// 			return tcpaddr.IP, tcpaddr.Port, nil
// 		}
// 	}
// 	return nil, -1, err
// }

// uint16 random number
func RandomNum() uint16 {
	rand.Seed(time.Now().UnixNano())
	number := uint16(rand.Intn(65536))
	return number
}

// is root?
func IsRoot() bool {
	return os.Geteuid() == 0
}

// Check whether you are an administrator
// Some functions must be run as an administrator.
func CheckRoot() {
	if !IsRoot() {
		fmt.Println("[!]Some functions of this operation must be run as an administrator.")
		os.Exit(0)
	}
}

// md5 function
func MD5(s string) (m string) {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func MakeTaskHash(k string) string {
	hash := MD5(k)
	return hash
}

func CheckTaskHash(hash string) bool {
	_, ok := vars.SuccessHash.Load(hash)
	return ok
}

func SetTaskHash(hash string) {
	vars.SuccessHash.Store(hash, true)
}

