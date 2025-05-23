GCP_ID:= GOOGLE_CLOUD_PROJECT_ID
BIN_FOLDER := ./build
BIN_NAME := BIN_NAME 
BIN_PATH := $(BIN_FOLDER)/$(BIN_NAME)
MAIN:= MAIN_FILE_DOT_GO
MAIN_PATH:= ./$(MAIN)

clean:
	rm -rf $(BIN_FOLDER)

# Sets the GCP project ID and quota
project/setup:
	gcloud config set project $(GCP_ID) 
	gcloud auth application-default set-quota-project  $(GCP_ID)

# Builds the project
project/build:
	go build -o $(BIN_PATH) $(MAIN_PATH)

# Builds and run the executable
project/run: project/build
	./$(BIN_PATH)

# Deploys the project to Google App Engine
project/deploy:
	gcloud app deploy

# Run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "make project/build" --build.bin "$(BIN_PATH)" --build.delay "500" \
	--build.exclude_dir "" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# Run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	npx @tailwindcss/cli -i ./static/tailwind.css -o ./static/dist/style.css --minify --watch

# Generates the style.css file
build/tailwind:
	tailwindcss -i ./static/tailwind.css -o ./static/dist/style.css --minify

# Run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ: clean
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

# start all watch processes in parallel.
live:
	make -j5 live/templ live/server live/tailwind
