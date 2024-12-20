# Changelog

## 0.10.0 - 2024-11-05

### 🐭 Required Go version

Required Go version is now 1.22.0 or later.

### ⬆️ Upgrading dependencies

- `github.com/aws/aws-sdk-go-v2/service/dynamodb` 1.36.1 -> 1.36.3
- `go.uber.org/mock` 0.4.0 -> 0.5.0

## 0.9.1 - 2024-10-09

⬆️ Upgrading dependencies

- `github.com/aws/aws-sdk-go-v2/service/dynamodb` 1.34.3 -> 1.36.1
- `gorm.io/gorm` 1.25.10 -> 1.25.12
- `github.com/miyamo2/sqldav` 0.2.0 -> 0.2.1

## 0.9.0 - 2024-07-13

### 💥 Breaking Changes

- **Removed custom types** 

    `dynmgrm.List`, `dynmgrm.Map`, `dynmgrm.Set`, and `dynmgrm.TypedList` are removed.  
    Please use [`miyamoto/sqldav`](https://github.com/miyamo2/sqldav) instead from now on.

- **Replace SQL driver**

## 0.8.2 - 2024-06-23

Only a few GoDoc fixes

## 0.8.1 - 2024-06-23

Deprecated unsupported methods of `Migrator`

## 0.8.0 - 2024-06-22

✨ New Features

- GSI creation by `Migrator.CreateIndex` is now supported

## 0.7.0 - 2024-05-11

✨ New Features

- Added `dynamo-nested`, the custom GORM serializer for nested struct.

💥 Breaking Change

- Renamed the key name of tag, `dynmgrm:lsi` to `dynmgrm:lsi-sk`

## 0.6.1 - 2024-04-28

📝 Fixed Dead Link in README

## 0.6.0 - 2024-04-28

✨ New Features

- Added support for `Migratior.CreateTable()`.

## 0.5.0 - 2024-04-16

✨ New Features

- Added `ListAppend()`, a helper to the `list_append` function

## 0.4.0 - 2024-04-14

✨ New Features

- Add custom types

    - `TypedList[T]`

## 0.3.0 - 2024-04-06

💥 Breaking Change

- Renamed `DataType` to `KeySchemaDataType`

⬆️ Upgrading dependencies

- `github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue` from `1.13.12` to `1.13.13`

- `gorm.io/gorm` from `1.25.8` to `1.25.9`

📝 Fixed references to external libraries in godoc.

  
## 0.2.0 - 2024-03-26

✨ New Features

- Add Support following GORM features
    - `DB.Transaction`

## 0.1.2 - 2024-03-26

♻️ Few Refactor to make it testable

📝 Added `GormValue` explain to List, Map, Set

## 0.1.1 - 2024-03-26

## 0.1.0 - 2024-03-26

✨ New Features

- Add Support following PartiQL operations
    - SELECT
      - With Secondary Index
      - With `begins_with` function
      - With `contains` function
      - With `size` function
      - With `attribute_type` function
      - With `MISSING` operator
    - UPDATE
      - With `list_append` function
      - With `set_add` function
      - With `set_delete` function 
    - DELETE
    - INSERT
  
- Add Support following GORM features
    - `Table`
    - `Model`
    - `Select`
    - `Where`
    - `Or`
    - `Not`
    - `Find`
    - `Scan`
    - `Update`
    - `Updates`
    - `Save`
    - `Create`
    - `Delete`
    - `Begin`
    - `Commit`
    - `Rollback`
  
- Add Custom Type for DynamoDB Document/Set types
    - `List`
    - `Map`
    - `Set[string | int | float64 | []byte]`

- Add Custom Cluser
    - `SecondaryIndex`
