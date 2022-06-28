package vars

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/cheggaaa/pb.v2"
)

// normal
var (
	Timeout   = 2
	ThreadNum = 1000
	StartTime time.Time
)

// progress bar
var (
	// port scan
	ProgressBarPS *pb.ProgressBar
	// when cracking password ,check the port is open
	ProgressBarPC *pb.ProgressBar
	// cracking password
	ProgressBarCP *pb.ProgressBar
)

// port scan or icmp scan
var (
	PortScanResult *sync.Map
	SrcIP          net.IP
	Host           string

	Port    = []int{61616, 50070, 50000, 37777, 27017, 11211, 9999, 9418, 9200, 9100, 9092, 9042, 9001, 9000, 8686, 8545, 8443, 8081, 8080, 7077, 7001, 6379, 6000, 5984, 5938, 5900, 5672, 5601, 5555, 5432, 5222, 5000, 4730, 3389, 3306, 3128, 2379, 2375, 2181, 2049, 1883, 1521, 1433, 1099, 1080, 902, 873, 636, 623, 548, 515, 500, 465, 445, 443, 389, 139, 135, 123, 111, 110, 80, 53, 25, 23, 22, 21}
	SrcPort = 30274

	PortScanMode     = "tcp"
	ModeFlag         = "TCP connection mode"
	UseToTestLocalIP = "114.114.114.114"
)

// port scan
func init() {
	PortScanResult = &sync.Map{}
}

// icmp scan
var (
	ICMPHost = []string{}
)

// password cracking
var (
	PortNames = map[int]string{
		22:    "SSH",
		3306:  "MYSQL",
		6379:  "REDIS",
		1433:  "MSSQL",
		5432:  "POSTGRESQL",
		27017: "MONGODB",
	}

	PCIPList = "ip_list.txt"
	UserDict = "user.txt"
	PassDict = "paswd.txt"

	// save result in a cache
	CacheService *cache.Cache
	ResultFile   = "password_result.txt"

	// Flag whether a particular user of a particular service has been cracked successfully.
	// If so, no further attempts are made to crack the user
	SuccessHash sync.Map

	SupportProtocols map[string]bool
)

func init() {
	CacheService = cache.New(cache.NoExpiration, cache.DefaultExpiration)

	SupportProtocols = make(map[string]bool)
	for _, proto := range PortNames {
		SupportProtocols[strings.ToUpper(proto)] = true
	}
}
