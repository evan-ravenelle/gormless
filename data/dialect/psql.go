package dialect

import (
	"fmt"
	"strings"
)

const (
	PsqlSerial          = "SERIAL"              // auto-incrementing four-byte integer
	PsqlSmallSerial     = "SMALLSERIAL"         // auto-incrementing two-byte integer
	PsqlBigSerial       = "BIGSERIAL"           // auto-incrementing eight-byte integer
	PsqlBigInt          = "BIGINT"              // No length or precision specifiers
	PsqlInt             = "INTEGER"             // No length or precision specifiers
	PsqlSmallInt        = "SMALLINT"            // No length or precision specifiers
	PsqlBoolean         = "BOOLEAN"             // No length or precision specifiers
	PsqlChar            = "CHAR(%d)"            // Length specifier
	PsqlVarChar         = "VARCHAR(%d)"         // Length specifier
	PsqlText            = "TEXT"                // No length or precision specifiers
	PsqlReal            = "REAL"                // No length or precision specifiers
	PsqlDoublePrecision = "DOUBLE PRECISION"    // No length or precision specifiers
	PsqlNumeric         = "NUMERIC(%d, %d)"     // Precision and Scale specifiers
	PsqlMoney           = "MONEY"               // No length or precision specifiers
	PsqlTimestamp       = "TIMESTAMP"           // No length or precision specifiers
	PsqlTimestampTz     = "TIMESTAMPTZ"         // No length or precision specifiers
	PsqlDate            = "DATE"                // No length or precision specifiers
	PsqlTime            = "TIME"                // No length or precision specifiers
	PsqlTimeTz          = "TIME WITH TIME ZONE" // No length or precision specifiers
	PsqlInterval        = "INTERVAL"            // No length or precision specifiers
	PsqlBytea           = "BYTEA"               // No length or precision specifiers
	PsqlBit             = "BIT(%d)"             // Length specifier
	PsqlVarBit          = "VARBIT(%d)"          // Length specifier
	PsqlUuid            = "UUID"                // No length or precision specifiers
	PsqlArray           = "ARRAY"               // No length or precision specifiers
	PsqlJson            = "JSON"                // No length or precision specifiers
	PsqlJsonB           = "JSONB"               // No length or precision specifiers
	PsqlXml             = "XML"                 // No length or precision specifiers
	PsqlPoint           = "POINT"               // No length or precision specifiers
	PsqlLine            = "LINE"                // No length or precision specifiers
	PsqlLseg            = "LSEG"                // No length or precision specifiers
	PsqlBox             = "BOX"                 // No length or precision specifiers
	PsqlPath            = "PATH"                // No length or precision specifiers
	PsqlPolygon         = "POLYGON"             // No length or precision specifiers
	PsqlCircle          = "CIRCLE"              // No length or precision specifiers
	PsqlCidr            = "CIDR"                // No length or precision specifiers
	PsqlInet            = "INET"                // No length or precision specifiers
	PsqlMacAddr         = "MACADDR"             // No length or precision specifiers
	PsqlTsVector        = "TSVECTOR"            // No length or precision specifiers
	PsqlTsQuery         = "TSQUERY"             // No length or precision specifiers
	PsqlTxidSnapshot    = "TXID_SNAPSHOT"       // No length or precision specifiers
	PsqlInt4Range       = "INT4RANGE"           // No length or precision specifiers
	PsqlInt8Range       = "INT8RANGE"           // No length or precision specifiers
	PsqlNumRange        = "NUMRANGE"            // No length or precision specifiers
	PsqlTsRange         = "TSRANGE"             // No length or precision specifiers
	PsqlTstzRange       = "TSTZRANGE"           // No length or precision specifiers
	PsqlDateRange       = "DATERANGE"           // No length or precision specifiers
	PsqlMacAddr8        = "MACADDR8"            // No length or precision specifiers
)

type PostgresDialect struct{}

func (p PostgresDialect) Sprintd(format string, args ...interface{}) string {
	// Process the args to quote any identifiers
	processedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		if s, ok := arg.(string); ok && format[strings.Index(format, "%"):strings.Index(format, "%")+2] == "%i" {
			// If the format specifier is %i, quote the string as an identifier
			processedArgs[i] = p.QuoteIdentifier(s)
		} else {
			processedArgs[i] = arg
		}
	}

	// Replace %i with %s in the format string
	format = strings.ReplaceAll(format, "%i", "%s")

	// Use fmt.Sprintf with the processed arguments
	return fmt.Sprintf(format, processedArgs...)
}

func (p PostgresDialect) Fprintd(builder *strings.Builder, format string, args ...interface{}) (int, error) {
	// Process the args to quote any identifiers
	processedArgs, processedFormat := p.processFormat(format, args...)

	// Use fmt.Fprintf with the processed arguments
	return fmt.Fprintf(builder, processedFormat, processedArgs...)
}

// Helper method to process format string and arguments
func (p PostgresDialect) processFormat(format string, args ...interface{}) ([]interface{}, string) {
	processedArgs := make([]interface{}, len(args))

	// Find all %i format specifiers
	var specifierIndices []int
	for i := 0; i < len(format)-1; i++ {
		if format[i] == '%' && i+1 < len(format) && format[i+1] == 'i' {
			specifierIndices = append(specifierIndices, i)
		}
	}

	// Process the arguments
	argIndex := 0
	for i, arg := range args {
		if argIndex < len(specifierIndices) {
			// If this argument corresponds to a %i specifier
			if s, ok := arg.(string); ok {
				processedArgs[i] = p.QuoteIdentifier(s)
			} else {
				processedArgs[i] = arg
			}
			argIndex++
		} else {
			processedArgs[i] = arg
		}
	}

	// Replace %i with %s in the format string
	processedFormat := strings.ReplaceAll(format, "%i", "%s")

	return processedArgs, processedFormat
}

func (p PostgresDialect) Serial() string            { return PsqlSerial }
func (p PostgresDialect) SmallSerial() string       { return PsqlSmallSerial }
func (p PostgresDialect) BigSerial() string         { return PsqlBigSerial }
func (p PostgresDialect) BigInt() string            { return PsqlBigInt }
func (p PostgresDialect) Int() string               { return PsqlInt }
func (p PostgresDialect) SmallInt() string          { return PsqlSmallInt }
func (p PostgresDialect) Boolean() string           { return PsqlBoolean }
func (p PostgresDialect) Char(length int) string    { return fmt.Sprintf(PsqlChar, length) }
func (p PostgresDialect) VarChar(length int) string { return fmt.Sprintf(PsqlVarChar, length) }
func (p PostgresDialect) Text() string              { return PsqlText }
func (p PostgresDialect) Real() string              { return PsqlReal }
func (p PostgresDialect) DoublePrecision() string   { return PsqlDoublePrecision }
func (p PostgresDialect) Numeric(precision, scale int) string {
	return fmt.Sprintf(PsqlNumeric, precision, scale)
}
func (p PostgresDialect) Money() string            { return PsqlMoney }
func (p PostgresDialect) Timestamp() string        { return PsqlTimestamp }
func (p PostgresDialect) TimestampTz() string      { return PsqlTimestampTz }
func (p PostgresDialect) Date() string             { return PsqlDate }
func (p PostgresDialect) Time() string             { return PsqlTime }
func (p PostgresDialect) TimeTz() string           { return PsqlTimeTz }
func (p PostgresDialect) Interval() string         { return PsqlInterval }
func (p PostgresDialect) Bytea() string            { return PsqlBytea }
func (p PostgresDialect) Bit(length int) string    { return fmt.Sprintf(PsqlBit, length) }
func (p PostgresDialect) VarBit(length int) string { return fmt.Sprintf(PsqlVarBit, length) }
func (p PostgresDialect) Uuid() string             { return PsqlUuid }
func (p PostgresDialect) Array() string            { return PsqlArray }
func (p PostgresDialect) Json() string             { return PsqlJson }
func (p PostgresDialect) JsonB() string            { return PsqlJsonB }
func (p PostgresDialect) Xml() string              { return PsqlXml }
func (p PostgresDialect) Point() string            { return PsqlPoint }
func (p PostgresDialect) Line() string             { return PsqlLine }
func (p PostgresDialect) Lseg() string             { return PsqlLseg }
func (p PostgresDialect) Box() string              { return PsqlBox }
func (p PostgresDialect) Path() string             { return PsqlPath }
func (p PostgresDialect) Polygon() string          { return PsqlPolygon }
func (p PostgresDialect) Circle() string           { return PsqlCircle }
func (p PostgresDialect) Cidr() string             { return PsqlCidr }
func (p PostgresDialect) Inet() string             { return PsqlInet }
func (p PostgresDialect) MacAddr() string          { return PsqlMacAddr }
func (p PostgresDialect) TsVector() string         { return PsqlTsVector }
func (p PostgresDialect) TsQuery() string          { return PsqlTsQuery }
func (p PostgresDialect) TxidSnapshot() string     { return PsqlTxidSnapshot }
func (p PostgresDialect) Int4Range() string        { return PsqlInt4Range }
func (p PostgresDialect) Int8Range() string        { return PsqlInt8Range }
func (p PostgresDialect) NumRange() string         { return PsqlNumRange }
func (p PostgresDialect) TsRange() string          { return PsqlTsRange }
func (p PostgresDialect) TstzRange() string        { return PsqlTstzRange }
func (p PostgresDialect) DateRange() string        { return PsqlDateRange }
func (p PostgresDialect) MacAddr8() string         { return PsqlMacAddr8 }

func (p PostgresDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

func (p PostgresDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf("\"%s\"", name)
}
