all: linux

linux:
	go build -o velosigmac ./src/*.go

windows:
	GOOS=windows go build -o velosigmac.exe .\src\

compile: compileThirdParty

compileThirdParty:  compileHayabusa compileHayabusaMonitoring compileChopChopGo compileWindowsETW compileHayabusaMonitoring compileLinuxEBPF

compileWindowsBaseDebug:
	dlv debug ./src/ -- compile --config ./config/windows_base.yaml --output ./output/Windows-Sigma-Base.zip --yaml ./output/Windows.Sigma.Base.yaml --docs ./docs/content/docs/models/windows_base/_index.md

compileWindowsBase:
	./velosigmac compile --config ./config/windows_base.yaml --output ./output/Windows-Sigma-Base.zip --yaml ./output/Windows.Sigma.Base.yaml --docs ./docs/content/docs/models/windows_base/_index.md
	./velosigmac compile --config ./config/windows_base_test.yaml --yaml ./output/Windows.Sigma.Base.CaptureTestSet.yaml

compileWindowsBaseEvents:
	./velosigmac compile --config ./config/windows_base_events.yaml --output ./output/Windows-Sigma-BaseEvents.zip --yaml ./output/Windows.Sigma.BaseEvents.yaml  --docs ./docs/content/docs/models/windows_base_events/_index.md
	./velosigmac compile --config ./config/windows_base_event_test.yaml --yaml ./output/Windows.Sigma.BaseEvents.CaptureTestSet.yaml

compileWindowsBaseETW:
	./velosigmac compile --config ./config/windows_etw_base.yaml --output ./output/Windows-Sigma-ETWBase.zip --yaml ./output/Windows.Sigma.ETWBase.yaml  --docs ./docs/content/docs/models/windows_etw_base/_index.md

compileLinuxEBPF: compileLinuxEBPFBase

compileLinuxEBPFBase:
	./velosigmac compile --config ./config/linux_ebpf_base.yaml --output ./output/Linux-Sigma-EBPFBase.zip --yaml ./output/Linux.Sigma.EBPFBase.yaml  --docs ./docs/content/docs/models/linux_ebpf_base/_index.md
	./velosigmac compile --config ./config/linux_ebpf_base_test.yaml --yaml ./output/Linux.Sigma.EBPFBase.CaptureTestSet.yaml

compileWindowsETW: compileWindowsBaseETW
	./velosigmac compile --config ./config/windows_etw_fibratus.yaml --output ./output/Windows-Sigma-Fibratus.zip --yaml ./output/Windows.Sigma.Fibratus.yaml

compileHayabusa: compileWindowsBase
	./velosigmac compile --config ./config/windows_hayabusa_rules.yaml --output ./output/Velociraptor-Hayabusa-Rules.zip --yaml ./output/Velociraptor.Hayabusa.Rules.yaml --rejects rejected/windows_hayabusa_rejects.json --ignore_previous_rejects

debugHayabusa:
	dlv debug ./src -- compile --config ./config/windows_hayabusa_rules.yaml --output ./output/Velociraptor-Hayabusa-Rules.zip --yaml ./output/Velociraptor-Hayabusa-Rules.yaml

compileHayabusaMonitoring: compileWindowsBaseEvents
	./velosigmac compile --config ./config/windows_hayabusa_event_monitoring.yaml --output ./output/Velociraptor-Hayabusa-Monitoring.zip --yaml ./output/Velociraptor-Hayabusa-Monitoring.yaml --rejects rejected/windows_hayabusa_rejects.json --ignore_previous_rejects

compileChopChopGo:
	./velosigmac compile --config ./config/ChopChopGo_rules.yaml --output ./output/Velociraptor-ChopChopGo-Rules.zip --yaml ./output/Velociraptor-ChopChopGo-Rules.yaml --rejects rejected/ChopChopGo_rules_rejects.json --ignore_previous_rejects

package:
	zip -v ./output/Velociraptor.Sigma.Artifacts.zip  ./output/*.yaml

test: compile
	go test -v ./...

golden:
	./tests/velociraptor -v --definitions ./output/ golden ./tests/testcases/ --config tests/golden.config.yaml --env testDir=`pwd`/tests/  --filter=${GOLDEN}

profile:
	./velosigmac gen_profiles --output ./output/profiles.json ./config/*base*.yaml
