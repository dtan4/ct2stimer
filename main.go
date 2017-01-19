package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/dtan4/ct2stimer/crontab"
	"github.com/dtan4/ct2stimer/systemd"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Please specify crontab file.")
		os.Exit(1)
	}
	filename := os.Args[1]

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
