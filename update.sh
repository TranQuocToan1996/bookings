#!/bin/bash

git pull

soda migrate
go build -o bookingsLinux cmd/web/*.go

sudo supervisorctl stop book
sudo supervisorctl start book

