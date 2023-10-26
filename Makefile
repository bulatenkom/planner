# GNU Make

BUILD_TARGET = planner

build ${BUILD_TARGET}:
	go build .

run: ${BUILD_TARGET}
	@ ./${BUILD_TARGET}

clean:
	@ rm ./${BUILD_TARGET}

dev:
	go build . && ./${BUILD_TARGET}

help:
	@ printf "'make build' \t\t- build application (producing executable file: '${BUILD_TARGET}')\n"
	@ printf "'make run' \t\t- run built application (to build application run 'make build')\n"
	@ printf "'make clean' \t\t- clean build artifacts\n"
	@ printf "'make dev' \t\t- build and run app in dev mode\n"
	@ printf "'make help' \t\t- short description of make's targets\n"

.PHONY: build run clean dev help
