package main

import "fmt"

// 策略模式（Strategy Pattern）
//
// 核心概念：把一組「可以互換的演算法」各自封裝成獨立物件，
// Context 持有一個策略介面，執行時委派給當前策略，
// 客戶端可以在執行期自由替換策略，不需改動 Context。
//
// 本範例：PriceCalculator 依不同「折扣策略」計算最終金額。

// -----------------------------------------------------------------------------
// Strategy 介面
// -----------------------------------------------------------------------------

type DiscountStrategy interface {
	Apply(iAmount int) int
}

// -----------------------------------------------------------------------------
// 具體策略
// -----------------------------------------------------------------------------

// NoDiscount 不打折
type NoDiscount struct {
}

func (oSelf *NoDiscount) Apply(iAmount int) int {
	return iAmount
}

// HalfDiscount 打五折
type HalfDiscount struct {
}

func (oSelf *HalfDiscount) Apply(iAmount int) int {
	return iAmount / 2
}

// -----------------------------------------------------------------------------
// Context：持有當前策略
// -----------------------------------------------------------------------------

type PriceCalculator struct {
	strategy DiscountStrategy
}

func (oSelf *PriceCalculator) SetStrategy(oStrategy DiscountStrategy) {
	oSelf.strategy = oStrategy
}

func (oSelf *PriceCalculator) Checkout(iAmount int) int {
	return oSelf.strategy.Apply(iAmount)
}

// -----------------------------------------------------------------------------

func main() {
	oCalculator := &PriceCalculator{}

	oCalculator.SetStrategy(&NoDiscount{})
	fmt.Println("無折扣：", oCalculator.Checkout(1000))

	oCalculator.SetStrategy(&HalfDiscount{})
	fmt.Println("打五折：", oCalculator.Checkout(1000))

	/*
			這樣看起來 strategy 模式跟 command 模式非常類似，他都是透過注入解耦，
		不一樣的是 strategy 關心 注入的程序邏輯， command 關心注入的東西的規格，

			**/
}
