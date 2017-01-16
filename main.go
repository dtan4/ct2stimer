package main

import (
	"fmt"
	"os"

	"github.com/coreos/go-systemd/dbus"
)

func main() {
	conn, err := dbus.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	units, err := conn.ListUnits()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, unit := range units {
		fmt.Println(unit.Name)
	}

	fmt.Println("ct2stimer")
}
