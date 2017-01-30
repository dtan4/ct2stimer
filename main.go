package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/dtan4/ct2stimer/crontab"
	"github.com/dtan4/ct2stimer/systemd"
	flag "github.com/spf13/pflag"
)

var opts = struct {
	after      string
	dryRun     bool
	filename   string
	nameRegexp string
	outdir     string
	reload     bool
	version    bool
}{}

func parseArgs(args []string) error {
	f := flag.NewFlagSet("ct2stimer", flag.ExitOnError)

	f.StringVar(&opts.after, "after", "", "unit dependencies (After=)")
	f.BoolVar(&opts.dryRun, "dry-run", false, "dry run")
	f.StringVarP(&opts.filename, "file", "f", crontab.DefaultCrontabFilename, "crontab file")
	f.StringVar(&opts.nameRegexp, "name-regexp", "", "regexp to extract scheduler name from crontab")
	f.StringVarP(&opts.outdir, "outdir", "o", systemd.DefaultUnitsDirectory, "directory to save systemd files")
	f.BoolVar(&opts.reload, "reload", false, "reload & start genreated timers")
	f.BoolVarP(&opts.version, "version", "v", false, "print version")

	f.Parse(args)

	if opts.version {
		return nil
	}

	if opts.filename == "" {
		return fmt.Errorf("Please specify crontab file.")
	}

	if opts.outdir == "" {
		return fmt.Errorf("Please specify directory to save systemd files.")
	}

	if _, err := os.Stat(opts.outdir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: directory does not exist", opts.outdir)
		}

		return err
	}

	return nil
}

func getScheduleName(schedule *crontab.Schedule, re *regexp.Regexp) string {
	name := schedule.NameByRegexp(re)
	if name == "" {
		name = "cron-" + schedule.SHA256Sum()[0:12]
	}

	return name
}

func reloadSystemd(timers []string) error {
	conn, err := systemd.NewConn()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	client := systemd.NewClient(conn)

	if err := client.Reload(); err != nil {
		return err
	}

	for _, timerUnit := range timers {
		if err := client.StartUnit(timerUnit); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := parseArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if opts.version {
		printVersion()
		os.Exit(0)
	}

	body, err := ioutil.ReadFile(opts.filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	schedules, err := crontab.Parse(string(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var re *regexp.Regexp

	if opts.nameRegexp == "" {
		re = nil
	} else {
		var err error

		re, err = regexp.Compile(opts.nameRegexp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	scMap := map[string]*crontab.Schedule{}

	for _, schedule := range schedules {
		name := getScheduleName(schedule, re)

		if sc, ok := scMap[name]; ok {
			fmt.Fprintln(os.Stderr, fmt.Errorf(`Schedule name %q already exists. Please consider another name regexp.
  Command A: %s
  Command B: %s`, name, sc.Command, schedule.Command))
			os.Exit(1)
		}

		scMap[name] = schedule
	}

	timers := []string{}

	for name, schedule := range scMap {
		calendar, err := schedule.ConvertToSystemdCalendar()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		service, err := systemd.GenerateService(name, schedule.Command, opts.after)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		servicePath := filepath.Join(opts.outdir, name+".service")

		if opts.dryRun {
			fmt.Printf("[dry-run] %q will be created\n", servicePath)
		} else {
			if err := ioutil.WriteFile(servicePath, []byte(service), 0644); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		timer, err := systemd.GenerateTimer(name, calendar)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		timerPath := filepath.Join(opts.outdir, name+".timer")

		if opts.dryRun {
			fmt.Printf("[dry-run] %q will be created\n", timerPath)
		} else {
			if err := ioutil.WriteFile(timerPath, []byte(timer), 0644); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		timers = append(timers, name+".timer")
	}

	if opts.reload && !opts.dryRun {
		if err := reloadSystemd(timers); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
