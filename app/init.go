package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
)

//Mongoconn  Mongo Connection variable
var Mongoconn *mongo.Client
var stopflag bool
var mgerr error

func Init() {
	fmt.Println("Starting Goserver")
	dir, _ := os.Getwd()
	catchSignals(dir + "/goserver.pid")
	CreatePidFile(dir + "/goserver.pid")
	stopflag = false
	//mctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	Mongoconn, mgerr = ConnectMongo(context.TODO(), MONGOHOST)
	if mgerr != nil {
		log.Fatalf("Mongo Connection Failed :%v", mgerr)
	}
}

//catchSignals just catched all INT TERM etc signals and exits after removing pidfile
func catchSignals(PIDFILE string) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-signalChan
		log.Printf("GOT SIGNAL %v", sig)
		stopflag = true
		ExitProg("All Done", PIDFILE)
	}()

}
