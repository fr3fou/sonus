MODULE_NAME=sonus
OUTPUT_DIR=sonus-out

ifeq ($(OS), Windows_NT)
	EXECUTABLE=$(MODULE_NAME).exe
	BUILD_FLAGS=-ldflags '-extldflags "-static" -H=windowsgui' .
else
	EXECUTABLE=$(MODULE_NAME)
	BUILD_FLAGS=.
endif

ifneq ("$(wildcard $(OUTPUT_DIR))","")
    OUTPUT_EXISTS = 1
endif

all: clean static-build bundle

static-build:
	go build $(BUILD_FLAGS)

bundle:
	mkdir $(OUTPUT_DIR)
	mv $(EXECUTABLE) $(OUTPUT_DIR)
	cp -r assets/ $(OUTPUT_DIR)

clean:
ifeq ($(OUTPUT_EXISTS), 1)
	rm -r $(OUTPUT_DIR)
endif