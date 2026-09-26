# 純 Go(沒有 cgo),直接用 GOOS/GOARCH 交叉編譯,不需要 Docker。
# 產物都放在 bin/,檔名: desktop-<os>-<arch>

GOOS     := $(shell go env GOOS)
GOARCH   := $(shell go env GOARCH)
EXT      := $(if $(filter windows,$(GOOS)),.exe,)
PROVIDER := ./bin/desktop-$(GOOS)-$(GOARCH)$(EXT)

LDFLAGS := -s -w
BUILD   := CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)"

# 預設目標
.PHONY: help
help:
	@echo "make build          編譯本機版本 -> bin/desktop (docker compose 用這個)"
	@echo "make build-mac      macOS   arm64 + amd64"
	@echo "make build-linux    Linux   amd64 + arm64"
	@echo "make build-windows  Windows amd64 + arm64"
	@echo "make all            以上全部"
	@echo "make up             編譯本機 OS/架構的版本,再 docker compose up"
	@echo "make down           docker compose down"
	@echo "make clean          刪除 bin/"

.PHONY: build
build:
	$(BUILD) -o bin/desktop .

.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=arm64 $(BUILD) -o bin/desktop-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 $(BUILD) -o bin/desktop-darwin-amd64 .

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/desktop-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(BUILD) -o bin/desktop-linux-arm64 .

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/desktop-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 $(BUILD) -o bin/desktop-windows-arm64.exe .

# 依目前的 OS/架構挑執行檔:編譯 bin/desktop-<os>-<arch>,再用 PROVIDER 變數告訴 compose.yaml
.PHONY: up
up:
	$(BUILD) -o $(PROVIDER) .
	@# 專案裡沒有容器時 docker compose up 一定回傳 1(截圖本身是成功的),所以前面加 - 忽略結束碼;
	@# 真的失敗時,錯誤訊息仍會顯示在輸出裡。
	@# --progress=plain:預設的動畫進度畫面會吃掉 provider 回報的訊息(看不到「截圖完成」)。
	-PROVIDER=$(PROVIDER) docker compose --progress=plain up

.PHONY: down
down:
	PROVIDER=$(PROVIDER) docker compose --progress=plain down

.PHONY: all
all: build-mac build-linux build-windows

.PHONY: clean
clean:
	rm -rf bin
