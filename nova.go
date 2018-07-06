package orm // import "github.com/novakit/orm"

import (
	"net/url"
	"os"

	"github.com/novakit/nova"
)

// ContextKey key in nova.Context
const ContextKey = "nova.orm"

// Options options structure
type Options struct {
	// Dedicated create dedicated orm.DB per request
	Dedicated bool
	// Dialect sql driver dialect
	Dialect string
	// URL database url
	URL string
}

func sanitizeOptions(opts ...Options) (opt Options) {
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.URL) == 0 {
		opt.URL = os.Getenv("DATABASE_URL")
	}
	if len(opt.Dialect) == 0 {
		opt.Dialect = os.Getenv("DATABASE_DIALECT")
	}
	if len(opt.Dialect) == 0 {
		if l, err := url.Parse(opt.URL); err != nil {
			opt.Dialect = l.Scheme
		}
	}
	return
}

// Handler create a nova.HandlerFunc injects *DB
func Handler(opts ...Options) nova.HandlerFunc {
	opt := sanitizeOptions(opts...)
	var db *DB
	var err error
	// create instance if not dedicated
	if !opt.Dedicated {
		if db, err = Open(opt.Dialect, opt.URL); err != nil {
			// panic if failed
			panic(err)
		}
	}
	return func(ctx *nova.Context) (err error) {
		// create instance if dedicated
		if opt.Dedicated {
			if db, err = Open(opt.Dialect, opt.URL); err != nil {
				return
			}
		}
		// inject to ctx
		ctx.Values[ContextKey] = db
		// invoke next handler
		ctx.Next()
		// close instance if dedicated
		if opt.Dedicated {
			db.Close()
		}
		return
	}
}

// Extract extract previous injected *DB
func Extract(ctx *nova.Context) (db *DB) {
	db, _ = ctx.Values[ContextKey].(*DB)
	return
}
