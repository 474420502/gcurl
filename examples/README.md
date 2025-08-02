# Examples Directory

This directory contains practical examples demonstrating various features of gcurl.

## Directory Structure

Each example is organized in its own subdirectory with a `main.go` file that can be run independently:

### 📂 `advanced_networking/`
Demonstrates advanced networking features including:
- `--connect-to` for connection redirection
- `-G/--get` for converting POST data to query parameters
- Complex integration scenarios

**Run the example:**
```bash
cd advanced_networking
go run main.go
```

### 📂 `authentication_demo/`
Demonstrates various authentication methods:
- HTTP Basic Authentication
- Bearer Token Authentication (JWT, OAuth2)
- API Key Authentication
- Complex multi-factor authentication scenarios

**Run the example:**
```bash
cd authentication_demo
go run main.go
```

### 📂 `readme_examples/`
Contains all the examples from the main README.md file:
- Basic GET requests
- POST with form data
- File uploads
- Cookie handling
- And more...

**Run the example:**
```bash
cd readme_examples
go run main.go
```

## Building Examples

Each example can be built as a standalone executable:

```bash
# Build advanced networking example
cd advanced_networking && go build -o advanced_networking .

# Build authentication demo
cd authentication_demo && go build -o authentication_demo .

# Build README examples
cd readme_examples && go build -o readme_examples .
```

## Go Module Structure

Each example directory uses the shared `go.mod` and `go.sum` files from the parent `examples` directory, ensuring consistent dependency management across all examples.
