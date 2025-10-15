# go-api-render

go-api-render is an API rendering module built on top of go-api-core. It provides a backend for testing API communication, user authentication (JWT), collection management, context, history, and includes an optional integrated frontend.

## Features

- JWT authentication and user management
- HTTP request execution and storage
- Import collections and OpenAPI specs
- System logs and metadata visualization
- Integrated SPA frontend (assets/front)
- TLS support and environment-based configuration

## Installation

```sh
git clone https://github.com/Rafael24595/go-api-render.git
cd go-api-render
go mod download
```

## Usage

```sh
go run main.go
```

## Configuration

Edit the .env file to set port, TLS, admin user, frontend, and other options.

## API Documentation

OpenAPI specification is available in swagger.yaml.
