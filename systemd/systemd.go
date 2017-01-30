package systemd

//go:generate go-bindata -pkg $GOPACKAGE templates/

const (
	// DefaultUnitsDirectory represents the directory of user-defined systemd units
	DefaultUnitsDirectory = "/run/systemd/system"
)
