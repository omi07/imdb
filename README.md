# IMDB API

IMDB API is a simple Golang Rest API platform Movie Library.

## Installation

```bash
     ./install.sh
```

This script installs Golang latest version set GOPATH & download all required packages to run this project



## Features


Register
--

- Register a user into the system. Requires username  which should be unique


Login
-- 

- Login user into system using username & password return a token for that user.


Add a Movie
--

- The only user with the admin role can add movies into the system


Rate or Comment a Movie
--
- All Identified user can  rate or comment on any movie


Search a Movie
--
- Search a Movie into the system by movie


Get all movie
--
- Gives all movies present into the system




##To Do
--- 
- User authorized based on token can be managed using a session.

- A rating can be stored using a timestamp to show the recent rated movie

- Incremental Userid can be managed by using caching technology like Redis.



 
