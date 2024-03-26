//go:generate mockgen -destination=internal/mocks/mock_dynmgrm.go -package=mocks -source=./dynmgrm.go
package dynmgrm

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// compatibility check
var _ gorm.Dialector = (*Dialector)(nil)

// KeySchemaDataType is the data type for the DynamoDB key schema.
type KeySchemaDataType string

// String returns the string representation of the KeySchemaDataType.
func (d KeySchemaDataType) String() string {
	return string(d)
}

const (
	// DriverName is the driver name for DynamoDB.
	DriverName = "godynamo"
	DBName     = "dynamodb"
)

// Define KeySchemaDataType
const (
	// KeySchemaDataTypeString is the data type for string.
	KeySchemaDataTypeString KeySchemaDataType = "string"
	// KeySchemaDataTypeNumber is the data type for number.
	KeySchemaDataTypeNumber KeySchemaDataType = "number"
	// KeySchemaDataTypeBinary is the data type for binary.
	KeySchemaDataTypeBinary KeySchemaDataType = "binary"
)

var (
	queryClauses   = []string{"SELECT", "FROM", "WHERE"}
	createClauses  = []string{"INSERT", "VALUES"}
	updateClauses  = []string{"UPDATE", "SET", "WHERE"}
	deleteClauses  = []string{"DELETE", "FROM", "WHERE"}
	clauseBuilders = map[string]clause.ClauseBuilder{
		"VALUES": toClauseBuilder(buildValuesClause),
		"SET":    toClauseBuilder(buildSetClause),
	}
)

// config is the configuration for the DynamoDB connection.
type config struct {
	region   string
	akId     string
	secret   string
	endpoint string
	timeout  int
	conn     gorm.ConnPool
}

// DBOpener is the interface for opening a database.
type DBOpener interface {
	DSN() string
	DriverName() string
	Apply() (*sql.DB, error)
}

type CallbacksRegisterer interface {
	Register(db *gorm.DB, config *callbacks.Config)
}

// Dialector gorm dialector for DynamoDB
type Dialector struct {
	conn gorm.ConnPool
	// dbOpener is used for testing
	dbOpener DBOpener
	// callbacksRegisterer is used for testing
	callbacksRegisterer CallbacksRegisterer
}

// DialectorOption is the option for the DynamoDB dialector.
type DialectorOption func(*config)

// WithRegion sets the region for the DynamoDB connection.
//
// Default: https://github.com/btnguyen2k/godynamo?tab=readme-ov-file#data-source-name-dsn-format-for-aws-dynamodb
func WithRegion(region string) func(*config) {
	return func(config *config) {
		config.region = region
	}
}

// WithAccessKeyID sets the access key ID for the DynamoDB connection.
//
// Default: https://github.com/btnguyen2k/godynamo?tab=readme-ov-file#data-source-name-dsn-format-for-aws-dynamodb
func WithAccessKeyID(accessKeyId string) func(*config) {
	return func(config *config) {
		config.akId = accessKeyId
	}
}

// WithSecretKey sets the secret key for the DynamoDB connection.
//
// Default: https://github.com/btnguyen2k/godynamo?tab=readme-ov-file#data-source-name-dsn-format-for-aws-dynamodb
func WithSecretKey(secretKey string) func(*config) {
	return func(config *config) {
		config.secret = secretKey
	}
}

// WithEndpoint sets the endpoint for the DynamoDB connection.
func WithEndpoint(endpoint string) func(*config) {
	return func(config *config) {
		config.endpoint = endpoint
	}
}

// WithTimeout sets the timeout milliseconds for the DynamoDB connection.
//
// Default: https://github.com/btnguyen2k/godynamo?tab=readme-ov-file#data-source-name-dsn-format-for-aws-dynamodb
func WithTimeout(timeout int) func(*config) {
	return func(config *config) {
		config.timeout = timeout
	}
}

// WithConnection sets the exist connection for the DynamoDB.
func WithConnection(conn gorm.ConnPool) func(*config) {
	return func(config *config) {
		config.conn = conn
	}
}

// Open returns a new DynamoDB dialector based on the DSN.
//
// e.g. "region=ap-northeast-1;AkId=<YOUR_ACCESS_KEY_ID>;SecretKey=<YOUR_SECRET_KEY>"
func Open(dsn string) gorm.Dialector {
	return &Dialector{
		dbOpener:            dbOpener{dsn: dsn, driverName: DriverName},
		callbacksRegisterer: &callbacksRegisterer{},
	}
}

// New returns a new DynamoDB dialector with options.
func New(option ...DialectorOption) gorm.Dialector {
	conf := config{}
	buildConfig(&conf, option...)
	return &Dialector{
		conn:                conf.conn,
		dbOpener:            dbOpener{dsn: parseConnectionString(conf), driverName: DriverName},
		callbacksRegisterer: &callbacksRegisterer{},
	}
}

func buildConfig(conf *config, option ...DialectorOption) {
	for _, opt := range option {
		opt(conf)
	}
}

func parseConnectionString(config config) string {
	dsnbuf := strings.Builder{}
	if config.region != "" {
		writeConnectionParameter(&dsnbuf, "region", config.region)
	}
	if config.akId != "" {
		writeConnectionParameter(&dsnbuf, "akId", config.akId)
	}
	if config.secret != "" {
		writeConnectionParameter(&dsnbuf, "secretKey", config.secret)
	}
	if config.endpoint != "" {
		writeConnectionParameter(&dsnbuf, "endpoint", config.endpoint)
	}
	if config.timeout != 0 {
		writeConnectionParameter(&dsnbuf, "timeout", strconv.Itoa(config.timeout))
	}
	return dsnbuf.String()
}

func writeConnectionParameter(dsnbuf *strings.Builder, key, value string) {
	if dsnbuf.Len() > 0 {
		dsnbuf.WriteString(";")
	}
	dsnbuf.WriteString(fmt.Sprintf("%s=%s", key, value))
}

// Name returns the name of the db.
func (dialector Dialector) Name() string {
	return DBName
}

// Initialize initializes the DynamoDB connection.
func (dialector Dialector) Initialize(db *gorm.DB) (err error) {
	if dialector.conn != nil {
		db.ConnPool = dialector.conn
	} else {
		conn, err := dialector.dbOpener.Apply()
		if err != nil {
			return err
		}
		db.ConnPool = conn
	}
	dialector.callbacksRegisterer.Register(
		db,
		&callbacks.Config{
			CreateClauses: createClauses,
			QueryClauses:  queryClauses,
			UpdateClauses: updateClauses,
			DeleteClauses: deleteClauses,
		})

	for k, v := range clauseBuilders {
		db.ClauseBuilders[k] = v
	}
	return
}

// DefaultValueOf returns the default value of the field.
func (dialector Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: ""}
}

// BindVarTo writes the bind variable of [goodynamo] to [clauses.Writer].
//
// [goodynamo]: https://pkg.go.dev/github.com/btnguyen2k/godynamo
// [clauses.Writer]: https://pkg.go.dev/gorm.io/gorm/clause#Writer
func (dialector Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteByte('?')
}

// QuoteTo escapes identifiers in SQL queries
func (dialector Dialector) QuoteTo(writer clause.Writer, str string) {
	writer.WriteString(fmt.Sprintf(`"%s"`, str))
}

// Explain returns the SQL string with the variables replaced.
// Explain is typically used only for logging, dry runs, and migration.
func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `"`, vars...)
}

// DataTypeOf maps GORM's data types to DynamoDB's data types.
// DataTypeOf works only with migration, so it will not return data types that are not allowed in PK, SK.
func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return KeySchemaDataTypeString.String()
	case schema.Int, schema.Uint, schema.Float:
		return KeySchemaDataTypeNumber.String()
	case schema.String:
		return KeySchemaDataTypeString.String()
	case schema.Time:
		return KeySchemaDataTypeString.String()
	case schema.Bytes:
		return KeySchemaDataTypeBinary.String()
	default:
		return KeySchemaDataTypeString.String()
	}
}

// Migrator returns the migrator for DynamoDB.
//
// Deprecated: Migration feature is not implemented.
func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return &Migrator{}
}
