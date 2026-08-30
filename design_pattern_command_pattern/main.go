package main

import "fmt"

// 命令模式（Command Pattern）
//
// 核心概念：把「一個請求」封裝成物件，讓請求可以被傳遞與觸發。
//   - Command  ：介面，宣告 Execute
//   - Receiver ：真正幹活的物件（本例的 Light）
//   - Invoker  ：持有命令並觸發它，不知道命令細節（本例的 Button）
//   - Client   ：組裝命令與接收者（本例的 main）

// -----------------------------------------------------------------------------
// Command 介面
// -----------------------------------------------------------------------------

type Command interface {
	Execute()
}

// -----------------------------------------------------------------------------
// Receiver：真正執行動作的物件
// -----------------------------------------------------------------------------

type Light struct {
}

func (oSelf *Light) On() {
	fmt.Println("[燈] 開")
}

func (oSelf *Light) Off() {
	fmt.Println("[燈] 關")
}

// -----------------------------------------------------------------------------
// 具體命令：把「動作 + 接收者」綁在一起
// -----------------------------------------------------------------------------

type LightOnCommand struct {
	light *Light
}

func (oSelf *LightOnCommand) Execute() {
	oSelf.light.On()
}

type LightOffCommand struct {
	light *Light
}

func (oSelf *LightOffCommand) Execute() {
	oSelf.light.Off()
}

// -----------------------------------------------------------------------------
// Invoker：持有命令、觸發命令
// -----------------------------------------------------------------------------

type Button struct {
	command Command
}

func (oSelf *Button) Press() {
	oSelf.command.Execute()
}

// -----------------------------------------------------------------------------

func main() {
	oLight := &Light{}

	oOnButton := &Button{
		command: &LightOnCommand{
			light: oLight,
		},
	}

	oOffButton := &Button{
		command: &LightOffCommand{
			light: oLight,
		},
	}

	oOnButton.Press()
	oOffButton.Press()

	/*
			這樣看起來 strategy 模式跟 command 模式非常類似，他都是透過注入解耦，
		不一樣的是 strategy 關心 注入的程序邏輯， command 關心注入的東西的規格，

			**/
}
