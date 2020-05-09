package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"

	uuid "github.com/satori/go.uuid"
)

//GetUUID generates a new UUID
func GetUUID() uuid.UUID {
	UUID := uuid.Must(uuid.NewV4())
	return UUID
}

//CreatePidFile generic utility to create a file with current pid
func CreatePidFile(pidfile string) {
	err := ioutil.WriteFile(pidfile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		panic(err)
	}
}

//OpenLog function opens a syslog handle and sets it to log
func OpenLog(syslogto syslog.Priority, servicename string) {
	logwriter, e := syslog.New(syslogto, servicename)
	if e == nil {
		log.SetOutput(logwriter)
	}
	log.SetFlags(0)
	log.Println("Starting Goserver")
}

//ExitProg Function is for removing pidfile on exit
func ExitProg(msg string, pidfile string) {
	log.Printf("Exiting goserver %s", msg)
	os.Remove(pidfile)
	os.Exit(0)
}
