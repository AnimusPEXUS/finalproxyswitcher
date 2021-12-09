VERSION=0.7.1.5

all:

"$(GOROOT)/misc/wasm/wasm_exec.js":
    

wasm_exec.js: "$(GOROOT)/misc/wasm/wasm_exec.js"
		cp "$(GOROOT)/misc/wasm/wasm_exec.js" .

build: wasm_exec.js
		make -C ../background

dist: build
		mv ../background/background.wasm .
		-mkdir ../../dist
		7z -mx=9 -tzip a ../../dist/finalproxyswitcher-$(VERSION).xpi .
