# Shatkon ⬡ षटकोन

Shatkon is a Go package that helps you quickly scaffold a new Go project with customizable configurations for web frameworks, databases, and middleware.

https://github.com/user-attachments/assets/6942a82d-f626-45bf-8e18-a49bb00e0882

## Features

- Interactive CLI for project configuration
- Support for multiple Go web frameworks (Echo, Gin, Fiber, Chi, and standard library)
- Database integration options (MongoDB, PostgreSQL, SQLite)
- Logging middleware setup (for Echo framework)
- Automatic project structure creation
- Git repository initialization

## Installation

To install Shatkon, use the following command:

```bash
go get -u github.com/sarthak0714/shatkon
```

## Usage

After installation, you can use Shatkon to create a new project:

```bash
shatkon
```

Follow the interactive prompts to configure your project:

1. Enter your GitHub UserID
2. Choose a project name
3. Select a web framework
4. Choose a database
5. Enable or disable logging middleware

Shatkon will create a new directory with your project name and set up the basic structure and configuration files based on your choices.

## Project Structure

The generated project will have the following structure:

```
your-project-name/
├── cmd/
│   └── main.go
├── internal/
│   ├── adapters/
│   │   ├── handlers/
│   │   └── repository/
│   ├── config/
│   │   └── config.go
│   └── core/
│       ├── domain/
│       ├── ports/
│       └── services/
├── pkg/
│   └── utils/
│       └── logger.go (if logging is enabled)
├── go.mod
└── .git/
```

## Configuration

The project is configured using environment variables. Make sure to set the following variables before running your application:

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USERNAME`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

## Contributing

Contributions to Shatkon are welcome! Please feel free to submit a Pull Request.

