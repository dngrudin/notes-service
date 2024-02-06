BINARY_NAME=notes-service

BUILD_DIR=build

BINARY_PATH=${BUILD_DIR}/${BINARY_NAME}

.PHONY: build
build:
	go build -o ${BINARY_PATH} cmd/${BINARY_NAME}/main.go

.PHONY: run
run: build
	${BINARY_PATH} -o ${BUILD_DIR}/

.PHONY: clean
clean:
	go clean
	rm ${BINARY_PATH}
