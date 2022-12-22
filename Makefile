ONESHELL:
SHELL = /bin/bash

gochecks:
	go mod tidy
	go fmt .
	# go vet fails with tiny go currently

gotests: gochecks
	pkgs_with_tests=("$$(find ./ -name '*_test.go' -printf "%h\n" | sort -ub)"); \
	scripts/runUnitTests.sh $$pkgs_with_tests

go-debug: gotests
	# May change file extension to .ll or .bin
	# See for more information: https://tinygo.org/docs/reference/usage/subcommands/
	scripts/build.sh -o debug.bc

flash: gotests
	scripts/launch_openocd.sh debug.bc

uf2: gotests
	scripts/build.sh -o release.uf2

terminal:
	if [ -a "/dev/ttyACM0" ]; then \
		minicom -D /dev/ttyACM0 -b 115200; \
	else \
		echo "No device /dev/ttyACM0 found"; \
	fi
