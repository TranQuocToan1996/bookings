@REM This file for windows user
go build -o bookings.exe cmd\web\main.go cmd\web\middleware.go cmd\web\routes.go
.\bookings.exe