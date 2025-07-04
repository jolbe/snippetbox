package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"snippetbox.gregor-pifko/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	user           *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:asdf@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	// Loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// DB connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize templates cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize new session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		infoLog:        infoLog,
		errorLog:       errorLog,
		snippets:       &models.SnippetModel{DB: db},
		user:           &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
	}

	srv := http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	infoLog.Printf("Starting server on %v", *addr)
	err = srv.ListenAndServeTLS("./tls/app.test+4.pem", "./tls/app.test+4-key.pem")
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
