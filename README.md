---

# Golang JWT Token Project

This is a simple project demonstrating how to implement JWT (JSON Web Tokens) authentication in a Golang web application using the Gin framework.

## Getting Started

### Prerequisites

- Go (Golang) installed on your machine.
- Understanding of basic Golang concepts.
- Basic knowledge of JWT (JSON Web Tokens).

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/Prasenjit43/golang-jwt-project.git
   ```

2. Navigate to the project directory:

   ```bash
   cd golang-jwt-project
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Usage

1. Start the server:

   ```bash
   go run main.go
   ```

2. By default, the server runs on port `8000`. You can specify a different port by setting the `PORT` environment variable:

   ```bash
   PORT=8080 go run main.go
   ```

3. Once the server is running, you can test the endpoints using tools like cURL or Postman.

## Endpoints

- `/auth`: Endpoint for authentication.
- `/user`: Endpoint for user authentication.
