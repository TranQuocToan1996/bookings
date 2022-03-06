package config

import (
	"html/template"
	"log"

	"github.com/TranQuocToan1996/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

// AppConfig: Create this package to get accessed from any parts of this project
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	MailChan      chan models.MailData
}
