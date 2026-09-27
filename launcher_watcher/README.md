

1. shell: 執行 `./bin/main http --watcher`：有帶 `--watcher`，走主程序的邏輯。
2. go 主程序: `os` 執行 `./bin/main http`：不帶 `--watcher`，運行副程序的邏輯。
3. go 主程序: 一直盯著副程序（檢查可靠性），結束就重啟。
4. go 副程序: 開服務
5. 注意 一件事情， 主程序 副程序 是在 linux 角度， 實際應用開發中 ，副程序才是我們的主要商務邏輯

```mermaid
flowchart TD
    S(["Shell"]) -- "1. 執行<br/>./bin/main http --watcher<br/>有帶 --watcher" --> B["主程序<br/>WatchLauncher"]
    B -- "2. os/exec 啟動<br/>./bin/main http（不帶 --watcher）" --> C["副程序<br/>router.New().Run()"]
    C --> D["4. 開服務<br/>提供 /ping"]

    D -. "3. 結束就等 5 秒重啟<br/>(cmd.Wait() 收到訊號)" .-> B

    E(("Ctrl-C / SIGTERM")) -- "先送到主程序<br/>(signal.Notify)" --> B
    B -. "stopLauncher 轉給副程序<br/>(oCommand.Process.Signal)" .-> D

    N["5. 主／副是 linux 角度<br/>副程序才是主要商務邏輯"]

    style N fill:none,stroke-dasharray: 5 5
```