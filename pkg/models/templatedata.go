package models

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // For other data structures, use an interface
	CSRFToken string                 // Cross site request forgery token. This token is called in <search-availability.page.tmpl>
	Flash     string                 // Flash message to the end-user
	Warning   string                 // Warning message to the end-user
	Error     string                 // Error message to the end-user
}
