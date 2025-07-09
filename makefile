# BUILD_NAME := bot
# EXPORT_PATH := ./.bin/

# ${EXPORT_PATH}${BUILD_NAME}:
# 	mkdir -p ${EXPORT_PATH}
# 	go build -o ${EXPORT_PATH}${BUILD_NAME} -v

.PHONY: run no-terminal
run:
	@echo "Executing \"main.go\"..."
	-go run main.go ||:

no-terminal:
	@echo "Executing \"main.go\" with terminal disabled..."
	-go run main.go --no-terminal ||:

dependencies:
	@echo "Getting dependencies..."
	go mod tidy
