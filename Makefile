

all: linux

linux:
	go build -o velosigmac ./src/*.go

windows:
	GOOS=windows GOARCH=amd64 go build -o velosigmac.exe ./src/*.go

compile:
	cd hayabusa/ && go run ../src/ compile --config ../config/windows_hayabusa_rules.yaml --output ../output/Velociraptor-Hayabusa-Rules.zip
