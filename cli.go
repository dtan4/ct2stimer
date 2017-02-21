package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/dtan4/ct2stimer/crontab"
	"github.com/dtan4/ct2stimer/systemd"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

const (
	exitCodeOK    = 0
	exitCodeError = 1
)

var opts = struct {
	after      string
	delete     bool
	dryRun     bool
	filename   string
	nameRegexp string
	outdir     string
	reload     bool
	user       string
	version    bool
}{}

func parseArgs(args []string) error {
	f := flag.NewFlagSet("ct2stimer", flag.ExitOnError)

	f.StringVar(&opts.after, "after", "", "unit dependencies (After=)")
	f.BoolVar(&opts.delete, "delete", false, "delete unused unit files")
	f.BoolVar(&opts.dryRun, "dry-run", false, "dry run")
	f.StringVarP(&opts.filename, "file", "f", crontab.DefaultCrontabFilename, "crontab file")
	f.StringVar(&opts.nameRegexp, "name-regexp", "", "regexp to extract scheduler name from crontab")
	f.StringVarP(&opts.outdir, "outdir", "o", systemd.DefaultUnitsDirectory, "directory to save systemd files")
	f.BoolVar(&opts.reload, "reload", false, "reload & start genreated timers")
	f.StringVar(&opts.user, "user", "", "unix username who executes process")
	f.BoolVarP(&opts.version, "version", "v", false, "print version")

	f.Parse(args)

	if opts.version {
		return nil
	}

	if opts.filename == "" {
		return errors.New("crontab file is required")
	}

	if opts.outdir == "" {
		return errors.New("directory to save systemd files is required")
	}

	if _, err := os.Stat(opts.outdir); err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("directory %q does not exist", opts.outdir)
		}

		return errors.Wrapf(err, "failed to read directory %q", opts.outdir)
	}

	return nil
}

func run(args []string) int {
	if err := parseArgs(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}

	if opts.version {
		printVersion()
		return exitCodeOK
	}

	body, err := ioutil.ReadFile(opts.filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}

	schedules, err := crontab.Parse(string(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}

	var re *regexp.Regexp

	if opts.nameRegexp == "" {
		re = nil
	} else {
		var err error

		re, err = regexp.Compile(opts.nameRegexp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitCodeError
		}
	}

	scMap := map[string]*crontab.Schedule{}

	for _, schedule := range schedules {
		name := getScheduleName(schedule, re)

		if sc, ok := scMap[name]; ok {
			fmt.Fprintln(os.Stderr, fmt.Errorf(`Schedule name %q already exists. Please consider another name regexp.
  Command A: %s
  Command B: %s`, name, sc.Command, schedule.Command))
			return exitCodeError
		}

		scMap[name] = schedule
	}

	timers := []string{}

	for name, schedule := range scMap {
		calendar, err := schedule.ConvertToSystemdCalendar()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitCodeError
		}

		service, err := systemd.GenerateService(name, schedule.Command, opts.after, opts.user)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitCodeError
		}

		servicePath := filepath.Join(opts.outdir, name+".service")

		if opts.dryRun {
			fmt.Printf("[dry-run] %q will be created\n", servicePath)
		} else {
			if err := ioutil.WriteFile(servicePath, []byte(service), 0644); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return exitCodeError
			}
		}

		timer, err := systemd.GenerateTimer(name, calendar)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitCodeError
		}

		timerPath := filepath.Join(opts.outdir, name+".timer")

		if opts.dryRun {
			fmt.Printf("[dry-run] %q will be created\n", timerPath)
		} else {
			if err := ioutil.WriteFile(timerPath, []byte(timer), 0644); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return exitCodeError
			}
		}

		timers = append(timers, name+".timer")
	}

	if opts.delete {
		deleted, err := deleteUnusedUnits(opts.outdir, scMap)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitCodeError
		}

		for _, path := range deleted {
			fmt.Printf("Deleted: %s\n", path)
		}
	}

	if opts.reload && !opts.dryRun {
		if err := reloadSystemd(timers); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitCodeError
		}
	}

	return exitCodeOK
}
