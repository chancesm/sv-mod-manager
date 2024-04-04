package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path"

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
		runtime.EventsEmit(a.ctx, "app:modFolder", a.Configs["modDir"])
		runtime.EventsEmit(a.ctx, "app:appReady")
	})
	runtime.EventsOn(a.ctx, "view:changeFolder", func(optionalData ...interface{}) {
		dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
		if err != nil {
			panic(err)
		}
		if dir != "" {
			_, err = a.db.Exec("UPDATE appconfig SET value = ? WHERE key = 'modDir';", dir)
			if err != nil {
				panic(err)
			}
			a.Configs["modDir"] = dir
			a.AssertDisabledModsFolder()
			runtime.EventsEmit(a.ctx, "app:modFolder", dir)
		}
		slog.Info("Selected Folder", "dir", dir)
	})
	runtime.EventsOn(a.ctx, "view:openFolder", func(optionalData ...interface{}) {
		a.OpenModsFolder()
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
	if modDir != "" {
		a.AssertDisabledModsFolder()
	}
}
func (a *App) OpenModsFolder() {
	modsFolder := a.Configs["modDir"]
	cmd := "open"
	if runtime.Environment(a.ctx).Platform == "windows" {
		cmd = "explorer"
	}
	exec.Command(cmd, modsFolder).Start()
}
func (a *App) AssertDisabledModsFolder() {
	modsFolder := a.Configs["modDir"]
	disabledModsFolder := ".SVMM_DISABLED_MODS/"
	enabledModsFolder := "SVMM_ENABLED_MODS/"
	disabledModsPath := path.Join(modsFolder, disabledModsFolder)
	enabledModsPath := path.Join(modsFolder, enabledModsFolder)
	// Create the disabled mods folder if it doesn't exist in the modsFolder
	err := os.Mkdir(disabledModsPath, 0755)
	if err != nil {
		if !errors.Is(err, fs.ErrExist) {
			slog.Error("Failed to create disabled mods folder", "error", err)
		} else {
			slog.Info("Disabled mods folder already exists")
		}
	}
	// Create the enabled mods folder if it doesn't exist in the modsFolder
	err = os.Mkdir(enabledModsPath, 0755)
	if err != nil {
		if !errors.Is(err, fs.ErrExist) {
			slog.Error("Failed to create enabled mods folder", "error", err)
		} else {
			slog.Info("Enabled mods folder already exists")
		}
	}
}
