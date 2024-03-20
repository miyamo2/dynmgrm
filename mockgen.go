//go:generate mockgen -destination=internal/mocks/mock_gorm_expression.go -package=mocks gorm.io/gorm/clause Expression
//go:generate mockgen -destination=internal/mocks/mock_gorm_builder.go -package=mocks gorm.io/gorm/clause Builder
package dynmgrm
