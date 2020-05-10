#!/bin/sh

set -x 
wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin


systemctl daemon-reload
go get "github.com/gorilla/mux" "github.com/dgrijalva/jwt-go"
go get "go.mongodb.org/mongo-driver/mongo"  "golang.org/x/crypto/bcrypt"
set +x

