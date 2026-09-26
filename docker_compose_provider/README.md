# docker_compose_provider

用 `docker compose` 啟動 **宿主機(Mac)上的程式**,而不是容器。
範例:`docker compose up` → 由宿主機的 ffmpeg 從攝影機截一張圖。

```sh
make build               # 編譯本機版本 -> bin/desktop
docker compose up        # 截圖,存到 ./runtime/desktop/shot-時間.jpg
docker compose down
```

## 依 OS / 架構選用不同的執行檔

`compose.yaml` 的 `type` 支援變數替換,預設是 `./bin/desktop`,設定 `PROVIDER` 就會換成別的:

```yaml
type: ${PROVIDER:-./bin/desktop}
```

```sh
make up      # 編譯「本機 OS/架構」的版本 bin/desktop-<os>-<arch>,再用它 docker compose up
make down
```

也可以手動指定,例如:

```sh
PROVIDER=./bin/desktop-linux-amd64 docker compose up
```

`make up` 是靠 `go env GOOS GOARCH` 判斷本機環境(Windows 會自動加 `.exe`)。
目前只在 macOS(darwin/arm64)上實測過;Linux、Windows 的路徑選擇沒有實測。

## 打包成各平台版本

純 Go(沒有 cgo),用 Go 內建的交叉編譯,不需要 Docker:

```sh
make all              # 全部
make build-mac        # bin/desktop-darwin-arm64、bin/desktop-darwin-amd64
make build-linux      # bin/desktop-linux-amd64、bin/desktop-linux-arm64
make build-windows    # bin/desktop-windows-amd64.exe、bin/desktop-windows-arm64.exe
make clean            # 刪除 bin/
```

**只有 macOS 版能真正截圖**:截圖用的是 ffmpeg 的 `avfoundation`(macOS 專屬),
其他平台編得出來、命令也能跑,但 `up` 會因為 ffmpeg 找不到裝置而失敗。

## 為什麼要這樣做

Docker Compose 本來只管理**容器**。容器跑在 Docker Desktop 的 Linux 虛擬機裡,
碰不到 Mac 的攝影機、視窗、選單列,所以桌面程式沒辦法放進容器。

要讓 Compose 去啟動一個「不是容器」的程式,唯一的辦法是 `provider`:
service 不寫 `image:`,改寫 `provider:`,Compose 就會呼叫你指定的可執行檔,由它自己決定要做什麼。

```yaml
services:
  screenshot:
    provider:
      type: ./bin/desktop          # 要呼叫的可執行檔(相對於 compose.yaml)
      options:                  # 會以 --device / --output 參數傳給程式
        device: default
        output: ./runtime/desktop/shot-{time}.jpg
```

## 規格:程式必須符合什麼

`provider` 是 Docker Compose 的概念,不是 Go 的概念。**只要要讓 Compose 啟動宿主機程式,
這個程式就必須符合下面的規格**,否則 Compose 呼叫時會失敗。

### 1. 要能被這樣呼叫

| 時機 | Compose 呼叫的指令 |
|---|---|
| 啟動前,詢問參數說明 | `bin/desktop compose metadata` |
| `docker compose up` | `bin/desktop compose --project-name=NAME up --device … --output … SERVICE` |
| `docker compose down` | `bin/desktop compose --project-name=NAME down SERVICE` |

- `up` 就是「把事情做起來」(這裡是截圖;若是長駐程式,就是啟動它)。
- `down` 就是「停掉」(這裡是一次性任務,什麼都不用停,但指令仍必須存在)。
- `compose.yaml` 的 `options` 每一項都會變成 `--key 值` 傳進來,程式要能接受(不認識的參數也不能報錯)。

### 2. 要用 JSON 回報

每行輸出一個 JSON 物件到標準輸出:

```json
{"type":"info","message":"screenshot: 截圖完成 /path/shot.jpg"}
{"type":"error","message":"screenshot: ffmpeg 失敗 ..."}
```

失敗時要以**非 0** 結束。`metadata` 則輸出一段描述參數的 JSON(見 `internal/logger/main.go`)。

### 3. 沒有的東西

Compose 只會呼叫 `up` / `down`。**掛掉自動重啟、日誌收集、`docker compose ps` 顯示狀態
都不會有**,要的話得由程式自己實作(或交給 launchd、pm2 這類工具)。

## 程式結構

跟一般 Go 專案一樣:根目錄 `main.go` 只負責啟動,命令放 `cmd/`,實際邏輯放 `internal/`。

```
main.go                          只做啟動:cmd.Execute()
cmd/                             命令樹(目錄是宿主機端命令的分組,命令名見右側)
├── root.go                      desktop
└── desktop/                     放宿主機端的命令(package desktop)
    ├── compose.go               desktop compose(帶 --project-name,組裝下面三個子命令)
    ├── metadata/metadata.go     desktop compose metadata
    ├── up/up.go                 desktop compose up   → 呼叫 screenshoter.Take 截圖
    └── down/down.go             desktop compose down
internal/
├── logger/main.go               輸出給 Compose 的 JSON 訊息、metadata(package logger)
└── screenshoter/main.go         用宿主機 ffmpeg 截一張圖(package screenshoter)
compose.yaml                     Compose 設定
```

每個命令一個目錄(一個 package),對外只公開 `Command`,由上一層 `AddCommand` 組裝,
沒有 `init()` 互相註冊的隱性順序。

## 注意事項

- `docker compose up` 會印出 `service "screenshot" has no container to start`,且結束碼為 1
  (專案裡沒有任何容器)。截圖本身是成功的,腳本判斷結果時要留意。
- 在互動終端機直接執行 `docker compose up`,預設的動畫進度畫面**看不到**「截圖完成」,
  只會顯示 `✔ screenshot Created`(截圖其實成功了,檔案在 `./runtime/desktop/`)。
  要看到訊息請加 `--progress=plain`,或用 `make up`(已內建)。
- 沒有實測過 `metadata` 是不是一定要有。
- `-framerate 30 -video_size 1280x720` 是這台 Mac 內建相機支援的模式,別台不一定支援。
- 第一次執行 macOS 可能詢問攝影機權限,權限算給「啟動它的程式」(終端機)。
  透過別的方式觸發時,權限算給誰沒有測過。
- 需要先安裝 ffmpeg(`brew install ffmpeg`)。
- 只在 Docker Compose v5.1.3、macOS 上驗證過。
