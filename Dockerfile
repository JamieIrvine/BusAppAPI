# Use official Go image with necessary tools
FROM golang:1.24.5

# Set working directory inside the container
WORKDIR /app

# Default command (overridable in docker-compose)
CMD ["go", "run", "."]
