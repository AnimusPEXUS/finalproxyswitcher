VERSION=0.4

all:

build:
		make -C ../background

dist: build
		mv ../background/background.wasm .
		-mkdir ../../dist
		7z -mx=9 -tzip a ../../dist/finalproxyswitcher-$(VERSION).xpi .
		#tar -c . | xz -9 > ../../dist/finalproxyswitcher-$(VERSION).xz
