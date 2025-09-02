# HTTP Server in Go

A lightweight custom HTTP server implemented in Go.
This project avoids the standard `net/http` package to provide a deeper understanding of how HTTP works under the hood by manually handling **requests**, **responses**, and **headers**.

---

## 🚀 Features

* Custom **request parser** (`RequestFromReader`)
* Custom **response writer** with status line, headers, and body
* Graceful shutdown on `SIGINT` / `SIGTERM`
* Example JSON response served from `test.json`
* Modular code structure with internal packages (`headers`, `request`, `response`, `server`)

---

## 📂 Project Structure

```
httpServer/
│── main.go
├── internal/
│   ├── headers/       # Handles HTTP headers
│   ├── request/       # Parses incoming HTTP requests
│   └── response/      # Builds and writes HTTP responses
└── server/
    └── server.go      # Core TCP server logic
```

---

## 🛠️ Installation & Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/httpServer.git
   cd httpServer/cmd/httpServer
   ```

2. Run the server:

   ```bash
   go run main.go
   ```

3. The server starts on port `42069`:

   ```bash
   curl http://localhost:42069
   ```

4. By default, it will try to serve the contents of `test.json`.

---

## 📝 Example Output

If you send a request:

```bash
curl -v http://localhost:42069
```

You’ll get a JSON response (from `test.json`) with headers like:

```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: <file_size>
```
