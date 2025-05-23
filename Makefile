all: linux

linux:
	go build -o velosigmac ./src/*.go

windows:
	GOOS=windows go build -o velosigmac.exe .\src\

compile: compileThirdParty

compileThirdParty:  compileHayabusa compileHayabusaMonitoring compileLinuxTriage compileWindowsETW compileHayabusaMonitoring compileLinuxEBPF compileWindowsVQL

compileWindowsBaseDebug:
	dlv debug ./src/ -- compile --config ./config/windows_base.yaml --output ./output/Windows-Sigma-Base.zip --yaml ./output/Windows.Sigma.Base.yaml --docs ./docs/content/docs/models/windows_base/_index.md

compileWindowsBase:
	./velosigmac compile --config ./config/windows_base.yaml --output ./output/Windows-Sigma-Base.zip --yaml ./output/Windows.Sigma.Base.yaml --docs ./docs/content/docs/models/windows_base/_index.md
	./velosigmac compile --config ./config/windows_base_test.yaml --yaml ./output/Windows.Sigma.Base.CaptureTestSet.yaml
	./velosigmac compile --config ./config/windows_base_test_replay.yaml --yaml ./output/Windows.Sigma.Base.ReplayTestSet.yaml


compileLinuxBase:
	./velosigmac compile --config ./config/linux_base.yaml --output ./output/Linux-Sigma-Base.zip --yaml ./output/Linux.Sigma.Base.yaml --docs ./docs/content/docs/models/linux_base/_index.md
	./velosigmac compile --config ./config/linux_base_test.yaml --yaml ./output/Linux.Sigma.Base.CaptureTestSet.yaml
	./velosigmac compile --config ./config/linux_base_test_replay.yaml --yaml ./output/Linux.Sigma.Base.ReplayTestSet.yaml

compileLinuxTriage: compileLinuxBase
	./velosigmac compile --config ./config/linux_sigma_triage.yaml --output ./output/Linux-Sigma-Triage.zip --yaml ./output/Linux.Sigma.Triage.yaml --rule_dir ./docs/content/docs/artifacts/Linux.Sigma.Triage/ --docs ./docs/content/docs/artifacts/Linux.Sigma.Triage/_index.md


compileWindowsBaseEvents:
	./velosigmac compile --config ./config/windows_base_events.yaml --output ./output/Windows-Sigma-BaseEvents.zip --yaml ./output/Windows.Sigma.BaseEvents.yaml  --docs ./docs/content/docs/models/windows_base_events/_index.md
	./velosigmac compile --config ./config/windows_base_event_test.yaml --yaml ./output/Windows.Sigma.BaseEvents.CaptureTestSet.yaml
	./velosigmac compile --config ./config/windows_base_event_test_replay.yaml --yaml ./output/Windows.Sigma.BaseEvents.ReplayTestSet.yaml


compileWindowsBaseETW:
	./velosigmac compile --config ./config/windows_etw_base.yaml --output ./output/Windows-Sigma-ETWBase.zip --yaml ./output/Windows.Sigma.ETWBase.yaml  --docs ./docs/content/docs/models/windows_etw_base/_index.md
	./velosigmac compile --config ./config/windows_etw_base_test.yaml --yaml ./output/Windows.Sigma.ETWBase.CaptureTestSet.yaml
	./velosigmac compile --config ./config/windows_etw_base_test_replay.yaml --yaml ./output/Windows.Sigma.ETWBase.ReplayTestSet.yaml

compileLinuxEBPF: compileLinuxEBPFBase
	./velosigmac compile --config ./config/linux_ebpf_monitoring.yaml --output ./output/Linux-Sigma-EBPF.zip --yaml ./output/Linux.Sigma.EBPF.yaml --rule_dir ./docs/content/docs/artifacts/Linux.EBPF.Monitoring/ --docs ./docs/content/docs/artifacts/Linux.EBPF.Monitoring/_index.md

compileLinuxEBPFBase:
	./velosigmac compile --config ./config/linux_ebpf_base.yaml --output ./output/Linux-Sigma-EBPFBase.zip --yaml ./output/Linux.Sigma.EBPFBase.yaml  --docs ./docs/content/docs/models/linux_ebpf_base/_index.md
	./velosigmac compile --config ./config/linux_ebpf_base_test.yaml --yaml ./output/Linux.Sigma.EBPFBase.CaptureTestSet.yaml
	./velosigmac compile --config ./config/linux_ebpf_base_test_replay.yaml --yaml ./output/Linux.Sigma.EBPFBase.ReplayTestSet.yaml

compileWindowsETW: compileWindowsBaseETW
	./velosigmac compile --config ./config/windows_etw_monitoring.yaml --output ./output/Windows-ETW-Monitoring.zip --yaml ./output/Windows.ETW.Monitoring.yaml --rule_dir ./docs/content/docs/artifacts/Windows.ETW.Monitoring/ --docs ./docs/content/docs/artifacts/Windows.ETW.Monitoring/_index.md

compileWindowsBaseVQL:
	./velosigmac compile --config ./config/windows_base_vql.yaml --output ./output/Windows-Sigma-BaseVQL.zip --yaml ./output/Windows.Sigma.BaseVQL.yaml  --docs ./docs/content/docs/models/windows_base_vql/_index.md

compileWindowsVQL: compileWindowsBaseVQL
	./velosigmac compile --config ./config/windows_vql_triage.yaml --output ./output/Windows-Sigma-Triage.zip --yaml ./output/Windows.Sigma.Triage.yaml --rule_dir ./docs/content/docs/artifacts/Windows.Sigma.Triage/ --docs ./docs/content/docs/artifacts/Windows.Sigma.Triage/_index.md

compileHayabusa: compileWindowsBase
	./velosigmac compile --config ./config/windows_hayabusa_rules.yaml --output ./output/Velociraptor-Hayabusa-Rules.zip --yaml ./output/Velociraptor.Hayabusa.Rules.yaml --rejects rejected/windows_hayabusa_rejects.json --ignore_previous_rejects --rule_dir ./docs/content/docs/artifacts/Windows.Hayabusa.Rules/ --docs ./docs/content/docs/artifacts/Windows.Hayabusa.Rules/_index.md

debugHayabusa:
	dlv debug ./src -- compile --config ./config/windows_hayabusa_rules.yaml --output ./output/Velociraptor-Hayabusa-Rules.zip --yaml ./output/Velociraptor-Hayabusa-Rules.yaml

compileHayabusaMonitoring: compileWindowsBaseEvents
	./velosigmac compile --config ./config/windows_hayabusa_event_monitoring.yaml --output ./output/Velociraptor-Hayabusa-Monitoring.zip --yaml ./output/Velociraptor-Hayabusa-Monitoring.yaml --rejects rejected/windows_hayabusa_rejects.json --ignore_previous_rejects --rule_dir ./docs/content/docs/artifacts/Windows.Hayabusa.Monitoring/ --docs ./docs/content/docs/artifacts/Windows.Hayabusa.Monitoring/_index.md

package:
	zip -v ./output/Velociraptor.Sigma.Artifacts.zip  ./output/*.yaml

test: compile
	go test -v ./...

golden:
	./tests/velociraptor -v --definitions ./output/ golden ./tests/testcases/ --config tests/golden.config.yaml --env testDir=`pwd`/tests/  --filter=${GOLDEN}

profile:
	./velosigmac gen_profiles --output ./output/profiles.json ./config/*base*.yaml

serve:
	cd docs && hugo serve --port 1314 --bind 0.0.0.0
