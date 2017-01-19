package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

	for _, schedule := range schedules {
		calendar, err := schedule.ConvertToSystemdCalendar()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		name := "cron-" + schedule.SHA256Sum()[0:12]

		service, err := systemd.GenerateService(name, schedule.Command)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		servicePath := filepath.Join(outdir, name+".service")
		if ioutil.WriteFile(servicePath, []byte(service), 0644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		timer, err := systemd.GenerateTimer(name, calendar)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		timerPath := filepath.Join(outdir, name+".timer")
		if ioutil.WriteFile(timerPath, []byte(timer), 0644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
