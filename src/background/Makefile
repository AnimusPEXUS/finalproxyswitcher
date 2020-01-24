all:
		# gopherjs build -v -o background.js .
		-rm background.wasm
		GOOS=js GOARCH=wasm go build -v -o background.wasm
		#gzip -9 background.wasm
		#mv background.wasm.gz background.wasm
