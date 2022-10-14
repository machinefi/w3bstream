build_log_wasm:
	@cd log && tinygo build -o log.wasm -scheduler=none --no-debug -target=wasi log.go
	@echo built to log.wasm

build_json_wasm:
	@cd json && tinygo build -o parse_json.wasm -scheduler=none --no-debug -target=wasi parse_json.go
	@echo built to parse_json.wasm

build_word_count_wasm:
	@cd word_count && tinygo build -o word_count.wasm -scheduler=none --no-debug -target=wasi word_count.go
	@echo built to word_count.wasm

build_word_count_v2_wasm:
	@cd word_count_v2 && tinygo build -o word_count_v2.wasm -scheduler=none --no-debug -target=wasi word_count_v2.go
	@echo built to word_count_v2.wasm

all: build_log_wasm build_word_count_wasm build_word_count_v2_wasm