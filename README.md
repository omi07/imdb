# imdb
Imdb

# IMDB API

IMDB API is a simple Golang Rest api platform Movie Library.

## Installation

```bash
 	./install.sh
```

This script install Golang latest version set GOPATH & download all required packages to run this project 



## Features 

 Register
Register a user into the system. Requires username  which should be unique 

Login 
Login user into system using username & password return a token for that user.


Add a Movie
Only user with admin role can add movies into the system 

 
Rate or Comment a Movie
All Identified user can  rate or comment any movie 


Search a Movie
Search a Movie into system by passing its movieid

 
Get all movie
Gives all movies present into the system




##TO DO 

User authorised on based on token can be managed using session.
Rating can be stored using time stamp to show recent rated movie 
Incremental Userid can be managed by using caching technology like Redis.

##Attachment
Database json 
Postman Collection 


 
