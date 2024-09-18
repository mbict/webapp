module webappv2/examples/cqs

go 1.23

replace github.com/mbict/webapp => ./../

require (
	github.com/mbict/go-commandbus/v2 v2.0.0-20221027201733-382e24555100
	github.com/mbict/go-querybus v0.0.0-20221027203450-0a5a5d72dc8d
	github.com/mbict/webapp v0.0.0
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
