

all: linux

linux:
	go build -o velosigmac ./src/*.go

windows:
	go build -o velosigmac.exe .\src\

compile:
	./velosigmac compile --config ./config/windows_hayabusa_rules.yaml --output ./output/Velociraptor-Hayabusa-Rules.zip --yaml ./output/Velociraptor-Hayabusa-Rules.yaml
	./velosigmac compile --config ./config/windows_hayabusa_event_monitoring.yaml --output ./output/Velociraptor-Hayabusa-Monitoring.zip --yaml ./output/Velociraptor-Hayabusa-Monitoring.yaml
	./velosigmac compile --config ./config/ChopChopGo_rules.yaml --output ./output/Velociraptor-ChopChopGo-Rules.zip --yaml ./output/Velociraptor-ChopChopGo-Rules.yaml

test: compile
	VELOCIRAPTOR_CONFIG= ../velociraptor/output/velociraptor-v0.7.1-rc1-linux-amd64 --definitions output/ artifacts list -v |grep 'Haya\|Chop'
