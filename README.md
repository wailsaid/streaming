# Simple Youtube Clone

This project is a simple Video streaming System built with the [Go](https://golang.org/) programming language and the [Gin](https://github.com/gin-gonic/gin) web framework. It demonstrates basic CRUD operations, dynamic HTML rendering, video streaming, and object storage using MinIO.

## Technologies Used

- [Go](https://golang.org/) - The primary programming language
- [GORM](https://gorm.io/) - An ORM library for Go, used for database operations
- [Plyr](https://plyr.io/) - A simple, accessible and customisable media player for Video
- [MinIO](https://min.io/) - High-performance object storage for video and thumbnail storage
- [Docker](https://www.docker.com/) - Containerization platform
- [MySQL](https://www.mysql.com/) - Database server

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install) (1.16 or later)

## Installation

1. **Clone the repository**:
   ```sh
   git clone https://github.com/wailsaid/streaming.git
   cd streaming
   ```

2. **Install dependencies**:
   ```sh
   go mod download
   ```

3. **Set up environment**:
   ```sh
   cp .env.example .env
   ```
   Configure the following environment variables in your .env file:
   - DB_HOST=db
   - DB_PORT=3306
   - DB_USER=root
   - DB_PASSWORD=rootpassword
   - DB_NAME=mydb
   
## Running with Docker

1. **Build and start the containers**:
   ```sh
   docker-compose up -d
   ```

2. **Access the application**:
   - Web Interface: http://localhost:8080
   - MinIO Console: http://localhost:9001
     - Username: minioadmin
     - Password: minioadmin

## Running Locally (Development)

1. **Run**:
   ```sh
   go run .
   ```
   Or build and run:
   ```sh
   go build
   ./streaming
   ```

## Features

- User authentication and authorization
- Video upload with thumbnail support
- Video streaming
- Object storage with MinIO
- Responsive video player
- Basic CRUD operations for videos

## Project Structure

```
.
├── controles/         # HTTP handlers and business logic
├── database/         # Database configuration and models
├── models/           # Data models
├── templ/           # HTML templates
├── utils/           # Utility functions and MinIO client
├── docker-compose.yml # Docker configuration
└── main.go          # Application entry point
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.