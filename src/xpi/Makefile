VERSION=0.7.0.2

all:

build:
		make -C ../background

dist: build
		mv ../background/background.wasm .
		-mkdir ../../dist
		7z -mx=9 -tzip a ../../dist/finalproxyswitcher-$(VERSION).xpi .
