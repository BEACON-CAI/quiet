package util

import (
	"fmt"
	"net"
	"sync"
	"log"
	"time"

	"quiet/crackpwd/models"
	"quiet/vars"

	"gopkg.in/cheggaaa/pb.v2"
)

var (
	AliveAddr []models.Target
	mutex     sync.Mutex
)

func init() {
	AliveAddr = make([]models.Target, 0)
}

func CheckAlive(ipList []models.Target) []models.Target {
	log.Printf("checking ip active")
	vars.ProgressBarPC= pb.StartNew(len(ipList))
	vars.ProgressBarPC.SetTemplate(`{{ rndcolor "Checking progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor}}  {{rtime . | rndcolor }}`)

	var wg sync.WaitGroup
	wg.Add(len(ipList))

	for _, addr := range ipList {
		go func(addr models.Target) {
			defer wg.Done()
			SaveAddr(check(addr))
		}(addr)
	}
	wg.Wait()
	vars.ProgressBarPC.Finish()

	return AliveAddr
}

func check(Target models.Target) (bool, models.Target) {
	alive := false
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", Target.IP, Target.Port),time.Duration(vars.Timeout) * time.Second)
	if err == nil {
		alive = true
	}

	vars.ProgressBarPC.Increment()
	return alive, Target
}

func SaveAddr(alive bool, Target models.Target) {
	if alive {
		mutex.Lock()
		AliveAddr = append(AliveAddr, Target)
		mutex.Unlock()
	}
}