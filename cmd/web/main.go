package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"snippetbox.gregor-pifko/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
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

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %v", *addr)
	err = srv.ListenAndServe()
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
