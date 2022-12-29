ONESHELL:
SHELL = /bin/bash

gochecks:
	# go vet fails with tiny go currently
	gofmt -w .
	go mod tidy

gotests: gochecks
	pkgs_with_tests=("$$(find ./ -name '*_test.go' -printf "%h\n" | sort -ub)"); \
	scripts/runUnitTests.sh $$pkgs_with_tests

# TODO add -no-debug and make a debug target
release: gotests
	if [[ ! -d build ]]; then \
		mkdir build ;\
	fi ;\
	build=$$(scripts/go_change_check.sh build/release); \
	if [ $$build == "true" ]; then \
		tinygo build -size full -target=pico -serial=uart -o build/release; \
	fi

flash: release
	scripts/launch_openocd.sh build/release

terminal:
	if [ -a "/dev/ttyACM0" ]; then \
		minicom -D /dev/ttyACM0 -b 115200; \
	else \
		echo "No device /dev/ttyACM0 found"; \
	fi
