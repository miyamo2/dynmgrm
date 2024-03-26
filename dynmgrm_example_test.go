package dynmgrm_test

import (
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm"
)

type Event struct {
	Name  string `gorm:"primaryKey"`
	Date  string `gorm:"primaryKey"`
	Host  string
	Guest dynmgrm.Set[string]
}

func Example() {
	db, err := gorm.Open(dynmgrm.New())
	if err != nil {
		panic(err)
	}

	var dynamoDBWorkshop Event
	db.Table("events").Where(`name=? AND date=?`, "DynamoDB Workshop", "2024/3/25").Scan(&dynamoDBWorkshop)

	dynamoDBWorkshop.Guest = append(dynamoDBWorkshop.Guest, "Alice")
	db.Save(&dynamoDBWorkshop)

	carolBirthday := Event{
		Name:  "Carol's Birthday",
		Date:  "2024/4/1",
		Host:  "Charlie",
		Guest: []string{"Alice", "Bob"},
	}
	db.Create(carolBirthday)

	var daveSchedule []Event
	db.Table("events").
		Where(`date=? AND ( host=? OR CONTAINS("guest", ?) )`, "2024/4/1", "Dave", "Dave").
		Scan(&daveSchedule)

	tx := db.Begin()
	for _, event := range daveSchedule {
		if event.Host == "Dave" {
			tx.Delete(&event)
		} else {
			tx.Model(&event).Update("guest", gorm.Expr("set_delete(guest, ?)", dynmgrm.Set[string]{"Dave"}))
		}
	}
	tx.Model(&carolBirthday).Update("guest", gorm.Expr("set_add(guest, ?)", dynmgrm.Set[string]{"Dave"}))
	tx.Commit()

	var hostDateIndex []Event
	db.Table("events").Clauses(
		dynmgrm.SecondaryIndex("host-date-index"),
	).Where(`host=?`, "Bob").Scan(&hostDateIndex)
}

func ExampleNew() {
	gorm.Open(dynmgrm.New())
}

func ExampleNew_withRegion() {
	gorm.Open(dynmgrm.New(dynmgrm.WithRegion("ap-northeast-1")))
}

func ExampleNew_withAccessKeyID() {
	gorm.Open(dynmgrm.New(dynmgrm.WithAccessKeyID("YourAccess")))
}

func ExampleNew_withSecretKey() {
	gorm.Open(dynmgrm.New(dynmgrm.WithSecretKey("YourSecretKey")))
}

func ExampleNew_withEndpoint() {
	gorm.Open(dynmgrm.New(dynmgrm.WithEndpoint("http://localhost:8000")))
}

func ExampleNew_withTimeout() {
	gorm.Open(dynmgrm.New(dynmgrm.WithTimeout(30000)))
}

func ExampleWithRegion() {
	dynmgrm.WithRegion("ap-northeast-1")
}

func ExampleWithAccessKeyID() {
	dynmgrm.WithAccessKeyID("YourAccess")
}

func ExampleWithSecretKey() {
	dynmgrm.WithSecretKey("YourSecretKey")
}

func ExampleWithEndpoint() {
	dynmgrm.WithEndpoint("http://localhost:8000")
}

func ExampleWithTimeout() {
	dynmgrm.WithTimeout(30000)
}

func ExampleOpen() {
	gorm.Open(dynmgrm.Open("region=ap-northeast-1;AkId=YourAccessKeyID;SecretKey=YourSecretKey"))
}
