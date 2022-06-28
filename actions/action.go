package actions

import (
	"fmt"
	"log"
	"net"
	"quiet/scanner"
	"quiet/util"
	"quiet/vars"
	"strings"

	"github.com/urfave/cli/v2"
)

func PortScan(ctx *cli.Context) error {
	// ip
	if ctx.IsSet("iplist") {
		vars.Host = ctx.String("iplist")
	}
	if ctx.IsSet("local") {
		if i, _, err := util.LocalIPPort(net.ParseIP(vars.UseToTestLocalIP)); err != nil {
			log.Fatal(err)
		} else {
			vars.Host = i.String()
		}
	}

	ips, err := util.GetIpList(vars.Host)
	if err != nil {
		log.Fatal(err)
	}

	if ctx.IsSet("port") {
		ports, err := util.GetPortList(ctx.String("port"))
		if err != nil {
			log.Fatal(err)
		}
		vars.Port = ports
	}

	if ctx.IsSet("mode") {
		vars.PortScanMode = ctx.String("mode")
	}

	if ctx.IsSet("timeout") {
		vars.Timeout = ctx.Int("timeout")
	}

	if ctx.IsSet("concurrency") {
		vars.ThreadNum = ctx.Int("concurrency")
	}

	if strings.ToLower(vars.PortScanMode) == "syn" {
		// TCP SYN mode need root
		vars.ModeFlag = "TCP SYN mode"
		util.CheckRoot()
		// Local ip port
		i, p, err := util.LocalIPPort(ips[0])
		if err != nil {
			log.Fatal(err)
		}
		vars.SrcIP = i
		vars.SrcPort = p
	}

	tasks, totalTask := scanner.GeneratePSTask(ips, vars.Port)
	// _ = n
	fmt.Println("Port scanning...")
	fmt.Printf("Total tasks: %v | Port scan mode: %s | Timeout: %d seconds\n", totalTask, vars.ModeFlag, vars.Timeout)
	// run
	scanner.RunPSTask(tasks, totalTask)
	// result
	scanner.PrintPSResult(vars.PortScanResult)
	return err
}

func ICMPScan(ctx *cli.Context) error {
	// ICMP must run as root
	util.CheckRoot()
	// ip
	if ctx.IsSet("iplist") {
		netips, err := util.GetIpList(ctx.String("iplist"))
		if err != nil {
			log.Fatal(err)
		}
		strips := util.Net2String(netips)
		vars.ICMPHost = strips
	}
	if ctx.IsSet("local") {
		if i, _, err := util.LocalIPPort(net.ParseIP(vars.UseToTestLocalIP)); err != nil {
			log.Fatal(err)
		} else {
			netips, err := util.GetIpList(i.String() + "/24")
			if err != nil {
				log.Fatal(err)
			}
			strips := util.Net2String(netips)
			vars.ICMPHost = strips
		}
	}

	if ctx.IsSet("domain") {
		vars.ICMPHost = append(vars.ICMPHost, ctx.String("domain"))
	}

	if ctx.IsSet("timeout") {
		vars.Timeout = ctx.Int("timeout")
	}

	if ctx.IsSet("concurrency") {
		vars.ThreadNum = ctx.Int("concurrency")
	}

	n := len(vars.ICMPHost)
	// _ = n
	fmt.Println("ICMP finding...")
	fmt.Printf("Total tasks: %v | Timeout: %d seconds\n", n, vars.Timeout)

	scanner.RunICMPTask(vars.ICMPHost)

	var err error
	return err
}
