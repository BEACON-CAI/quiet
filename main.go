package main

import (
	"log"
	"os"
	"quiet/flag"
	"runtime"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "quiet"
	app.Version = "2.0.1"
	app.Usage = "port/icmp scanner"
	app.Commands = []*cli.Command{flag.PortScanCom, flag.ICMPScanCom,flag.PasswordCrack}

	banner := `
             _      _
            ( )    | |
  __ _ _   _ _  ___| |_
 / _\ | | | | |/ _ \ __|
| (_| | |_| | |  __/ |_
 \__, |\__,_|_|\___|\__|
    | |		By：` + aurora.Red("Beacon").String() + `
    |_|		Version: ` + aurora.Green(app.Version).String() + `
Github: ` + aurora.Blue("https://github.com/BEACON-CAI/quiet").String() + ``

	print(banner + "\n\n")
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
