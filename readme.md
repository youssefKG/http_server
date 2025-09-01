# HTTP Server in Go

A lightweight HTTP server implementation built from scratch in Go, featuring custom request parsing and graceful shutdown handling.

## 🚀 Features

- **Custom HTTP Server**: Built without using Go's standard `net/http` package
- **Request Parsing**: Custom request parser handling HTTP methods, headers, and body
- **Graceful Shutdown**: Proper signal handling for SIGINT and SIGTERM
- **Modular Architecture**: Clean separation of concerns with internal packages
- **Error Handling**: Custom error handling with status codes and messages

## 📁 Project Structure

```
httpServer/
├── main.go
├── cmd/
│   └── httpServer/
│       ├── internal/
│       │   └── request/
│       │       └── request.go
│       └── server/
│           └── server.go
└── README.md
```

## 🛠️ Components

### Main Server (`main.go`)
- Starts HTTP server on port `42069`
- Implements request handler with debugging output
- Handles graceful shutdown with OS signals

### Request Parser (`internal/request`)
- Parses HTTP request line (method, target, version)
- Extracts headers and body
- Provides structured access to request components

### Server Package (`server`)
- Core server implementation
- Handler interface with custom error handling
- Connection management and graceful shutdown

## 🏃‍♂️ Running the Server

1. **Clone and navigate to the project:**
   ```bash
   git clone <your-repo-url>
   cd httpServer
   ```

2. **Run the server:**
   ```bash
   go run main.go
   ```

3. **Server will start on port 42069:**
   ```
   Server started on port 42069
   ```

## 🧪 Testing

Test the server with curl:

```bash
# Simple GET request
curl http://localhost:42069

# POST request with data
curl -X POST http://localhost:42069 \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, Server!"}'

# Custom headers
curl http://localhost:42069 \
  -H "X-Custom-Header: test-value" \
  -H "User-Agent: MyClient/1.0"
```

## 📊 Example Output

When a request is received, the server logs:

```
GET / HTTP/1.1
Host: localhost:42069
User-Agent: curl/7.68.0
Accept: */*

Response: All good, frfr
```

## 🛑 Graceful Shutdown

Stop the server gracefully with `Ctrl+C` or `SIGTERM`:

```bash
^C
Server gracefully stopped
```

## 🔧 Configuration

- **Port**: Modify the `port` constant in `main.go` (default: 42069)
- **Handler**: Customize the `handler` function for different responses
- **Error Handling**: Modify `HandlerError` struct for custom error responses

## 📝 Handler Function

The handler function signature:

```go
func handler(w io.Writer, req *request.Request) *server.HandlerError
```

- **w**: Response writer
- **req**: Parsed HTTP request
- **Returns**: `nil` for success, or `*HandlerError` for errors

## 🎯 Future Enhancements

- [ ] Add routing support
- [ ] Implement middleware system
- [ ] Add HTTPS/TLS support
- [ ] Create configuration file support
- [ ] Add request/response logging
- [ ] Implement rate limiting
- [ ] Add authentication middleware

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built as a learning project to understand HTTP protocol internals
- Inspired by the desire to understand web servers from the ground up
