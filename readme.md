# Booking and reservations

This is the repository for my bookings and reservations project.

- Go version: 1.18

- Dependencies:
  - [chi router](https://github.com/go-chi/chi): is a lightweight, idiomatic and composable router for building Go HTTP services. It's especially good at helping you write large REST API services that are kept maintainable as your project grows and changes. 
  - [SCS - session management](https://github.com/alexedwards/scs/v2)
  - [nosurf - CSRF protection](https://github.com/justinas/nosurf)
  - [Vanillajs Datepicker](https://mymth.github.io/vanillajs-datepicker/#/): Awsome Date Range Picker.
  - [Notie](https://github.com/jaredreich/notie): notification, input, and selection suite for Javascript, with no dependencies.
  - [Sweet Alert](https://github.com/t4t5/sweetalert): A beautiful replacement for JavaScript's "alert".
  - [GoValidator](https://github.com/asaskevich/govalidator): Validator for Go.
  - [Soda CLI](https://gobuffalo.io/en/docs/db/toolbox/): A small CLI toolbox to manage your database. It can help you to create a new database, drop existing ones, and so on.
  - [pgx](https://github.com/jackc/pgx): PostgreSQL Driver and Toolkit.
  - [MailHog](https://github.com/mailhog/MailHog): is an email testing tool.

- To do:
    - Install Soda CLI
        + go get github.com/gobuffalo/pop/...
        + go install github.com/gobuffalo/pop/soda
    - Create a Postgresql database
    - Change name "database - Copy.yml" -> "database.yml", after that correct the values inside.
    - Run command "soda migrate"

- To build and run the application, from the root level of the project, refer the file: windowsRun.sh and linuxBuild.sh
or refer this below command

    ```
    go build -o bookings ./cmd/web/ && ./bookings -dbname=yourDatabaseName -dbuser=yourDatabaseUserName
    ```

for full list of command use "./bookings -h"


- For the testing:
    - Run go test: 

        ```
        go test -v
        ```

    - Check your coverage with this command:

        ```
        go test -cover
        ```

    - Get your coverage in the browser with this command:

        ```
        go test -coverprofile=coverage.out && go tool cover -html=coverage.out
        ```
