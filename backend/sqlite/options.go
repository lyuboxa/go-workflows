package sqlite

import (
	"github.com/cschleiden/go-workflows/backend"
)

type options struct {
	*backend.Options

	// ApplyMigrations automatically applies database migrations on startup.
	ApplyMigrations bool

	MaxOpenConnections int

	SQLiteOptions []sqliteOption
}

type sqliteOption struct {
	key, value string
}

type option func(*options)

// WithApplyMigrations automatically applies database migrations on startup.
func WithApplyMigrations(applyMigrations bool) option {
	return func(o *options) {
		o.ApplyMigrations = applyMigrations
	}
}

// WithBackendOptions allows to pass generic backend options.
func WithBackendOptions(opts ...backend.BackendOption) option {
	return func(o *options) {
		for _, opt := range opts {
			opt(o.Options)
		}
	}
}

func WithSQLiteOption(name, setting string) option {
	return func(o *options) {
		o.SQLiteOptions = append(o.SQLiteOptions, sqliteOption{key: name, value: setting})
	}
}

// WithMaxOpenConnections sets the number of open connection per sqlite backend
func WithMaxOpenConnections(n int) option {
	return func(o *options) {
		o.MaxOpenConnections = n
	}

}
