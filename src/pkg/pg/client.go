package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log/slog"
	"strings"
	"time"
	"tms/src/pkg/logger/sl"
)

type Config struct {
	URL         string
	AutoMigrate bool
	Migrations  string
}

type Client struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

// New создает новый экземпляр Client
func New(log *slog.Logger, cfg Config) (*Client, error) {
	const op = "Client.New"

	connString := fmt.Sprintf(cfg.URL)
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error("Error parsing pool config", slog.String("operation", op), sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	dbPool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Error("Error connecting to database", slog.String("operation", op), sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	p := &Client{
		db:  dbPool,
		log: log,
	}

	if cfg.AutoMigrate {
		if err = p.AutoMigrate(cfg.Migrations); err != nil {
			p.Close()
			log.Error("Error during auto migration", slog.String("operation", op), sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("Client intance created successfully")
	return p, nil
}

// Close закрывает подключение к базе данных
func (p *Client) Close() {
	p.db.Close()
	p.log.Info("Client connection closed")
}

// AutoMigrate выполняет автоматическую миграцию базы данных
func (p *Client) AutoMigrate(migrationsPath string) error {
	const op = "Client.AutoMigrate"

	migrationScript, err := ioutil.ReadFile(migrationsPath)
	if err != nil {
		p.log.Error("Failed to read migration file", slog.String("operation", op), sl.Err(err))
		return fmt.Errorf("%s: failed to read migration file: %w", op, err)
	}

	sqlCommands := strings.Split(string(migrationScript), ";")

	for _, cmd := range sqlCommands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		if _, err = p.Exec(context.Background(), cmd); err != nil {
			p.log.Error("Failed to execute migration command", slog.String("operation", op), slog.String("command", cmd), sl.Err(err))
			return fmt.Errorf("%s: failed to execute migration command: %w", op, err)
		}
	}

	p.log.Info("Auto migration completed successfully")
	return nil
}

// Exec обертка для выполнения SQL команд с логированием
func (p *Client) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	start := time.Now()

	result, err := p.db.Exec(ctx, sql, args...)
	duration := time.Since(start)

	if err != nil {
		p.log.Error("Error executing query",
			slog.String("query", formatSQLQuery(sql)),
			slog.Any("args", args),
			slog.Duration("duration", duration),
			sl.Err(err))
	} else {
		p.log.Info("Query executed successfully",
			slog.String("query", formatSQLQuery(sql)),
			slog.Any("args", args),
			slog.Duration("duration", duration))
	}
	return result, err
}

// Query обертка для выполнения SQL запросов с логированием
func (p *Client) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	start := time.Now()

	rows, err := p.db.Query(ctx, sql, args...)
	duration := time.Since(start)

	if err != nil {
		p.log.Error("Error executing query",
			slog.String("query", formatSQLQuery(sql)),
			slog.Any("args", args),
			slog.Duration("duration", duration),
			sl.Err(err))
	} else {
		p.log.Info("Query executed successfully",
			slog.String("query", formatSQLQuery(sql)),
			slog.Any("args", args),
			slog.Duration("duration", duration))
	}
	return rows, err
}

// QueryRow обертка для выполнения SQL запросов с логированием
func (p *Client) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	start := time.Now()

	row := p.db.QueryRow(ctx, sql, args...)
	duration := time.Since(start)

	p.log.Info("Query executed", slog.String("query", formatSQLQuery(sql)), slog.Any("args", args), slog.Duration("duration", duration))
	return row
}

// BeginTx начинает транзакцию
func (p *Client) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	start := time.Now()
	p.log.Info("Beginning transaction")

	tx, err := p.db.BeginTx(ctx, opts)
	if err != nil {
		p.log.Error("Error beginning transaction", slog.Duration("duration", time.Since(start)), sl.Err(err))
		return nil, err
	}

	p.log.Info("Transaction begun successfully", slog.Duration("duration", time.Since(start)))
	return tx, nil
}

// formatSQLQuery форматирует SQL запрос для лучшей читаемости в логах
func formatSQLQuery(query string) string {
	query = strings.ReplaceAll(query, "\n", " ")
	return strings.ReplaceAll(query, "\r", " ")
}
