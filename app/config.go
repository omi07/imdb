package app

const (
	//MAXCONN Maxconnection
	MAXCONN = 10

	//HTTPPORT The port where go server is running, type string
	HTTPPORT = ":9080"

	// DEFAULTMSG is the message sent to all wrong urls
	DEFAULTMSG = "Custom http"

	//PIDFILE location where the pid file exists and also for lock
	PIDFILE = "$HOME/goserver.pid"

	//MAXRETRY Maximum retry for Mongo Connection
	MAXRETRY = 3

	//MONGOHOST Mongo host ....
	MONGOHOST = "mongodb://localhost:27017"

	//MAXCONNPOOL Maximum number of connections
	MAXCONNPOOL = 5

	//Dbname ...
	Dbname = "imdb"

	//UserCollection  ....
	UserCollection = "userdetails"

	//MovieCollection
	MovieCollection = "movies"
)
