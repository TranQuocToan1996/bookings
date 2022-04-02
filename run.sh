#!/bin/bash
# This file for windows user with gitBash
# build ingore test file
# "host=localhost port=5432 dbname=bookings user=postgres password=postgres"
# env GOOS=linux go build -o bookings cmd/web/*.go
go build -o bookings cmd/web/*.go
./bookings -dbname=bookings -dbuser=postgres -cache=false -production=true -dbpass=postgres