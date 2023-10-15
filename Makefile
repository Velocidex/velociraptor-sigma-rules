compile:
	cd hayabusa/ && go run ../src/ compile --config ../config/logsources.yaml --output ../output/Velociraptor-Hayabusa-Rules.zip hayabusa/builtin/ hayabusa/sysmon/ sigma/builtin/ sigma/builtin/sysmon
