package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"

	"huyduong.com/snippetbox/internal/models"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "snippetbox:huyduong1401@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

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

	infoLog.Printf("Starting server on %s", *addr)

	srv := http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
