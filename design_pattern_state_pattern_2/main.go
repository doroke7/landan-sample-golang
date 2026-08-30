package main

import "fmt"

// 狀態模式（State Pattern）—— 只能 next 的版本
//
// 核心概念：Context 只持有當前狀態，把請求委派給它，
// 由 State 自己決定下一個狀態，呼叫 Context.setState 換掉自己。
// 新增狀態只要多一個 State 實作，不用改 Context 內的 if/switch。
//
// 本範例：一個人的姿勢狀態，而且只提供一個動作 —— Next。
// 姿勢只能照固定順序往下走：站立 -> 坐下 -> 蹲著 -> 跳 -> 站立 -> ...
// 每個狀態自己決定「做完這個動作會變成什麼姿勢」。

// -----------------------------------------------------------------------------
// State 介面：每個姿勢只要能回應 Next
// -----------------------------------------------------------------------------

type StatePort interface {
	Name() string
	Next(oPerson *Person)
}

// -----------------------------------------------------------------------------
// Context：持有當前姿勢
// -----------------------------------------------------------------------------

type Person struct {
	state StatePort
}

func NewPerson() *Person {
	oPerson := &Person{
		state: &StandingState{},
	}

	fmt.Printf("初始姿勢：%s\n", oPerson.state.Name())

	return oPerson
}

func (oSelf *Person) setState(oState StatePort) {
	fmt.Printf("  %s -> %s\n", oSelf.state.Name(), oState.Name())
	oSelf.state = oState
}

func (oSelf *Person) Next() {
	oSelf.state.Next(oSelf)
}

// -----------------------------------------------------------------------------
// 具體狀態：站立
// -----------------------------------------------------------------------------

type StandingState struct {
}

func (oSelf *StandingState) Name() string {
	return "站立"
}

func (oSelf *StandingState) Next(oPerson *Person) {
	fmt.Println("[站立] 坐下")
	oPerson.setState(&SittingState{})
}

// -----------------------------------------------------------------------------
// 具體狀態：坐下
// -----------------------------------------------------------------------------

type SittingState struct {
}

func (oSelf *SittingState) Name() string {
	return "坐下"
}

func (oSelf *SittingState) Next(oPerson *Person) {
	fmt.Println("[坐下] 蹲下")
	oPerson.setState(&CrouchingState{})
}

// -----------------------------------------------------------------------------
// 具體狀態：蹲著
// -----------------------------------------------------------------------------

type CrouchingState struct {
}

func (oSelf *CrouchingState) Name() string {
	return "蹲著"
}

func (oSelf *CrouchingState) Next(oPerson *Person) {
	fmt.Println("[蹲著] 蓄力後往上跳")
	oPerson.setState(&JumpingState{})
}

// -----------------------------------------------------------------------------
// 具體狀態：跳（在空中）
// -----------------------------------------------------------------------------

type JumpingState struct {
}

func (oSelf *JumpingState) Name() string {
	return "跳"
}

func (oSelf *JumpingState) Next(oPerson *Person) {
	fmt.Println("[跳] 落地站好")
	oPerson.setState(&StandingState{})
}

// -----------------------------------------------------------------------------

func main() {
	oPerson := NewPerson()

	oPerson.Next() // 站立 -> 坐下
	oPerson.Next() // 坐下 -> 蹲著
	oPerson.Next() // 蹲著 -> 跳
	oPerson.Next() // 跳 -> 站立
	oPerson.Next() // 站立 -> 坐下
}
