# ct2stimer

[![Build Status](https://travis-ci.org/dtan4/ct2stimer.svg?branch=master)](https://travis-ci.org/dtan4/ct2stimer)

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
