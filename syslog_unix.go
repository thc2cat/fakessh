//go:build linux || freebsd
// +build linux freebsd

package main

import (
	"log"
	"log/syslog"
)

var xlog *syslog.Writer // Used by initSyslog
func initSyslog(exename string) {

	var err error
	xlog, err = syslog.New(syslog.LOG_DAEMON|syslog.LOG_INFO, exename)
	if err == nil {
		log.SetOutput(xlog)
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime)) // remove timestamp
	}

}
