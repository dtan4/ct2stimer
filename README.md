# ct2stimer

[![Build Status](https://travis-ci.org/dtan4/ct2stimer.svg?branch=master)](https://travis-ci.org/dtan4/ct2stimer)
[![codecov](https://codecov.io/gh/dtan4/ct2stimer/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/ct2stimer)

Convert crontab to systemd timer

```bash
ubuntu@ubuntu-xenial:~/src/github.com/dtan4/ct2stimer$ sudo ct2stimer -f sample.cron --reload
ubuntu@ubuntu-xenial:~/src/github.com/dtan4/ct2stimer$ systemctl list-timers
NEXT                         LEFT                   LAST PASSED UNIT                         ACTIVATES
Fri 2017-01-20 07:50:00 UTC  4min 16s left          n/a  n/a    cron-77e2fb273c45.timer      cron-77e2fb273c45.service
Fri 2017-01-20 07:56:01 UTC  10min left             n/a  n/a    systemd-tmpfiles-clean.timer systemd-tmpfiles-clean.service
Fri 2017-01-20 08:00:00 UTC  14min left             n/a  n/a    cron-1b33d99b7dda.timer      cron-1b33d99b7dda.service
Fri 2017-01-20 10:00:00 UTC  2h 14min left          n/a  n/a    cron-b60fe106ef63.timer      cron-b60fe106ef63.service
Fri 2017-01-20 12:16:09 UTC  4h 30min left          n/a  n/a    snapd.refresh.timer          snapd.refresh.service
Fri 2017-01-20 19:11:59 UTC  11h left               n/a  n/a    apt-daily.timer              apt-daily.service
Wed 2017-02-01 00:00:00 UTC  1 weeks 4 days left    n/a  n/a    cron-fcd6d8377d9d.timer      cron-fcd6d8377d9d.service
Sat 2017-12-02 01:23:00 UTC  10 months 11 days left n/a  n/a    cron-d3c507cb2439.timer      cron-d3c507cb2439.service

8 timers listed.
Pass --all to see loaded but inactive timers, too.
```

## Installation

TBD

## Usage

ct2stimer reads crontab file at `/etc/crontab` by default. You can specify crontab file with `-f FILE` flag.

systemd unit file are saved at `/etc/systemd/system` by default. You can specify save directory with `-o OUTDIR` flag.

```bash
$ ct2stimer
$ ct2stimer -f sample.cron -o unitfiles
```

### Reload systemd and start all timers automatically

If `--reload` is provided, ct2stimer reloads systemd unit files (= `systemctl daemon-reload`) and starts all generated timers (= `systemctl start foo.timer`). Maybe `sudo` is required to execute.

```bash
$ sudo ct2stimer -f sample.cron --reload
```

### Determine unit name from command to execute

As you know, crontab does not have the concept of "task name". However, task name is required to identify each systemd unit.
You can extract task name from original command using regular expression. `--name-regexp REGEXP` flag is used for this.
Regular expression must have one [capturing group](http://www.regular-expressions.info/brackets.html).

If regular expression is not provided or command does not match to the given regular expression, hash value, which is calculated from command, is used for unit name.

```bash
$ ct2stimer -f sample.cron --name-regexp '--name ([a-zA-Z0-9_-]+)'
```

### Delete unregistered unit files

If `--delete` is provided, ct2timer deletes unit files which are no longer written in the given crontab file.

```bash
$ ct2stimer -f tmp/scheduler -o /run/systemd/system --delete
Deleted: /run/systemd/system/cron-19fb9c164fe8.service
Deleted: /run/systemd/system/cron-19fb9c164fe8.timer
Deleted: /run/systemd/system/cron-4f76a3902132.service
Deleted: /run/systemd/system/cron-4f76a3902132.timer
```

### Specify unit dependencies

You can specify unit dependencies (`After=`) with `--after AFTER` flag.

```bash
$ ct2stimer -f sample.crom --after docker.service
```

## Development

Building and executing on Ubuntu 16.04 VM is easy so that macOS does not have systemd.

```bash
$ go get -d github.com/dtan4/ct2stimer
$ cd $GOPATH/src/github.com/dtan4/ct2stimer
$ vagrant up
$ vagrant ssh

ubuntu@ubuntu-xenial:~/src/github.com/dtan4/ct2timer$ make deps
ubuntu@ubuntu-xenial:~/src/github.com/dtan4/ct2timer$ make
ubuntu@ubuntu-xenial:~/src/github.com/dtan4/ct2timer$ bin/ct2stimer
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
