package flag

import (
	"quiet/actions"
	"quiet/crackpwd"

	"github.com/urfave/cli/v2"
)

// ./main ps --iplist ip_list --port port_list --mode syn  --timeout 2 --concurrency 10
var PortScanCom = &cli.Command{
	Name:        "portscan",
	Usage:       "tcp syn/connect port scanner",
	Aliases:     []string{"ps", "p", "port"},
	Description: "start to scan port",
	Action:      actions.PortScan,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "iplist",
			Aliases: []string{"ip", "i"},
			Value:   "",
			Usage:   "ip list",
		},
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "port list",
		},
		&cli.StringFlag{
			Name:    "mode",
			Aliases: []string{"m"},
			Value:   "",
			Usage:   "port scan mode",
		},
		&cli.IntFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   2,
			Usage:   "timeout",
		},
		&cli.IntFlag{
			Name:    "concurrency",
			Aliases: []string{"c"},
			Value:   1000,
			Usage:   "concurrency",
		},
		&cli.BoolFlag{
			Name:    "local",
			Aliases: []string{"l"},
			Usage:   "local port scan",
		},
	},
}

// ./main ping --iplist ip_list --timeout 2 --concurrency 10
// ./main ping --domain domain --timeout 2 --concurrency 10
var ICMPScanCom = &cli.Command{
	Name:        "ICMPscan",
	Usage:       "ICMP scanner",
	Aliases:     []string{"icmpscan", "is", "ping"},
	Description: "start to ping a host",
	Action:      actions.ICMPScan,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "iplist",
			Aliases: []string{"ip", "i"},
			Usage:   "ip list",
		},
		&cli.StringFlag{
			Name:    "domain",
			Aliases: []string{"d"},
			Usage:   "domain",
		},
		&cli.IntFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   2,
			Usage:   "timeout",
		},
		&cli.IntFlag{
			Name:    "concurrency",
			Aliases: []string{"c"},
			Value:   1000,
			Usage:   "concurrency",
		},
		&cli.BoolFlag{
			Name:    "local",
			Aliases: []string{"l"},
			Usage:   "local ICMP scan",
		},
	},
}

var PasswordCrack = &cli.Command{
	Name:        "PasswordCrack",
	Usage:       "Password crack",
	Aliases:     []string{"crack", "pc"},
	Description: "Password crack",
	Action:      crackpwd.PasswordCrack,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "ip_list",
			Aliases: []string{"ip", "i"},
			Value:   "ip_list.txt",
			Usage:   "ip list",
		},
		&cli.StringFlag{
			Name:    "user_dict",
			Aliases: []string{"ud", "u","user"},
			Value:   "user.txt",
			Usage:   "user dict",
		},
		&cli.StringFlag{
			Name:    "pass_dict",
			Aliases: []string{"pd", "p","password"},
			Value:   "paswd.txt",
			Usage:   "password dict",
		},
		&cli.StringFlag{
			Name:    "outfile",
			Aliases: []string{"of", "o"},
			Value:   "password_result.txt",
			Usage:   "password result",
		},
		&cli.IntFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   2,
			Usage:   "timeout",
		},
		&cli.IntFlag{
			Name:    "concurrency",
			Aliases: []string{"c"},
			Value:   1000,
			Usage:   "concurrency",
		},
	},
}

