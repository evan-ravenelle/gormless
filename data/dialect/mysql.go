package dialect

import (
	"fmt"
	"strings"
)

const (
	// MySQL data type constants
	MySqlTinyInt            = "TINYINT"            // 1 byte (-128 to 127)
	MySqlSmallInt           = "SMALLINT"           // 2 bytes (-32,768 to 32,767)
	MySqlMediumInt          = "MEDIUMINT"          // 3 bytes (-8,388,608 to 8,388,607)
	MySqlInt                = "INT"                // 4 bytes (-2,147,483,648 to 2,147,483,647)
	MySqlBigInt             = "BIGINT"             // 8 bytes (-9,223,372,036,854,775,808 to 9,223,372,036,854,775,807)
	MySqlDecimal            = "DECIMAL(%d, %d)"    // Precision and scale specifiers
	MySqlFloat              = "FLOAT"              // 4 bytes
	MySqlDouble             = "DOUBLE"             // 8 bytes
	MySqlBit                = "BIT(%d)"            // Up to 64 bits
	MySqlChar               = "CHAR(%d)"           // Fixed-length string
	MySqlVarChar            = "VARCHAR(%d)"        // Variable-length string
	MySqlTinyText           = "TINYTEXT"           // Up to 255 characters
	MySqlText               = "TEXT"               // Up to 65,535 characters
	MySqlMediumText         = "MEDIUMTEXT"         // Up to 16,777,215 characters
	MySqlLongText           = "LONGTEXT"           // Up to 4,294,967,295 characters
	MySqlBinary             = "BINARY(%d)"         // Fixed-length binary string
	MySqlVarBinary          = "VARBINARY(%d)"      // Variable-length binary string
	MySqlTinyBlob           = "TINYBLOB"           // Up to 255 bytes
	MySqlBlob               = "BLOB"               // Up to 65,535 bytes
	MySqlMediumBlob         = "MEDIUMBLOB"         // Up to 16,777,215 bytes
	MySqlLongBlob           = "LONGBLOB"           // Up to 4,294,967,295 bytes
	MySqlEnum               = "ENUM(%s)"           // Enumeration type
	MySqlSet                = "SET(%s)"            // Set type
	MySqlDate               = "DATE"               // YYYY-MM-DD
	MySqlTime               = "TIME"               // HH:MM:SS
	MySqlDateTime           = "DATETIME"           // YYYY-MM-DD HH:MM:SS
	MySqlTimestamp          = "TIMESTAMP"          // YYYY-MM-DD HH:MM:SS
	MySqlYear               = "YEAR"               // YYYY
	MySqlBoolean            = "BOOLEAN"            // Alias for TINYINT(1)
	MySqlJson               = "JSON"               // JSON data type (MySQL 5.7.8 and up)
	MySqlGeometry           = "GEOMETRY"           // Spatial data type
	MySqlPoint              = "POINT"              // Spatial data type
	MySqlLineString         = "LINESTRING"         // Spatial data type
	MySqlPolygon            = "POLYGON"            // Spatial data type
	MySqlMultiPoint         = "MULTIPOINT"         // Spatial data type
	MySqlMultiLineString    = "MULTILINESTRING"    // Spatial data type
	MySqlMultiPolygon       = "MULTIPOLYGON"       // Spatial data type
	MySqlGeometryCollection = "GEOMETRYCOLLECTION" // Spatial data type
)

// MySQLDialect implements Dialect for MySQL
type MySQLDialect struct{}

func (m MySQLDialect) Sprintd(format string, args ...interface{}) string {
	// Process the args to quote any identifiers
	processedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		if s, ok := arg.(string); ok && format[strings.Index(format, "%"):strings.Index(format, "%")+2] == "%i" {
			// If the format specifier is %i, quote the string as an identifier
			processedArgs[i] = m.QuoteIdentifier(s)
		} else {
			processedArgs[i] = arg
		}
	}

	// Replace %i with %s in the format string
	format = strings.ReplaceAll(format, "%i", "%s")

	// Use fmt.Sprintf with the processed arguments
	return fmt.Sprintf(format, processedArgs...)
}

func (m MySQLDialect) Fprintd(builder *strings.Builder, format string, args ...interface{}) (int, error) {
	// Process the args to quote any identifiers
	processedArgs, processedFormat := m.processFormat(format, args...)

	// Use fmt.Fprintf with the processed arguments
	return fmt.Fprintf(builder, processedFormat, processedArgs...)
}

// Helper method to process format string and arguments
func (m MySQLDialect) processFormat(format string, args ...interface{}) ([]interface{}, string) {
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
				processedArgs[i] = m.QuoteIdentifier(s)
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

// Data type implementations using constants
func (m MySQLDialect) Serial() string            { return fmt.Sprintf("%s AUTO_INCREMENT", MySqlInt) }
func (m MySQLDialect) SmallSerial() string       { return fmt.Sprintf("%s AUTO_INCREMENT", MySqlSmallInt) }
func (m MySQLDialect) BigSerial() string         { return fmt.Sprintf("%s AUTO_INCREMENT", MySqlBigInt) }
func (m MySQLDialect) BigInt() string            { return MySqlBigInt }
func (m MySQLDialect) Int() string               { return MySqlInt }
func (m MySQLDialect) SmallInt() string          { return MySqlSmallInt }
func (m MySQLDialect) Boolean() string           { return MySqlBoolean }
func (m MySQLDialect) Char(length int) string    { return fmt.Sprintf(MySqlChar, length) }
func (m MySQLDialect) VarChar(length int) string { return fmt.Sprintf(MySqlVarChar, length) }
func (m MySQLDialect) Text() string              { return MySqlText }
func (m MySQLDialect) Real() string              { return MySqlFloat }
func (m MySQLDialect) DoublePrecision() string   { return MySqlDouble }
func (m MySQLDialect) Numeric(precision, scale int) string {
	return fmt.Sprintf(MySqlDecimal, precision, scale)
}
func (m MySQLDialect) Money() string         { return fmt.Sprintf(MySqlDecimal, 19, 4) } // MySQL has no direct MONEY type
func (m MySQLDialect) Timestamp() string     { return MySqlTimestamp }
func (m MySQLDialect) TimestampTz() string   { return MySqlTimestamp } // MySQL handles timezone differently
func (m MySQLDialect) Date() string          { return MySqlDate }
func (m MySQLDialect) Time() string          { return MySqlTime }
func (m MySQLDialect) TimeTz() string        { return MySqlTime } // MySQL handles timezone differently
func (m MySQLDialect) Bytea() string         { return MySqlBlob }
func (m MySQLDialect) Bit(length int) string { return fmt.Sprintf(MySqlBit, length) }
func (m MySQLDialect) Json() string          { return MySqlJson }
func (m MySQLDialect) Point() string         { return MySqlPoint }
func (m MySQLDialect) Line() string          { return MySqlLineString } // Approximation
func (m MySQLDialect) Lseg() string          { return MySqlLineString } // Approximation
func (m MySQLDialect) Box() string           { return MySqlPolygon }    // Approximation
func (m MySQLDialect) Path() string          { return MySqlLineString } // Approximation
func (m MySQLDialect) Polygon() string       { return MySqlPolygon }
func (m MySQLDialect) Circle() string        { return MySqlGeometry } // Approximation

func (m MySQLDialect) Interval() string {
	panic("MySQL does not support INTERVAL data type")
}

func (m MySQLDialect) VarBit(length int) string {
	panic("MySQL does not support VARBIT data type")
}

func (m MySQLDialect) Uuid() string {
	// Note: MySQL 8.0+ has UUID_TO_BIN and BIN_TO_UUID functions, but no native type
	panic("MySQL does not have a native UUID data type (use CHAR(36) instead)")
}

func (m MySQLDialect) Array() string {
	panic("MySQL does not support ARRAY data type")
}

func (m MySQLDialect) JsonB() string {
	panic("MySQL does not support JSONB data type (use JSON instead)")
}

func (m MySQLDialect) Xml() string {
	panic("MySQL does not have a native XML data type (use TEXT instead)")
}

func (m MySQLDialect) Cidr() string {
	panic("MySQL does not support CIDR data type")
}

func (m MySQLDialect) Inet() string {
	panic("MySQL does not support INET data type")
}

func (m MySQLDialect) MacAddr() string {
	panic("MySQL does not support MACADDR data type")
}

func (m MySQLDialect) TsVector() string {
	panic("MySQL does not support TSVECTOR data type (use full-text search instead)")
}

func (m MySQLDialect) TsQuery() string {
	panic("MySQL does not support TSQUERY data type (use full-text search instead)")
}

func (m MySQLDialect) TxidSnapshot() string {
	panic("MySQL does not support TXID_SNAPSHOT data type")
}

func (m MySQLDialect) Int4Range() string {
	panic("MySQL does not support range data types")
}

func (m MySQLDialect) Int8Range() string {
	panic("MySQL does not support range data types")
}

func (m MySQLDialect) NumRange() string {
	panic("MySQL does not support range data types")
}

func (m MySQLDialect) TsRange() string {
	panic("MySQL does not support range data types")
}

func (m MySQLDialect) TstzRange() string {
	panic("MySQL does not support range data types")
}

func (m MySQLDialect) DateRange() string {
	panic("MySQL does not support range data types")
}

func (m MySQLDialect) MacAddr8() string {
	panic("MySQL does not support MACADDR8 data type")
}

// Placeholder returns the placeholder for a given parameter index
func (m MySQLDialect) Placeholder(index int) string {
	return "?" // MySQL uses ? for all placeholders
}

// QuoteIdentifier quotes an identifier
func (m MySQLDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf("`%s`", name)
}
