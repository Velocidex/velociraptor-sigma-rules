

all: linux

linux:
	go build -o velosigmac ./src/*.go

windows:
	GOOS=windows GOARCH=amd64 go build -o velosigmac.exe ./src/*.go

compile:
	./velosigmac compile --config ./config/windows_hayabusa_rules.yaml --output ./output/Velociraptor-Hayabusa-Rules.zip
	./velosigmac compile --config ./config/windows_hayabusa_event_monitoring.yaml --output ./output/Velociraptor-Hayabusa-Monitoring.zip
	./velosigmac compile --config ./config/ChopChopGo_rules.yaml --output ./output/Velociraptor-ChopChopGo-Rules.zip
