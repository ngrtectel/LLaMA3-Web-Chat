# LLaMA3-Web-Chat

A simple web-based chat application built with Go. This project serves static files and provides a basic chat interface using LLaMA3.

## Features
- Web chat interface
- Static file serving (HTML, CSS, JS)
- Go backend
- Basic unit tests

## Getting Started

### Prerequisites
- Go 1.18 or newer

### Installation
1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd LLaMA3-Web-Chat
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application
```bash
go run main.go
```
The server will start and serve the chat interface at `http://localhost:8080` (default).

### Running Tests
```bash
go test
```

### Project Structure
```
LLaMA3-Web-Chat/
  go.mod         # Go module definition
  go.sum         # Go dependencies
  main.go        # Main server code
  main_test.go   # Unit tests
  static/
    index.html   # Chat frontend
  README.md      # Project documentation
```

## License
MIT
