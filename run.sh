#!/bin/bash
# This file for linux and MAC
# build ingore test file
go build -o bookings.exe cmd/web/*.go
./bookings.exe