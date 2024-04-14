# Changelog

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
