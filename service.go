package webappv2

type Service interface {
	Binder() Binder
	Validator() Validator
	Router() Router
	Routes() Routes
	Renderer() Renderer
	Logger() Logger

	JsonEncoder() JSONEncoding
	XmlEncoder() XMLEncoding
}
