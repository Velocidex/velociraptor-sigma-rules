

all: linux

linux:
	go build -o velosigmac ./src/*.go

windows:
	GOOS=windows go build -o velosigmac.exe .\src\

compile: compileThirdParty compileCurated

compileThirdParty:  compileHayabusa compileHayabusaMonitoring compileChopChopGo

compileHayabusa:
	./velosigmac compile --config ./config/windows_hayabusa_rules.yaml --output ./output/Velociraptor-Hayabusa-Rules.zip --yaml ./output/Velociraptor-Hayabusa-Rules.yaml --rejects rejected/windows_hayabusa_rejects.json --ignore_previous_rejects

debugHayabusa:
	dlv debug ./src -- compile --config ./config/windows_hayabusa_rules.yaml --output ./output/Velociraptor-Hayabusa-Rules.zip --yaml ./output/Velociraptor-Hayabusa-Rules.yaml

compileHayabusaMonitoring:
	./velosigmac compile --config ./config/windows_hayabusa_event_monitoring.yaml --output ./output/Velociraptor-Hayabusa-Monitoring.zip --yaml ./output/Velociraptor-Hayabusa-Monitoring.yaml --rejects rejected/windows_hayabusa_rejects.json --ignore_previous_rejects

compileChopChopGo:
	./velosigmac compile --config ./config/ChopChopGo_rules.yaml --output ./output/Velociraptor-ChopChopGo-Rules.zip --yaml ./output/Velociraptor-ChopChopGo-Rules.yaml --rejects rejected/ChopChopGo_rules_rejects.json --ignore_previous_rejects

compileCurated:  compilePostProcess compileWindowsRules

compilePostProcess:
	./velosigmac compile --config ./config/velociraptor_post_process.yaml --output ./output/Velociraptor-Post-Process.zip --yaml ./output/Velociraptor-Post-Process.yaml

compileWindowsRules:
	./velosigmac compile --config ./config/velociraptor_windows_rules.yaml --output ./output/Velociraptor-Windows-Rules.zip --yaml ./output/Velociraptor-Windows-Rules.yaml

test: compile
	go test -v ./...

golden:
	./tests/velociraptor -v --definitions ./output/ golden ./tests/testcases/ --config tests/golden.config.yaml --env testDir=`pwd`/tests/  --filter=${GOLDEN}
