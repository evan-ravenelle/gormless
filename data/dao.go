package data

import (
	"database/sql"
)

type IDAO interface {
	Upsert() error
	Get() (sql.Row, error)
	GetMultiple() error
	Delete() error
}

type DAO[T any] struct {
	ISession
	id     string
	Value  string
	Column Column
	Table  Table
}

func (dao *DAO[T]) Upsert(entity *T) error {
	query := dao.ISession.Dialect().Sprintd(
		"UPDATE %i SET %i = %s WHERE id = %s",
		dao.Table.Name,
		dao.Column.Name,
		dao.Value,
		dao.id)
	_, err := dao.ISession.Exec(query)
	return err
}

func (dao *DAO[T]) Get() error {
	query := dao.ISession.Dialect().Sprintd(
		"SELECT * FROM %i WHERE %i = %s",
		dao.Table.Name,
		dao.Column.Name,
		dao.id)
	err := dao.ISession.QueryRow(query).Scan(&dao.Value)
	return err
}

func (dao *DAO[T]) Delete() error {
	query := dao.ISession.Dialect().Sprintd(
		"DELETE FROM %i WHERE id = %s",
		dao.Table.Name,
		dao.id)
	_, err := dao.ISession.Exec(query, dao.id)
	return err
}
