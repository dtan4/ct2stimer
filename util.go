package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"./crontab"
	"./systemd"
	"github.com/pkg/errors"
)

const (
	serviceExt = ".service"
	timerExt   = ".timer"
)

func deleteUnusedUnits(outdir string, scMap map[string]*crontab.Schedule) ([]string, error) {
	files, err := ioutil.ReadDir(outdir)
	if err != nil {
		return []string{}, errors.Wrapf(err, "failed to read %q", outdir)
	}

	deleted := []string{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		deletable := false

		if strings.HasSuffix(file.Name(), serviceExt) {
			unitName := strings.TrimSuffix(filepath.Base(file.Name()), serviceExt)

			if _, ok := scMap[unitName]; ok {
				continue
			}

			deletable = true
		} else if strings.HasSuffix(file.Name(), timerExt) {
			unitName := strings.TrimSuffix(filepath.Base(file.Name()), timerExt)

			if _, ok := scMap[unitName]; ok {
				continue
			}

			deletable = true
		} else {
			continue
		}

		if deletable {
			path := filepath.Join(outdir, file.Name())

			if err := os.Remove(path); err != nil {
				return []string{}, errors.Wrapf(err, "failed to delete %q", path)
			}

			deleted = append(deleted, path)
		}
	}

	return deleted, nil
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
		return errors.Wrap(err, "cannot establish new systemd connection")
	}
	defer conn.Close()

	client := systemd.NewClient(conn)

	if err := client.Reload(); err != nil {
		return errors.Wrap(err, "cannot reload systemd unit files")
	}

	for _, timerUnit := range timers {
		if err := client.StartUnit(timerUnit); err != nil {
			return errors.Wrap(err, "cannot reload systemd timer unit")
		}
	}

	return nil
}
