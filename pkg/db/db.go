package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Исполняемая функция в транзакции
type Handler func(ctx context.Context) error

// Client для работы с БД
type Client interface {
	DB() DB
	Close() error
}

type TxManager interface {
	ReadCommited(ctx context.Context, f Handler) error
}

// Обертка над запросом с именем запроса и самим запросом
type Query struct {
	Name     string
	QueryRow string
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// NamedExecer + QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// Интерфейс для рботы с именнованными запросами с помощью тегов в структурах
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// Интерфейс для работы с обычными запросами
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecer
	Pinger
	Transactor
	Close()
}
