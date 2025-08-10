.PHONY: _ _chalk dependencies exec b
# make with no target defaults to "default".
_: run

EX_FILENAME := "discordgo-bot"
OUTPUT_PATH := "./.bin/"

_chalk:
	@command -v go >/dev/null 2>&1 || { \
		echo "golang is not installed."; \
		echo "install go here: https://golang.org/dl/"; \
		exit 1; \
	}
	@echo Checking for dependencies...
	@go get
	@go mod tidy
.env:
	@echo "\".env\" file does not exist in root! cannot continue."
	@echo "BOT_TOKEN = token_goes_here # https://discord.com/developers/applications" > .env
	@echo "created \".env\", update the file and run 'make' once more."
	@echo "                  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ "
	@echo 
	@echo 
	@echo 
	@exit 1
# get dependencies
dependencies: _chalk .env

exec: dependencies
	@echo "Executing \"main.go\"..."
	-go run main.go ${PARAMETERS} ||:
b: dependencies
	@echo "Building to \"${OUTPUT_PATH}${EX_FILENAME}${EXT}\"" 
	@go build -o ${OUTPUT_PATH}${EX_FILENAME}${EXT}

.PHONY: run run-no-terminal run-verbose
run:
	$(MAKE) exec TARGET=$@ PARAMETERS="${pmt}"
# note: i suggest using these only for checks and debugging
# 	 	do "go run main.go (parameters)" otherwise.
run-no-terminal:
	$(MAKE) exec TARGET=$@ PARAMETERS="--no-terminal"
run-verbose:
	$(MAKE) exec TARGET=$@ PARAMETERS="--verbose"

.PHONY: build build-exe
build:
	$(MAKE) b TARGET=$@ EXT=""
build-exe:
	$(MAKE) b TARGET=$@ EXT=".exe"