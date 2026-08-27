// 用官方 centrifuge-js SDK 改寫，寫法跟 landan-admin-refine 專案裡
// sample/centrifuge-1-context/centrifuge-subject.ts 的 CentrifugeSubject 同一套慣例對齊。
//
// 需要先安裝依賴：npm install centrifuge
// Node 22.6+ 原生支援 strip type，直接 `node client.ts` 就能跑，不用額外裝 ts-node/tsx。
import { Centrifuge, type Subscription } from "centrifuge";

const CENTRIFUGE_URL = "ws://localhost:4021/connection/websocket";

// cmd/centrifuge.go 的 OnConnect 是 server 端主動把每個連進來的 client
// Subscribe("all")（server-side subscription），不是等 client 自己送 subscribe
// 指令，這點跟 centrifuge-subject.ts 範例裡「client 自己 newSubscription 訂閱」
// 的用法不同。實測發現：即使是 server-side subscription，只要本地也用同一個
// channel 名稱開一個 Subscription 物件，publication 一樣會投遞進來，
// 所以介面用法還是可以維持跟參考範例一致 —— 差別只在於這個 Subscription
// 物件本身永遠不會走到 "subscribed" 狀態（會停在 "unsubscribed"），
// 真正代表「server 已經訂閱成功」的事件要看下面 centrifuge-level 的 "subscribed"。
const CHANNEL = "all";

let centrifuge: Centrifuge | null = null;
let subscription: Subscription | null = null;

function connect(sUrl: string) {
  const oCentrifuge = new Centrifuge(sUrl);

  oCentrifuge.on("connecting", (oCtx) => {
    console.log("[Centrifuge] 連線中:", oCtx.reason);
  });

  oCentrifuge.on("connected", (oCtx) => {
    console.log("[Centrifuge] 已連線, client id:", oCtx.client);
  });

  oCentrifuge.on("disconnected", (oCtx) => {
    console.log("[Centrifuge] 連線已關閉:", oCtx.reason);
  });

  oCentrifuge.on("error", (oCtx) => {
    console.log("[Centrifuge] 發生錯誤:", oCtx);
  });

  // server-side subscription 的確認事件是從 centrifuge-level 觸發，
  // 不是從下面的 Subscription 物件觸發。
  oCentrifuge.on("subscribed", (oCtx) => {
    console.log("[Centrifuge] server 端已訂閱 channel:", oCtx.channel);
  });

  const oSubscription = oCentrifuge.newSubscription(CHANNEL);

  oSubscription.on("publication", (oContext) => {
    console.log(`[push] channel=${CHANNEL} data=`, oContext.data);
  });

  oSubscription.on("error", (oCtx) => {
    console.log("[Centrifuge] 訂閱發生錯誤:", oCtx);
  });

  oSubscription.subscribe();

  centrifuge = oCentrifuge;
  subscription = oSubscription;

  oCentrifuge.connect();
}

connect(CENTRIFUGE_URL);
