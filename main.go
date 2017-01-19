package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/dtan4/ct2stimer/crontab"
	"github.com/dtan4/ct2stimer/systemd"
	flag "github.com/spf13/pflag"
)

func main() {
	var (
		filename string
		outdir   string
	)

	f := flag.NewFlagSet("ct2stimer", flag.ExitOnError)

	f.StringVarP(&filename, "file", "f", "", "crontab file")
	f.StringVarP(&outdir, "outdir", "o", "", "Directory to save systemd files")

	f.Parse(os.Args[1:])

	if filename == "" {
		fmt.Fprintln(os.Stderr, "Please specify crontab file.")
		os.Exit(1)
	}

	if outdir == "" {
		fmt.Fprintln(os.Stderr, "Please specify directory to save systemd files.")
		os.Exit(1)
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	schedules, err := crontab.Parse(string(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, schedule := range schedules {
		calendar, err := schedule.ConvertToSystemdCalendar()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		name := "systemd" + strconv.Itoa(i)

		service, err := systemd.GenerateService(name, schedule.Command)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		timer, err := systemd.GenerateTimer(name, calendar)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("----- Service")
		fmt.Println(service)
		fmt.Println("")

		fmt.Println("----- Timer")
		fmt.Println(timer)
		fmt.Println("")
	}
}
