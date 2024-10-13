package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type ProjectConfig struct {
	GithubUserID string
	ProjectName  string
	Framework    string
	Database     string
	Logging      bool
}

func main() {
	var config ProjectConfig

	form := huh.NewForm(

		// user info
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your GitHub UserID").
				Description("This will be used to create the project repository.").
				Placeholder("johndoe").
				Value(&config.GithubUserID).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("GitHub UserID cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Enter your Project Name").
				Description("Choose a name for your new Go project.").
				Placeholder("my-awesome-project").
				Value(&config.ProjectName).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("project name cannot be empty")
					}
					return nil
				}),
		),

		// Framework Selection
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Go framework").
				Options(
					huh.NewOption("StdLib", "stdlib"),
					huh.NewOption("Gin", "gin"),
					huh.NewOption("Echo", "echo"),
					huh.NewOption("Fiber", "fiber"),
					huh.NewOption("Chi", "chi"),
				).
				Value(&config.Framework),
		),

		// Database Selection
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a database").
				Options(
					huh.NewOption("PostgreSQL", "postgresql"),
					huh.NewOption("MongoDB", "mongodb"),
					huh.NewOption("SQLite", "sqlite"),
				).
				Value(&config.Database),
		),

		// Middleware Options
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Logging Middleware?").
				Value(&config.Logging).
				Validate(func(b bool) error {
					if b && config.Framework != "echo" {
						return errors.New("logging middleware is only available for Echo framework")
					}
					return nil
				}),
		),

		// Confirmation
		huh.NewGroup(
			huh.NewConfirm().
				Title("Create this project?").
				Description("Review your choices and confirm to create the project."),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	InitProject(config)
	mainPath := config.ProjectName + "/cmd/main.go"

	switch config.Framework {
	case "stdlib":
		CreateFile(stdLibTemplate, mainPath)
	case "echo":
		if config.Logging {
			addEchoLogger(config)
			CreateFile(echoTemplateWithLogger, mainPath)
		} else {
			CreateFile(echoTemplate, mainPath)
		}
	case "gin":
		CreateFile(ginTemplate, mainPath)
	case "chi":
		CreateFile(chiTemplate, mainPath)
	case "fiber":
		CreateFile(fiberTempalte, mainPath)
	}

	dbFilepath := config.ProjectName + "/internal/adapters/repository/db.go"
	switch config.Database {
	case "sqlite":
		CreateFile(sqliteTemplate, dbFilepath)
	case "postgresql":
		CreateFile(pgSqlTemplate, dbFilepath)
	case "mongodb":
		CreateFile(mongoDBTemplate, dbFilepath)

	}

	goModCmd := exec.Command("go", "mod", "tidy")
	goModCmd.Dir = "./" + config.ProjectName
	if err := goModCmd.Run(); err != nil {
		panic(err)
	}

	printProjectSummary(config)
}

func printProjectSummary(config ProjectConfig) {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render(s)
	}

	fmt.Fprintf(&sb, "%s\n\n"+
		"GitHub UserID: %s\n"+
		"Project Name: %s\n"+
		"Framework: %s\n"+
		"Database: %s\n"+
		"Logging Middleware: %s",
		titleStyle.Render("Project Configuration Summary"),
		keyword(config.GithubUserID),
		keyword(config.ProjectName),
		keyword(config.Framework),
		keyword(config.Database),
		keyword(fmt.Sprintf("%v", config.Logging)),
	)
	fmt.Println(lipgloss.NewStyle().
		Width(60).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(sb.String()))
}

func InitProject(config ProjectConfig) error {
	if err := exec.Command("mkdir", config.ProjectName).Run(); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	dirs := []string{
		config.ProjectName + "/internal/adapters",
		config.ProjectName + "/internal/config",
		config.ProjectName + "/internal/core",
		config.ProjectName + "/internal/adapters/handlers",
		config.ProjectName + "/internal/adapters/repository",
		config.ProjectName + "/internal/core/domain",
		config.ProjectName + "/internal/core/ports",
		config.ProjectName + "/internal/core/services",
	}

	for _, dir := range dirs {
		if err := exec.Command("mkdir", "-p", dir).Run(); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	projectName := "github.com/" + config.GithubUserID + "/" + config.ProjectName

	goInitCmd := exec.Command("go", "mod", "init", projectName)
	goInitCmd.Dir = "./" + config.ProjectName
	if err := goInitCmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	gitInitCmd := exec.Command("git", "init")
	gitInitCmd.Dir = "./" + config.ProjectName
	if err := gitInitCmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}
	cfgFilePath := config.ProjectName + "/internal/config/config.go"

	if err := CreateFile(cfgTemplate, cfgFilePath); err != nil {
		return err
	}

	return nil

}

func CreateFile(content, filePath string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the plain string content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func addEchoLogger(cfg ProjectConfig) error {

	filePath := cfg.ProjectName + "/pkg/utils/logger.go"
	return CreateFile(loggerTemplate, filePath)
}

const loggerTemplate = `
package utils

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	colorRed       = "\033[31m"
	colorGreen     = "\033[32m"
	colorYellow    = "\033[33m"
	colorBlue      = "\033[34m"
	colorPurple    = "\033[35m"
	colorCyan      = "\033[36m"
	colorGray      = "\033[37m"
	colorReset     = "\033[0m"
	colorLightCyan = "\033[96m"
	colorMagenta   = "\033[35m"
)

// Returns color ASNII for the specified http status code
func statusColor(code int) string {
	switch {
	case code >= 100 && code < 200:
		return colorYellow
	case code >= 200 && code < 300:
		return colorGreen
	case code >= 300 && code < 400:
		return colorBlue
	case code >= 400 && code < 500:
		return colorRed
	case code >= 500:
		return colorPurple
	default:
		return colorReset
	}
}

// Custom Middleware function for Pretty logging :).
func CustomLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			logMessage := fmt.Sprintf("%s[%s]%s %s%s%s%s%s %s%s%s %s%s%d%s%s %s%v%s %s",
				colorLightCyan, time.Now().Format("2006-01-02 15:04:05"), colorReset,
				"\033[1m", colorGray, req.Method, colorReset, "\033[0m",
				colorCyan, req.URL.Path, colorReset,
				"\033[1m", statusColor(res.Status), res.Status, colorReset, "\033[0m",
				colorGray, time.Since(start), colorReset,
				id,
			)

			fmt.Println(logMessage)

			return nil
		}
	}
}

// Custom Middleware logger to indicate the perodic fetch afetr completion
func FetchLogger() {
	logMessage := fmt.Sprintf("%s[%s]%s %s%s%s%s%s",
		colorLightCyan, time.Now().Format("2006-01-02 15:04:05"), colorReset,
		"\033[1m", colorMagenta, "API FETCHED", colorReset, "\033[0m",
	)
	fmt.Println(logMessage)
}
`

const cfgTemplate = `
package config

type Config struct {}

func LoadConfig() *Config {
	return &Config{	}
}

`

const echoTemplate = `
package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
`

const echoTemplateWithLogger = `
package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner=true
	e.Use(utils.CustomLogger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
`

const chiTemplate = `
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("works"))
	})
	http.ListenAndServe(":8080", r)
}
`

const fiberTempalte = `
import (
    "log"

    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("works")
    })

    log.Fatal(app.Listen(":8080"))
}
`
const ginTemplate = `
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "works",
		})
	})
	r.Run() 
}
`

const stdLibTemplate = `
package main

import (
    "fmt"
    "net/http"
)



func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    		fmt.Fprintln(w, "Works")
		},
	)
    fmt.Println("Server is running at http://localhost:8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
`

const sqliteTemplate = `
package repositories

type sqliteDB struct {
	db *gorm.DB
}

// this will return a new sqlite struct
func NewStore(connectionString string) (*sqliteDB, error) {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &sqliteDB{
		db: db,
	}, nil
}

`

const pgSqlTemplate = `
package repository

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PGStore struct {
	db *gorm.DB
}

func NewStore(dsn string) (*PGStore, error) {
	// dsn := "host=localhost user=postgres dbname=postgres password=jomum port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PGStore{
		db: db,
	}, nil
}
`

const mongoDBTemplate = `
package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoStore(dsn string, dbName string) (*MongoStore, error) {
	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	db := client.Database(dbName)

	return &MongoStore{
		client: client,
		db:     db,
	}, nil
}

func (store *MongoStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return store.client.Disconnect(ctx)
}
`
