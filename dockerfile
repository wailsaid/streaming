FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o main -ldflags "-s -w" .

# Step 2: Create the final image
FROM alpine:3.18

# Set the working directory in the final image
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/templates /app/templates

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
