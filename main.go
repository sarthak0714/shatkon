package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type ProjectConfig struct {
	GithubUserID string
	ProjectName  string
	Framework    string
	Database     string
	Auth         bool
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

		// Part 4: Middleware Options
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Authentication Middleware?").
				Value(&config.Auth),
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
		"Authentication Middleware: %s\n"+
		"Logging Middleware: %s",
		titleStyle.Render("Project Configuration Summary"),
		keyword(config.GithubUserID),
		keyword(config.ProjectName),
		keyword(config.Framework),
		keyword(config.Database),
		keyword(fmt.Sprintf("%v", config.Auth)),
		keyword(fmt.Sprintf("%v", config.Logging)),
	)

	fmt.Println(lipgloss.NewStyle().
		Width(60).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(sb.String()))
}
