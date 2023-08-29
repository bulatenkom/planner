# GNU Make
#
# File created on 20-08-2023
# (c) bulatenkom
#
# Build requirements:
# - Go 1.7.8
# - golang.org/x/tools/cmd/goimports (branch release-branch.go1.7)
# - Old-Go means old build approach with GOPATH


BUILD_TARGET = planner
GOIMPORTS = ${GOPATH}/bin/goimports

goimport:
	@ test -x ${GOIMPORTS} || echo "goimports tool is required for 'goimport' target!"
	set -u; ${GOIMPORTS} -w .

build ${BUILD_TARGET}:
	go build .

run: ${BUILD_TARGET}
	@ ./${BUILD_TARGET}

clean:
	@ rm ./${BUILD_TARGET}

help:
	@ printf "'make goimport, make' \t- fix imports and formatting in files tree\n"
	@ printf "'make build' \t\t- build application (producing executable file: '${BUILD_TARGET}')\n"
	@ printf "'make run' \t\t- run built application (to build application run 'make build')\n"
	@ printf "'make clean' \t\t- clean build artifacts\n"
	@ printf "'make help' \t\t- short description of make's targets\n"

.PHONY: goimport build run clean help
