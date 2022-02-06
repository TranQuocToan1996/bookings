package models

// TemplateData hold data set from handlers to template
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} //other struct
	CSRFToken string
	Flash     string // Message send to user
	Warning   string
	Error     string
}
