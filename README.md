# dynmgrm - GORM DynamoDB Driver

<img src=".assets/logo/svg/dynmgrm_logo_with_caption.svg" width="400" alt="logo">

[![Go Reference](https://pkg.go.dev/badge/github.com/miyamo2/dynmgrm.svg)](https://pkg.go.dev/github.com/miyamo2/dynmgrm)
[![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/miyamo2/dynmgrm?logo=go)](https://img.shields.io/github/go-mod/go-version/miyamo2/dynmgrm?logo=go)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/miyamo2/dynmgrm)](https://img.shields.io/github/v/release/miyamo2/dynmgrm)
[![codecov](https://codecov.io/gh/miyamo2/dynmgrm/graph/badge.svg?token=QLIVB3ESVD)](https://codecov.io/gh/miyamo2/dynmgrm)
[![Go Report Card](https://goreportcard.com/badge/github.com/miyamo2/dynmgrm)](https://goreportcard.com/report/github.com/miyamo2/dynmgrm)
[![GitHub License](https://img.shields.io/github/license/miyamo2/dynmgrm?&color=blue)](https://img.shields.io/github/license/miyamo2/dynmgrm?&color=blue)

## Features

Supports the following PartiQL operations:

- [x] Select
  - [x] With Secondary Index
  - [x] With `begins_with` function
  - [x] With `contains` function
  - [x] With `size` function
  - [x] With `attribute_type` function
  - [x] With `MISSING` operator
- [x] Insert
- [x] Update
  - [x] With `SET` clause
    - [x] With `list_append` function
      - [x] `ListAppend()`
    - [x] With `set_add` function
    - [x] With `set_delete` function
  - [ ] With `REMOVE` clause
- [x] Delete
- [ ] Create (Table | Index)

Supports the following GORM features:

<details>
<summary>Query</summary>

- [x] `Select`
- [x] `Find`
- [x] `Scan`

</details>

<details>
<summary>Update</summary>

- [x] `Update`
- [x] `Updates`
- [x] `Save`

</details>

<details>
<summary>Create</summary>

- [x] `Create`

</details>

<details>
<summary>Delete</summary>

- [x] `Delete`

</details>

<details>
<summary>Condition</summary>

- [x] `Where`
- [x] `Not`
- [x] `Or`

</details>

<details>
<summary>Table/Model</summary>

- [x] `Table`
- [x] `Model` ※ Combination with Secondary Index are not supported.

</details>

<details>
<summary>Transaction</summary>

  - [x] `Begin`
  - [x] `Commit`
  - [x] `Rollback`
  - [x] `Transaction`

※ Supports only Insert, Update, and Delete.

</details>

<details>
<summary>Migration</summary>

- [ ] `AutoMigrate`
- [ ] `CurrentDatabase`
- [ ] `FullDataTypeOf`
- [ ] `CreateTable`
- [ ] `DropTable`
- [ ] `HasTable`
- [ ] `RenameTable`
- [ ] `GetTables`
- [ ] `AddColumn`
- [ ] `DropColumn`
- [ ] `AlterColumn`
- [ ] `MigrateColumn`
- [ ] `HasColumn`
- [ ] `RenameColumn`
- [ ] `ColumnTypes`
- [ ] `CreateView`
- [ ] `DropView`
- [ ] `CreateConstraint`
- [ ] `DropConstraint`
- [ ] `HasConstraint`
- [ ] `CreateIndex`
- [ ] `DropIndex`
- [ ] `HasIndex`
- [ ] `RenameIndex`

</details>

Custom Clause:

- `SecondaryIndex`

Custom Data Types:

- `Set[string | int | float64 | []byte]`

- `List`

- `Map`

- `TypedList[T]`

## Quick Start

### Installation

```sh
go get github.com/miyamo2/dynmgrm
```

### Usage

```go
package main

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

func main() {
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
```

## Contributing

Feel free to open a PR or an Issue.

## License

**dynmgrm** released under the [MIT License](https://github.com/miyamo2/dynmgrm/blob/master/LICENSE)

## Credits

### Go gopher

The Go gopher was designed by [Renee French.](http://reneefrench.blogspot.com/)
The design is licensed under the Creative Commons 3.0 Attributions license.
Read this article for more [details](https://go.dev/blog/gopher)

### Special Thanks

- [JetBrainsMono](https://github.com/JetBrains/JetBrainsMono)

	JetBrainsMono is used for the caption of the dynmgrm logo.