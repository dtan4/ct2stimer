# ct2stimer

[![Build Status](https://travis-ci.org/dtan4/ct2stimer.svg?branch=master)](https://travis-ci.org/dtan4/ct2stimer)

Convert crontab to systemd timer

## Developemnt

Building and executing on Ubuntu 16.04 VM is easy so that macOS does not have Systemd.

```bash
$ go get -d github.com/dtan4/ct2stimer
$ cd $GOPATH/src/github.com/dtan4/ct2stimer
$ vagrant up
$ vagrant ssh

ubuntu@ubuntu-xenial:~$ cd $GOPATH/src/github.com/dtan4/ct2timer
ubuntu@ubuntu-xenial:~/src/github.com/dtan4/ct2timer$ make
ubuntu@ubuntu-xenial:~/src/github.com/dtan4/ct2timer$ bin/ct2stimer
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
