package scanner

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"quiet/vars"

	"gopkg.in/cheggaaa/pb.v2"
)

// Generate port scan task
func GeneratePSTask(ipList []net.IP, ports []int) ([]map[string]int, int) {
	// []map[string]int
	tasks := make([]map[string]int, 0)

	// Generate a scan task list
	for _, ip := range ipList {
		for _, port := range ports {
			ipPort := map[string]int{ip.String(): port}
			tasks = append(tasks, ipPort)
		}
	}

	return tasks, len(tasks)
}

// Port scan task
func PortScanTask(taskChan chan map[string]int, wg *sync.WaitGroup) {
	for task := range taskChan {
		vars.ProgressBarPS.Increment()
		for ip, port := range task {
			if strings.ToLower(vars.PortScanMode) == "syn" {
				err := SavePSResult(TcpSYN(vars.SrcIP, vars.SrcPort, ip, port))
				_ = err
			} else {
				err := SavePSResult(TcpConnect(ip, port))
				_ = err
			}
			wg.Done()
		}
	}
}

// Run port scan task
func RunPSTask(tasks []map[string]int, totalTask int) {
	vars.ProgressBarPS = pb.StartNew(totalTask)
	vars.ProgressBarPS.SetTemplate(`Custom template: {{counters . }}`)
	wg := &sync.WaitGroup{}

	// Create a channel which buffer is vars.threadNum
	taskChan := make(chan map[string]int, vars.ThreadNum)

	// Create vars.ThreadNum coroutines
	for i := 0; i < vars.ThreadNum; i++ {
		go PortScanTask(taskChan, wg)
	}

	// The producer keeps sending data to taskChan until taskChan blocks.
	// The coroutine is first created,
	// then communicated with and transmitted to the coroutine
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()
}

// Save port scan result
func SavePSResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}
	// When local port scan,
	// even if scan to a not open port,
	// err = nil, so scanner.TcpSYN() will return 0
	if port != 0 {
		log.Printf("Open port found: %s:%d", ip, port)
		v, ok := vars.PortScanResult.Load(ip)
		if ok {
			ports, ook := v.([]int)
			if ook {
				ports = append(ports, port)
				vars.PortScanResult.Store(ip, ports)
			}
		} else {
			ports := make([]int, 0)
			ports = append(ports, port)
			vars.PortScanResult.Store(ip, ports)
		}
	}
	return err
}

// Print port scan result
func PrintPSResult(result *sync.Map) {
	fmt.Println("\nPort scan completed!")
	result.Range(func(key, value interface{}) bool {
		fmt.Printf("%v : ", key)
		fmt.Printf("%v\n", value)
		return true
	})
}
