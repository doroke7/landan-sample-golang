package main

import "fmt"

// 狀態模式（State Pattern）
//
// 核心概念：Context 只持有當前狀態，把請求委派給它，
// 由 State 自己決定下一個狀態，呼叫 Context.setState 換掉自己。
// 新增狀態只要多一個 State 實作，不用改 Context 內的 if/switch。
//
// 本範例：一個人的姿勢狀態 —— 站立 / 坐下 / 蹲著 / 跳。
// 每個狀態自己決定「從這裡能做什麼、會變成什麼姿勢」。

// -----------------------------------------------------------------------------
// State 介面：每個姿勢都要能回應這四個動作
// -----------------------------------------------------------------------------

type StatePort interface {
	Name() string
	Stand(oPerson *Person)
	Sit(oPerson *Person)
	Crouch(oPerson *Person)
	Jump(oPerson *Person)
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

func (oSelf *Person) Stand() {
	oSelf.state.Stand(oSelf)
}

func (oSelf *Person) Sit() {
	oSelf.state.Sit(oSelf)
}

func (oSelf *Person) Crouch() {
	oSelf.state.Crouch(oSelf)
}

func (oSelf *Person) Jump() {
	oSelf.state.Jump(oSelf)
}

// -----------------------------------------------------------------------------
// 具體狀態：站立
// -----------------------------------------------------------------------------

type StandingState struct {
}

func (oSelf *StandingState) Name() string {
	return "站立"
}

func (oSelf *StandingState) Stand(oPerson *Person) {
	fmt.Println("[站立] 已經站著了")
}

func (oSelf *StandingState) Sit(oPerson *Person) {
	fmt.Println("[站立] 坐下")
	oPerson.setState(&SittingState{})
}

func (oSelf *StandingState) Crouch(oPerson *Person) {
	fmt.Println("[站立] 蹲下")
	oPerson.setState(&CrouchingState{})
}

func (oSelf *StandingState) Jump(oPerson *Person) {
	fmt.Println("[站立] 起跳")
	oPerson.setState(&JumpingState{})
}

// -----------------------------------------------------------------------------
// 具體狀態：坐下
// -----------------------------------------------------------------------------

type SittingState struct {
}

func (oSelf *SittingState) Name() string {
	return "坐下"
}

func (oSelf *SittingState) Stand(oPerson *Person) {
	fmt.Println("[坐下] 起來")
	oPerson.setState(&StandingState{})
}

func (oSelf *SittingState) Sit(oPerson *Person) {
	fmt.Println("[坐下] 已經坐著了")
}

func (oSelf *SittingState) Crouch(oPerson *Person) {
	fmt.Println("[坐下] 坐著沒辦法直接蹲，要先站起來")
}

func (oSelf *SittingState) Jump(oPerson *Person) {
	fmt.Println("[坐下] 坐著沒辦法跳，要先站起來")
}

// -----------------------------------------------------------------------------
// 具體狀態：蹲著
// -----------------------------------------------------------------------------

type CrouchingState struct {
}

func (oSelf *CrouchingState) Name() string {
	return "蹲著"
}

func (oSelf *CrouchingState) Stand(oPerson *Person) {
	fmt.Println("[蹲著] 起來")
	oPerson.setState(&StandingState{})
}

func (oSelf *CrouchingState) Sit(oPerson *Person) {
	fmt.Println("[蹲著] 順勢坐下")
	oPerson.setState(&SittingState{})
}

func (oSelf *CrouchingState) Crouch(oPerson *Person) {
	fmt.Println("[蹲著] 已經蹲著了")
}

func (oSelf *CrouchingState) Jump(oPerson *Person) {
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

func (oSelf *JumpingState) Stand(oPerson *Person) {
	fmt.Println("[跳] 落地站好")
	oPerson.setState(&StandingState{})
}

func (oSelf *JumpingState) Sit(oPerson *Person) {
	fmt.Println("[跳] 在空中沒辦法坐")
}

func (oSelf *JumpingState) Crouch(oPerson *Person) {
	fmt.Println("[跳] 在空中沒辦法蹲")
}

func (oSelf *JumpingState) Jump(oPerson *Person) {
	fmt.Println("[跳] 已經在空中了，不能二段跳")
}

// -----------------------------------------------------------------------------

func main() {
	oPerson := NewPerson()

	oPerson.Sit()    // 站立 -> 坐下
	oPerson.Jump()   // 坐著不能跳
	oPerson.Stand()  // 坐下 -> 站立
	oPerson.Crouch() // 站立 -> 蹲著
	oPerson.Jump()   // 蹲著 -> 跳
	oPerson.Jump()   // 空中不能二段跳
	oPerson.Stand()  // 跳 -> 站立
}
