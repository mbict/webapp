module github.com/mbict/webapp/binder

go 1.23

replace github.com/mbict/webapp => ./../

require (
	github.com/mbict/webapp v0.0.0-20230630153911-700d05ab545f
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0
)

require github.com/stretchr/testify v1.8.4 // indirect
