package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================
// gorm Many To Many（多對多）
// ============================================================
//
// Student ── n:n ── Class，中間靠 join table 串起來。
//
// 1. 最單純：兩邊各放一個 slice，加 `gorm:"many2many:<join_table>"`
//      type Student struct{ Classes []Class   `gorm:"many2many:student_classes"` }
//      type Class   struct{ Students []Student `gorm:"many2many:student_classes"` }
//    gorm 會自動建一張只有 (student_id, class_id) 兩欄的 join table。
//
// 2. 自訂 join table（join 上要放額外欄位，例如成績、加入時間）：
//      type StudentClass struct {
//          StudentId  uint `gorm:"primaryKey"`
//          ClassId    uint `gorm:"primaryKey"`
//          Score      int
//          EnrolledAt time.Time
//      }
//    AutoMigrate 前，「雙向都要」呼叫 SetupJoinTable：
//      db.SetupJoinTable(&Student{}, "Classes", &StudentClass{})
//      db.SetupJoinTable(&Class{},   "Students", &StudentClass{})
//    然後 AutoMigrate 只傳 &Student{}, &Class{}——不要傳 &StudentClass{}，
//    不然 gorm 會把它當普通表另外遷移，跟 join table 打架。
//
// 3. 操作關聯（不用手動塞 join table）：
//      db.Model(&s).Association("Classes").Append(&c)   // 加
//      db.Model(&s).Association("Classes").Delete(&c)   // 移除（不刪 Class 本身）
//      db.Model(&s).Association("Classes").Replace(...) // 整組換掉
//      db.Model(&s).Association("Classes").Count()      // 數量
//
// 4. 查詢：Preload("Classes") / Preload("Students")；
//    join table 的額外欄位 Preload 不會帶出來，要自己查 StudentClass。
// ============================================================

type Student struct {
	Id   uint
	Name string
	Age  int

	// many2many：透過 student_classes 這張 join table 關聯 Class
	Classes []Class `gorm:"many2many:student_classes"`
}

type Class struct {
	Id   uint
	Name string

	// 反向：同一張 join table，反過來拿 Student
	Students []Student `gorm:"many2many:student_classes"`
}

// StudentClass 是自訂的 join table：除了兩個外鍵，額外放 Score / EnrolledAt。
// 複合主鍵 (student_id, class_id) 保證同一個學生不會重複加入同一個班。
type StudentClass struct {
	StudentId  uint `gorm:"primaryKey"`
	ClassId    uint `gorm:"primaryKey"`
	Score      int
	EnrolledAt time.Time
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug()

	// 雙向 SetupJoinTable，都在 AutoMigrate 之前
	if err := db.SetupJoinTable(&Student{}, "Classes", &StudentClass{}); err != nil {
		log.Fatal(err)
	}
	if err := db.SetupJoinTable(&Class{}, "Students", &StudentClass{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&Student{}, &Class{}); err != nil {
		log.Fatal(err)
	}

	// ===== 先建 Class（讓它們有 Id），再建 Student 時直接引用，gorm 就只補 join、不會重複建課 =====
	oMath := Class{Name: "數學"}
	oEnglish := Class{Name: "英文"}
	oPE := Class{Name: "體育"}
	db.Create(&oMath)
	db.Create(&oEnglish)
	db.Create(&oPE)

	db.Create(&Student{Name: "Alice", Age: 18, Classes: []Class{oMath, oEnglish}})
	db.Create(&Student{Name: "Bob", Age: 19, Classes: []Class{oMath, oPE}})

	// ===== 1. Preload：Student -> 他修的 Classes =====
	var oAlice Student
	db.Preload("Classes").First(&oAlice, "name = ?", "Alice")
	fmt.Printf("\n[Student->Classes] %s 修：", oAlice.Name)
	for _, o := range oAlice.Classes {
		fmt.Printf(" %s", o.Name)
	}
	fmt.Println()

	// ===== 2. 反向 Preload：Class -> 修課的 Students =====
	var oMathClass Class
	db.Preload("Students").First(&oMathClass, "name = ?", "數學")
	fmt.Printf("[Class->Students]  %s 班：", oMathClass.Name)
	for _, o := range oMathClass.Students {
		fmt.Printf(" %s(%d)", o.Name, o.Age)
	}
	fmt.Println()

	// ===== 3. Association API：Alice 加簽體育、退選英文 =====
	db.Model(&oAlice).Association("Classes").Append(&oPE)
	db.Model(&oAlice).Association("Classes").Delete(&oEnglish)
	fmt.Printf("[Association]      %s 加體育退英文後，共 %d 門\n",
		oAlice.Name, db.Model(&oAlice).Association("Classes").Count())

	// ===== 4. join table 額外欄位：直接操作 StudentClass =====
	db.Model(&StudentClass{}).
		Where("student_id = ? AND class_id = ?", oAlice.Id, oMath.Id).
		Updates(map[string]any{"score": 95, "enrolled_at": time.Now()})

	var oSC StudentClass
	db.First(&oSC, "student_id = ? AND class_id = ?", oAlice.Id, oMath.Id)
	fmt.Printf("[join extra]       %s 的數學成績：%d（%s 加入）\n",
		oAlice.Name, oSC.Score, oSC.EnrolledAt.Format("2006-01-02"))

	// ===== 5. Replace：Bob 的課整組換成只剩數學 =====
	var oBob Student
	db.First(&oBob, "name = ?", "Bob")
	db.Model(&oBob).Association("Classes").Replace(&oMath)

	db.Preload("Classes").First(&oBob, oBob.Id)
	fmt.Printf("[Replace]          %s 現在只修：", oBob.Name)
	for _, o := range oBob.Classes {
		fmt.Printf(" %s", o.Name)
	}
	fmt.Println()
}
