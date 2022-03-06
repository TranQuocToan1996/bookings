#!/bin/bash
# This file for linux and MAC
# build ingore test file
# "host=localhost port=5432 dbname=bookings user=postgres password=postgres"
go build -o bookings.exe cmd/web/*.go
./bookings.exe -dbname=bookings -dbuser=postgres -cache=false -production=false -dbpass=postgres