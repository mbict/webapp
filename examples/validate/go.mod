module webappv2/examples/default

go 1.23

replace github.com/mbict/webapp => ./../../

replace github.com/mbict/webapp/validator => ./../../validator

require (
	github.com/mbict/webapp v0.0.0-20240918184906-9581e3fee779
	github.com/mbict/webapp/validator v0.0.0-00010101000000-000000000000
)

require (
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
