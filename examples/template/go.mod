module webappv2/examples/default

go 1.23

replace github.com/mbict/webapp => ./../../

replace github.com/mbict/webapp/template => ./../../template

require (
	github.com/mbict/webapp v0.0.0-20230630153911-700d05ab545f
	github.com/mbict/webapp/template v0.0.0-00010101000000-000000000000
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.0 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
)
