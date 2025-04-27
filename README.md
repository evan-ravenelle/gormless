# gormless: A Type-Safe Database Schema Manager for Go

gormless is a lightweight database schema management tool for Go applications. Unlike traditional ORMs, gormless embraces SQL and provides a thin layer of abstraction to work with your database schema in a type-safe manner.

## Philosophy

gormless was created as an "anti-ORM" solution, based on the belief that:

- Database models and application models should be separate concerns
- Tables should be first-class concepts, not relegated to data storage for business objects
- SQL is a powerful language that developers should leverage, not hide from
- Schema migrations should be explicit and version-controlled
- Type safety should not come at the cost of performance or flexibility

## Features

- **Multi-dialect Support**: Works with PostgreSQL and MySQL out of the box
- **Type-Safe Operations**: Uses Go generics for compile-time type checking
- **Schema Management**: Create tables, add/modify columns, and manage constraints
- **SQL Injection Prevention**: Built-in validation of SQL identifiers
- **Transaction Support**: Manage database transactions with ease
- **Minimal Dependencies**: Relies on standard library where possible
- **Comprehensive Testing**: Includes both unit and integration tests

## Installation

```bash
go get github.com/evan-ravenelle/gormless
```

## Quick Start

### Define Your Database Configuration

```yaml
# db_config.yml
database:
  host: localhost
  port: 5432
  username: postgres
  password: your_secure_password
  dbname: myapp
  sslmode: disable
```

### Create a Session

```go
import (
    "fmt"
    "gormless/data"
    "gormless/data/dialect"
)

func main() {
    // Load configuration
    conf, err := data.LoadConfig("db_config.yml")
    if err != nil {
        panic(err)
    }
    
    // Create connection string
    dsn := fmt.Sprintf(
        "user=%s dbname=%s sslmode=%s",
        conf.Database.Username,
        conf.Database.DBName,
        conf.Database.SSLMode)
    
    // Initialize session
    session, err := data.GetDbSession(dsn, dialect.POSTGRES)
    if err != nil {
        panic(err)
    }
    defer session.Close()
    
    // Now you can use the session for database operations
}
```

### Define and Create Tables

```go
import (
    "fmt"
    "gormless/data"
    "gormless/data/dialect"
)

func InitUserTable(session data.ISession) error {
    // Define column types
    idType := dialect.PsqlSerial
    nameType := fmt.Sprintf(dialect.PsqlVarChar, 64)
    emailType := fmt.Sprintf(dialect.PsqlVarChar, 128)
    
    // Create table definition
    userTable := data.Table{
        Name: "users",
        Columns: &[]data.Column{
            {Name: "id", Type: &idType, PrimaryKey: true},
            {Name: "name", Type: &nameType},
            {Name: "email", Type: &emailType, Indexed: true},
        },
    }
    
    // Create the table
    return data.CreateTable(session, userTable)
}
```

### Working with Data

```go
// Define your data model
type User struct {
    dao      *data.DAO[User]
    ID       int
    Name     string
    Email    string
}

// Implement a method to get the DAO
func (u *User) GetDAO(session data.ISession) *data.DAO[User] {
    dao := data.DAO[User]{
        ISession: session,
        Table:    data.Table{Name: "users"},
    }
    return &dao
}

// Create a new user
func CreateUser(session data.ISession, name, email string) (*User, error) {
    user := &User{
        Name:  name,
        Email: email,
    }
    
    dao := user.GetDAO(session)
    err := dao.Upsert(*user)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

## Advanced Usage

### Creating Migrations

```go
// Define a migration to add a column
func AddUserStatusColumn() data.Migration {
    statusType := fmt.Sprintf(dialect.PsqlVarChar, 20)
    column := data.Column{
        Name: "status", 
        Type: &statusType,
    }
    
    return data.AddColumn(data.Table{Name: "users"}, column)
}

// Apply the migration
func ApplyMigration(session data.ISession) error {
    migration := AddUserStatusColumn()
    return migration(data.Table{Name: "users"}, session)
}
```

### Working with Foreign Keys

```go
func CreateTablesWithRelationship(session data.ISession) error {
    // First create the parent table
    roleIdType := dialect.PsqlSerial
    roleNameType := fmt.Sprintf(dialect.PsqlVarChar, 32)
    
    roleTable := data.Table{
        Name: "roles",
        Columns: &[]data.Column{
            {Name: "id", Type: &roleIdType, PrimaryKey: true},
            {Name: "name", Type: &roleNameType},
        },
    }
    
    err := data.CreateTable(session, roleTable)
    if err != nil {
        return err
    }
    
    // Now create the child table with a foreign key
    userIdType := dialect.PsqlSerial
    userNameType := fmt.Sprintf(dialect.PsqlVarChar, 64)
    roleIdFkType := dialect.PsqlInt
    
    userRoleFk := data.ForeignKey{
        Table:  &roleTable,
        Column: &data.Column{Name: "id"},
    }
    
    userTable := data.Table{
        Name: "users",
        Columns: &[]data.Column{
            {Name: "id", Type: &userIdType, PrimaryKey: true},
            {Name: "name", Type: &userNameType},
            {Name: "role_id", Type: &roleIdFkType, ForeignKey: &userRoleFk},
        },
    }
    
    return data.CreateTable(session, userTable)
}
```

## Dialect Support

gormless supports multiple SQL dialects through its dialect package. Currently supported dialects are:

- PostgreSQL
- MySQL

Each dialect implements the `Dialect` interface which provides methods for generating SQL specific to that database system.

## Testing

gormless includes unit tests and integration tests. To run the tests:

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run tests with coverage
make test-coverage
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.
