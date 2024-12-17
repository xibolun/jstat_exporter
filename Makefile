# Go编译器
GO := go

# 可执行文件名
BINARY_NAME := jstat_exporter

# 源代码文件
SOURCES := $(wildcard *.go)

.PHONY: all build run test clean

all: build

# 编译可执行文件
build:
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_NAME) $(SOURCES)

# 运行可执行文件
run:
	$(GO) run $(SOURCES)

# 清理生成的文件
clean:
	rm -f $(BINARY_NAME)