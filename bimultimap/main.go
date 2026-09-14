package main

import (
	"fmt"
	"sort"

	"github.com/cornelk/hashmap"
)

// Hashable 是 cornelk/hashmap 接受的 key 型別集合（數值與字串）。
// 它自己的 key constraint 沒有匯出，這裡列一份等價的給 BiMultiMap 當型別參數約束用。
type Hashable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

type BiMultiMap[L, R Hashable] struct {
	leftToRights *hashmap.Map[L, *hashmap.Map[R, struct{}]]
	rightToLefts *hashmap.Map[R, *hashmap.Map[L, struct{}]]
}

func NewBiMultiMap[L, R Hashable]() *BiMultiMap[L, R] {
	oLeftToRights := hashmap.New[L, *hashmap.Map[R, struct{}]]()
	oRightToLefts := hashmap.New[R, *hashmap.Map[L, struct{}]]()
	oBiMultiMap := &BiMultiMap[L, R]{
		leftToRights: oLeftToRights,
		rightToLefts: oRightToLefts,
	}

	return oBiMultiMap
}

// Insert 建立 oLeft/oRight 的雙向關聯，重複呼叫同一組 oLeft/oRight 不會有副作用。
func (oSelf *BiMultiMap[L, R]) Insert(oLeft L, oRight R) {
	oNewRightsSet := hashmap.New[R, struct{}]()
	oRights, _ := oSelf.leftToRights.GetOrInsert(oLeft, oNewRightsSet)
	oRights.Set(oRight, struct{}{})

	oNewLeftsSet := hashmap.New[L, struct{}]()
	oLefts, _ := oSelf.rightToLefts.GetOrInsert(oRight, oNewLeftsSet)
	oLefts.Set(oLeft, struct{}{})
}

// Left 回傳這個 left 目前關聯到的所有 right。
func (oSelf *BiMultiMap[L, R]) Left(oLeft L) []R {
	oRights, bGotten := oSelf.leftToRights.Get(oLeft)
	if !bGotten {
		return nil
	}

	iRightsLen := oRights.Len()
	aResult := make([]R, 0, iRightsLen)
	oRights.Range(func(oRight R, _ struct{}) bool {
		aResult = append(aResult, oRight)
		return true
	})

	return aResult
}

// Right 回傳這個 right 目前關聯到的所有 left。
func (oSelf *BiMultiMap[L, R]) Right(oRight R) []L {
	oLefts, bGotten := oSelf.rightToLefts.Get(oRight)
	if !bGotten {
		return nil
	}

	iLeftsLen := oLefts.Len()
	aResult := make([]L, 0, iLeftsLen)
	oLefts.Range(func(oLeft L, _ struct{}) bool {
		aResult = append(aResult, oLeft)
		return true
	})

	return aResult
}

// Remove 拆掉單一一組 oLeft/oRight 的關聯，其他跟 oLeft 或 oRight 相關的關聯不受影響。
func (oSelf *BiMultiMap[L, R]) Remove(oLeft L, oRight R) {
	if oRights, bGotten := oSelf.leftToRights.Get(oLeft); bGotten {
		oRights.Del(oRight)
		if oRights.Len() == 0 {
			oSelf.leftToRights.Del(oLeft)
		}
	}

	if oLefts, bGotten := oSelf.rightToLefts.Get(oRight); bGotten {
		oLefts.Del(oLeft)
		if oLefts.Len() == 0 {
			oSelf.rightToLefts.Del(oRight)
		}
	}
}

// RemoveLeft 拆掉這個 left 的全部關聯，連帶清掉對面每個 right 記著這個 left 的紀錄。
func (oSelf *BiMultiMap[L, R]) RemoveLeft(oLeft L) {
	oRights, bGotten := oSelf.leftToRights.Get(oLeft)
	if !bGotten {
		return
	}

	oRights.Range(func(oRight R, _ struct{}) bool {
		if oLefts, bLefts := oSelf.rightToLefts.Get(oRight); bLefts {
			oLefts.Del(oLeft)
			if oLefts.Len() == 0 {
				oSelf.rightToLefts.Del(oRight)
			}
		}
		return true
	})

	oSelf.leftToRights.Del(oLeft)
}

// RemoveRight 拆掉這個 right 的全部關聯，連帶清掉對面每個 left 記著這個 right 的紀錄。
func (oSelf *BiMultiMap[L, R]) RemoveRight(oRight R) {
	oLefts, bGotten := oSelf.rightToLefts.Get(oRight)
	if !bGotten {
		return
	}

	oLefts.Range(func(oLeft L, _ struct{}) bool {
		if oRights, bRights := oSelf.leftToRights.Get(oLeft); bRights {
			oRights.Del(oRight)
			if oRights.Len() == 0 {
				oSelf.leftToRights.Del(oLeft)
			}
		}
		return true
	})

	oSelf.rightToLefts.Del(oRight)
}

func main() {
	oBimultimap := NewBiMultiMap[string, string]()

	// conn-1 訂了 #general 跟 #random；conn-2 訂了 #general；conn-3 訂了 #random。
	oBimultimap.Insert("conn-1", "#general")
	oBimultimap.Insert("conn-1", "#random")
	oBimultimap.Insert("conn-2", "#general")
	oBimultimap.Insert("conn-3", "#random")
	oBimultimap.Insert("conn-1", "#general") // 重複 Insert 同一組沒有副作用

	// Left：某條連線訂了哪些頻道。
	aConn1Channels := oBimultimap.Left("conn-1")
	aSortedConn1Channels := sorted(aConn1Channels)
	fmt.Println("conn-1 訂閱的頻道:", aSortedConn1Channels)

	// Right：某個頻道裡有哪些連線（推播就是查這個方向）。
	aGeneralSubscribers := oBimultimap.Right("#general")
	aSortedGeneralSubscribers := sorted(aGeneralSubscribers)
	fmt.Println("#general 的訂閱者:", aSortedGeneralSubscribers)

	aRandomSubscribers := oBimultimap.Right("#random")
	aSortedRandomSubscribers := sorted(aRandomSubscribers)
	fmt.Println("#random 的訂閱者:", aSortedRandomSubscribers)

	// 查不存在的 key 回傳 nil（可以直接 range）。
	aMissingKeyChannels := oBimultimap.Left("conn-404")
	fmt.Println("查不存在的 key:", aMissingKeyChannels == nil)
}

// demoRemove：三種拆關聯的方式。
func demoRemove() {
	oBimultimap := NewBiMultiMap[string, string]()

	oBimultimap.Insert("conn-1", "#general")
	oBimultimap.Insert("conn-1", "#random")
	oBimultimap.Insert("conn-2", "#general")
	oBimultimap.Insert("conn-2", "#random")
	oBimultimap.Insert("conn-3", "#random")

	// Remove：只拆單一一組 oLeft/oRight，其他不動。
	oBimultimap.Remove("conn-1", "#random")
	aConn1ChannelsAfterRemove := oBimultimap.Left("conn-1")
	aSortedConn1ChannelsAfterRemove := sorted(aConn1ChannelsAfterRemove)
	fmt.Println("Remove(conn-1, #random) 後，conn-1 的頻道:", aSortedConn1ChannelsAfterRemove)

	aRandomSubscribersAfterRemove := oBimultimap.Right("#random")
	aSortedRandomSubscribersAfterRemove := sorted(aRandomSubscribersAfterRemove)
	fmt.Println("Remove(conn-1, #random) 後，#random 的訂閱者:", aSortedRandomSubscribersAfterRemove)

	// RemoveLeft：某條連線斷線，拆掉它的所有訂閱（對面頻道也會同步移除這條連線）。
	oBimultimap.RemoveLeft("conn-2")
	aGeneralSubscribersAfterRemoveLeft := oBimultimap.Right("#general")
	aSortedGeneralSubscribersAfterRemoveLeft := sorted(aGeneralSubscribersAfterRemoveLeft)
	fmt.Println("RemoveLeft(conn-2) 後，#general 的訂閱者:", aSortedGeneralSubscribersAfterRemoveLeft)

	// RemoveRight：整個頻道關閉，拆掉所有人對它的訂閱。
	oBimultimap.RemoveRight("#random")
	aConn1ChannelsAfterRemoveRight := oBimultimap.Left("conn-1")
	aSortedConn1ChannelsAfterRemoveRight := sorted(aConn1ChannelsAfterRemoveRight)
	fmt.Println("RemoveRight(#random) 後，conn-1 的頻道:", aSortedConn1ChannelsAfterRemoveRight)

	aConn3ChannelsAfterRemoveRight := oBimultimap.Left("conn-3")
	aSortedConn3ChannelsAfterRemoveRight := sorted(aConn3ChannelsAfterRemoveRight)
	fmt.Println("RemoveRight(#random) 後，conn-3 的頻道:", aSortedConn3ChannelsAfterRemoveRight)
}

// sorted 把 Range 出來的無序 slice 排好，讓範例輸出穩定。
func sorted(a []string) []string {
	sort.Strings(a)
	return a
}
