module github.com/Velocidex/velociraptor-sigma-rules

go 1.21

toolchain go1.23.2

require (
	github.com/Velocidex/ordereddict v0.0.0-20230909174157-2aa49cc5d11d
	github.com/Velocidex/sigma-go v0.0.0-20241113062227-c1c5ea4b5250
	github.com/Velocidex/yaml/v2 v2.2.8
	github.com/alecthomas/kingpin/v2 v2.3.2
	github.com/davecgh/go-spew v1.1.1
	github.com/sebdah/goldie/v2 v2.5.3
	github.com/stretchr/testify v1.8.4
	golang.org/x/exp v0.0.0-20240213143201-ec583247a57a
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Velocidex/json v0.0.0-20220224052537-92f3c0326e5a // indirect
	github.com/alecthomas/participle v0.7.1 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

// replace github.com/Velocidex/sigma-go => ../sigma-go
