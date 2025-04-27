package dialect

import (
	"strings"
)

const (
	POSTGRES = "postgres"
	MYSQL    = "mysql"
)

type Dialect interface {
	Sprintd(format string, args ...interface{}) string
	Fprintd(builder *strings.Builder, format string, args ...interface{}) (int, error)
	Serial() string
	SmallSerial() string
	BigSerial() string
	BigInt() string
	Int() string
	SmallInt() string
	Boolean() string
	Char(length int) string
	VarChar(length int) string
	Text() string
	Placeholder(index int) string
	QuoteIdentifier(name string) string
	Real() string
	DoublePrecision() string
	Numeric(precision, scale int) string
	Money() string
	Timestamp() string
	TimestampTz() string
	Date() string
	Time() string
	TimeTz() string
	Interval() string
	Bytea() string
	Bit(length int) string
	VarBit(length int) string
	Uuid() string
	Array() string
	Json() string
	JsonB() string
	Xml() string
	Point() string
	Line() string
	Lseg() string
	Box() string
	Path() string
	Polygon() string
	Circle() string
	Cidr() string
	Inet() string
	MacAddr() string
	TsVector() string
	TsQuery() string
	TxidSnapshot() string
	Int4Range() string
	Int8Range() string
	NumRange() string
	TsRange() string
	TstzRange() string
	DateRange() string
	MacAddr8() string
}
