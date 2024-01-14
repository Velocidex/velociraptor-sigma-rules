build_indexes: hayabusa_index chopchopgo_index

hayabusa_index:
	python3 ./scripts/curated_rules_reference_index.py --github_link https://github.com/Yamato-Security/hayabusa-rules/blob/main/ --root ../hayabusa/ ../hayabusa/hayabusa/builtin/ ../hayabusa/hayabusa/sysmon/ ../hayabusa/sigma/builtin ../hayabusa/sigma/builtin/sysmon static/index/hayabusa_index.json

chopchopgo_index:
	python3 ./scripts/curated_rules_reference_index.py --github_link https://github.com/M00NLIG7/ChopChopGo/blob/master/ --root ../ChopChopGo/ ../ChopChopGo/rules/linux static/index/chopchopgo_index.json
