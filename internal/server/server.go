// Package server provides the HTTP server for the application.
package server

import (
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"net/url"

	_ "github.com/throskam/memo/internal/translations"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/throskam/ki"
	"github.com/throskam/ki/middlewares"
	"github.com/throskam/kix/auth"
	"github.com/throskam/kix/i18n"
	"github.com/throskam/kix/sess"
	"github.com/throskam/memo/internal/lib"
	"github.com/throskam/memo/internal/orm"
	"golang.org/x/text/language"
)

func NewServer(ctx context.Context) *http.Server {
	// Config

	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	ki.SetLoggerLevelByText(config.Logger.Level)

	// Database

	pool, err := pgxpool.New(ctx, config.Database.ConnectionString)
	if err != nil {
		panic(err)
	}

	queries := orm.New(pool)

	// Session

	k := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}

	store := sess.NewSecureCookieSessionStore(
		config.Cookie.HashKey,
		config.Cookie.BlockKey,
	)

	// JWKs

	jwks := auth.JWKS{}

	jwks.Add(
		auth.JWK{
			Kid:   "1234",
			Value: config.JWT.Secret,
		},
	)

	// Services

	ts := lib.NewTopicService(queries)
	ps := lib.NewProjectService(queries)
	ams := lib.NewAuthenticationMethodService(queries)
	us := lib.NewUserService(queries, ams)

	// Router

	r := ki.NewRouter()

	r.Use(middlewares.RequestLogger())
	r.Use(middlewares.Recoverer())
	r.Use(middlewares.RequestID())
	r.Use(middlewares.RealIP())
	r.Use(middlewares.ContentCharset("utf-8"))
	r.Use(middlewares.ContentType("application/x-www-form-urlencoded"))
	r.Use(middlewares.ContentEncoding())
	r.Use(middlewares.ContentSecurityPolicy(url.Values{
		"default-src": {"'none'"},
		"script-src": {
			"'self'",
			"'unsafe-eval'",
			"https://unpkg.com",
			"https://cdn.jsdelivr.net",
		},
		"style-src": {
			"'self'",
			"https://fonts.googleapis.com",
			"https://cdn.jsdelivr.net",
			"'unsafe-inline'",
		},
		"font-src":        {"https://fonts.gstatic.com"},
		"connect-src":     {"'self'"},
		"img-src":         {"'self'", "data:"},
		"object-src":      {"'none'"},
		"base-uri":        {"'self'"},
		"frame-ancestors": {"'none'"},
	}))
	r.Use(middlewares.Locator(r))
	r.Use(middlewares.Language(language.MustParse("en-US"), language.MustParse("fr-FR")))
	r.Use(middlewares.OverrideLanguage("lang"))
	r.Use(i18n.Translator())

	addRoutes(
		r,
		store,
		ts,
		ps,
		us,
		jwks,
		config,
	)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	srv.RegisterOnShutdown(func() {
		pool.Close()
	})

	return srv
}