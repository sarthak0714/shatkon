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
					huh.NewOption("MySQL", "mysql"),
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
		return err // Return the error if directory creation fails
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err // Return the error instead of panicking
	}
	defer file.Close()

	// Write the plain string content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err // Return the error instead of panicking
	}

	return nil
}

const cfgTemplate = `
package config

type Config struct {}

func LoadConfig() *Config {
	return &Config{	}
}

`
