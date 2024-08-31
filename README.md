# Simple Youtube Clone

This project is a simple Video streaming System built with the [Go](https://golang.org/) programming language and the [Gin](https://github.com/gin-gonic/gin) web framework. It demonstrates basic CRUD operations, dynamic HTML rendering, and the use of GORM for database interactions.

## Technologies Used

- [Go](https://golang.org/) - The primary programming language.
- [Gin Web Framework](https://github.com/gin-gonic/gin) - A lightweight web framework for building APIs.
- [GORM](https://gorm.io/) - An ORM library for Go, used for database operations.
- [Plyr](https://www.google.com/url?sa=t&source=web&rct=j&opi=89978449&url=https://plyr.io/) - A simple, accessible and customisable media player for Video
## Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/wailsaid/streaming.git
   cd streaming

2. **Set Env Variables**:

   ```sh
   cp .env.example .env

set port and GIN_MODE ( debug or release ).

3. **RUN**:

   ```sh
   go run .

or build.