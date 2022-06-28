package crackpwd

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"quiet/crackpwd/models"
	"quiet/crackpwd/plugins"
	"quiet/util"
	"quiet/vars"

	"github.com/urfave/cli/v2"
	"gopkg.in/cheggaaa/pb.v2"
)

// Generate crack task
func GenerateTask(ipList []models.Target, users []string, passwords []string) (tasks []models.Service, taskNum int) {
	tasks = make([]models.Service, 0)

	for _, user := range users {
		for _, password := range passwords {
			for _, addr := range ipList {
				service := models.Service{IP: addr.IP, Port: addr.Port, Protocol: addr.Protocol, Username: user, Password: password}
				tasks = append(tasks, service)
			}
		}
	}

	return tasks, len(tasks)
}

// Run task
func RunTask(tasks []models.Service) {
	totalTask := len(tasks)
	vars.ProgressBarPC = pb.StartNew(totalTask)
	vars.ProgressBarPC.SetTemplate(`{{ rndcolor "Scanning progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor }} {{rtime . | rndcolor}} `)

	wg := &sync.WaitGroup{}

	taskChan := make(chan models.Service, vars.ThreadNum)

	for i := 0; i < vars.ThreadNum; i++ {
		go CrackPassword(taskChan, wg)
	}

	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	close(taskChan)
	waitTimeout(wg, time.Duration(vars.Timeout)*time.Second)

	// Save result
	{
		_ = util.SaveResultToFile()
		util.ResultTotal()
		_ = util.DumpToFile(vars.ResultFile)
	}

}

// Crack Password
func CrackPassword(taskChan chan models.Service, wg *sync.WaitGroup) {
	for task := range taskChan {
		vars.ProgressBarPC.Increment()

		var k string
		protocol := strings.ToUpper(task.Protocol)

		if protocol == "REDIS" {
			k = fmt.Sprintf("%v-%v-%v", task.IP, task.Port, task.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", task.IP, task.Port, task.Username)
		}

		h := util.MakeTaskHash(k)
		if util.CheckTaskHash(h) {
			wg.Done()
			continue
		}

		fn := plugins.ScanFuncMap[protocol]
		util.SaveResult(fn(task))
		wg.Done()
	}
}

func PasswordCrack(ctx *cli.Context) (err error) {
	if ctx.IsSet("timeout") {
		vars.Timeout = ctx.Int("timeout")
	}

	if ctx.IsSet("scan_num") {
		vars.ThreadNum = ctx.Int("scan_num")
	}

	if ctx.IsSet("ip_list") {
		vars.PCIPList = ctx.String("ip_list")
	}

	if ctx.IsSet("user_dict") {
		vars.UserDict = ctx.String("user_dict")
	}

	if ctx.IsSet("pass_dict") {
		vars.PassDict = ctx.String("pass_dict")
	}

	if ctx.IsSet("outfile") {
		vars.ResultFile = ctx.String("outfile")
	}

	vars.StartTime = time.Now()

	userDict, uErr := util.ReadUserDict(vars.UserDict)
	passDict, pErr := util.ReadPasswordDict(vars.PassDict)

	ipList := util.ReadIpList(vars.PCIPList)

	aliveIpList := util.CheckAlive(ipList)
	if uErr == nil && pErr == nil {
		tasks, _ := GenerateTask(aliveIpList, userDict, passDict)
		RunTask(tasks)
	}
	return err
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
