package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration, which is accessible to every package that imports the <config>
type AppConfig struct {
	// UseCache defines if the application should get the template cache from the <AppConfig> (UseCache == true)
	// or if the application should create the template cache again (UseCache == false)
	UseCache bool

	// TemplateCache maps
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
