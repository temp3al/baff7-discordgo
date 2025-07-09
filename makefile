# BUILD_NAME := bot
# EXPORT_PATH := ./.bin/

# ${EXPORT_PATH}${BUILD_NAME}:
# 	mkdir -p ${EXPORT_PATH}
# 	go build -o ${EXPORT_PATH}${BUILD_NAME} -v

.PHONY: run
run:
	@echo "Running executable..."
	@go run main.go

dep:
	@echo "Getting dependencies..."
	@go mod tidy
