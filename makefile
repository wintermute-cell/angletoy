.PHONY: build build-win dev run clean

PROJECT="game-gorl"
BUILD_PATH="./build"

init:
	mkdir build
	mkdir runtime

build:
	mkdir -p $(BUILD_PATH)
	cp -r assets/* $(BUILD_PATH)
	go build -o $(BUILD_PATH)/$(PROJECT) -v cmd/game/main.go

build-debug:
	mkdir -p $(BUILD_PATH)
	cp -r assets/* $(BUILD_PATH)
	CGO_CFLAGS='-O0 -g' go build -a -v -gcflags="all=-N -l" -o $(BUILD_PATH)/$(PROJECT) cmd/game/main.go 

run:
	cd $(BUILD_PATH); ./$(PROJECT)

dev:
	@make build && make run || echo "build failed!"

clean:
	rm -r $(BUILD_PATH)/*
	mkdir -p $(BUILD_PATH)
	cp -r assets/* $(BUILD_PATH)

lint:
	nilaway main.go
