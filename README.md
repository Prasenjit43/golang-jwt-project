```
# Golang JWT Token Project

This is a simple Golang project that demonstrates the usage of JSON Web Tokens (JWT) for authentication.

## Prerequisites

Before running this project, make sure you have the following installed:

- Go (version 1.13 or higher)
- Git

## Installation

To install and run this project locally, follow these steps:

1. Clone this repository:

```

git clone https://github.com/Prasenjit43/golang-jwt-project.git

```

2. Navigate to the project directory:

```

cd golang-jwt-project

```

3. Build and run the project:

```

go run main.go

```

## Usage

Once the project is running, you can access the endpoints using a tool like cURL or Postman. Here are the available routes:

- `POST /auth/login`: Endpoint for user authentication. Requires a valid username and password.
- `GET /user/profile`: Endpoint to fetch user profile. Requires a valid JWT token.

Make sure to set the appropriate environment variables such as `PORT` before running the project.

## Environment Variables

- `PORT`: The port on which the server will listen. Defaults to `8000` if not set.

## Libraries Used

- `github.com/gin-gonic/gin`: HTTP web framework for building web applications in Go.
- `github.com/Prasenjit43/golang-jwt-project/routes`: Custom routes for authentication and user management.

```
