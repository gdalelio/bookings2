package models

// TemplateData holds data sent from the handler to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[int]int
	FloatMap  map[float32]float32
	Data      map[string]interface{} //any data type is called an interface
	CSRFToken string    //cross site request forge token 
	Flash     string
	Warning   string
	Error     string
}
