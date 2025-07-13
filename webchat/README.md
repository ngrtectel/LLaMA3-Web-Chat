# LLaMA3-Web-Chat

A simple web-based chat application built with Go. This project serves static files and provides a basic chat interface using LLaMA3.

## Features
- Web chat interface
- Static file serving (HTML, CSS, JS)
- Go backend

## Getting Started

### Prerequisites
- Go 1.18 or newer

### Installation
1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd LLaMA3-Web-Chat/webchat
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

### Project Structure
```
webchat/
  go.mod        # Go module definition
  go.sum        # Go dependencies
  main.go       # Main server code
  static/
    index.html  # Chat frontend
```

## License
MIT
