# dynmgrm - GORM DynamoDB Driver

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
  - [ ] With `REMOVE` clause
- [x] Delete
- [ ] Create Table | Index

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
- [ ] `Not`
- [ ] `Or`

</details>

<details>
<summary>Commons</summary>

- [x] `Table`
- [ ] `Model`

</details>

<details>
<summary>Transaction</summary>

  - [x] `Begin`
  - [x] `Commit`
  - [x] `Rollback`
  - [ ] `Transaction`

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

And Custom Data Types:

- `Sets[string | int | float64 | []byte]`

- `List`

- `Map`

## Quick Start

### Installation

```.sh
go get github.com/miyamo2/dynmgrm
```

### Usage

```.go
```

## Contributing

Feel free to open a PR or an Issue.

## License

**dynmgrm** released under the [MIT License](https://github.com/miyamo2/dynmgrm/blob/master/LICENSE)