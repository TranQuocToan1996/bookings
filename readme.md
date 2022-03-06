# Booking and reservations

This is the repository for my bookings and reservations project.

- Go version: 1.17
- Use:

  - [chi router](https://github.com/go-chi/chi)
  - [SCS - session management](https://github.com/alexedwards/scs/v2)
  - [nosurf - CSRF protection](https://github.com/justinas/nosurf)
  - [Vanillajs Datepicker](https://mymth.github.io/vanillajs-datepicker/#/): Awsome Date Range Picker.
  - [Notie](https://github.com/jaredreich/notie): notification, input, and selection suite for Javascript, with no dependencies.
  - [Sweet Alert](https://github.com/t4t5/sweetalert): A beautiful replacement for JavaScript's "alert".
  - [GoValidator](https://github.com/asaskevich/govalidator): Validator for Go.
  - [Soda CLI](https://gobuffalo.io/en/docs/db/toolbox/): A small CLI toolbox to manage your database. It can help you to create a new database, drop existing ones, and so on.
  - [pgx](https://github.com/jackc/pgx): PostgreSQL Driver and Toolkit.
  - [MailHog](https://github.com/mailhog/MailHog): is an email testing tool.
  - [Foundation for Emails 2](https://get.foundation/emails.html): Quickly create responsive HTML emails.
https://github.com/BootstrapDash/RoyalUI-Free-Bootstrap-Admin-Template
 https://github.com/fiduswriter/Simple-DataTables
 https://github.com/joho/godotenv


- For the testing:
 
Things you can do:
    - Run go test: 

        ```go test -v```
        
    - Check your coverage with this command:

        ```go test -cover```


    - Get your coverage in the browser with this command:
    
        ```go test -coverprofile=coverage.out && go tool cover -html=coverage.out```
