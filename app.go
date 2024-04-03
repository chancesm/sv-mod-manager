package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx     context.Context
	db      *sql.DB
	Configs map[string]string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		Configs: make(map[string]string),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	file := "test.db"
	a.ctx = ctx
	// Get Configurations
	db, err := sql.Open("sqlite", file)
	if err != nil {
		panic(err)
	}

	a.db = db
	runtime.EventsOn(a.ctx, "view:setupComplete", func(optionalData ...interface{}) {
		a.MustLoadOrInitDB()
		runtime.EventsEmit(a.ctx, "modFolder", a.Configs["modDir"])
		runtime.EventsEmit(a.ctx, "appReady")
	})
	// a.MustLoadOrInitDB()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) MustLoadOrInitDB() {
	// Create the database if it doesn't exist
	_, err := a.db.Exec(`
		CREATE TABLE IF NOT EXISTS appconfig (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);
	`)
	if err != nil {
		panic(err)
	}

	modDir := ""
	row := a.db.QueryRow("SELECT value FROM appconfig WHERE key = 'modDir';")
	err = row.Scan(&modDir)
	if err != nil {
		// If the entry doesn't exist, create it
		_, err = a.db.Exec("INSERT INTO appconfig (key, value) VALUES ('modDir', '');")
		if err != nil {
			panic(err)
		}
		modDir = ""
	}
	a.Configs["modDir"] = modDir
}
