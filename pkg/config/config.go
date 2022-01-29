package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// Create this package to get access from any part of this project

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
}
