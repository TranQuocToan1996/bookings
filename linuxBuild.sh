#!/bin/bash
# This file for linux
# build ingore test file
# "host=localhost port=5432 dbname=bookings user=postgres password=postgres"
env GOOS=linux go build -o bookingsLinux cmd/web/*.go
# ./bookingsLinux -dbname=bookings -dbuser=postgres -cache=false -production=false -dbpass=postgres