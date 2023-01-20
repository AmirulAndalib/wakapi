package utils

import (
	"net"
	"strings"

	"github.com/emvi/logbuch"
)

// CheckEmailMX takes an e-mail address and verifies that an MX DNS record exists for its domain
func CheckEmailMX(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	records, err := net.LookupMX(parts[1])
	if len(records) == 0 {
		logbuch.Info("0 dns records")
	}
	if err != nil {
		logbuch.Info("dns:" + err.Error())
	}
	return len(records) > 0 && err == nil
}
