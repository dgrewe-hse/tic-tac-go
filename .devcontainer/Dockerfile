FROM mcr.microsoft.com/devcontainers/go:1-bookworm

# Set the working directory inside the container
WORKDIR /workspace

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./

RUN go mod download

# Copy the rest of the source code
COPY . .

# Default command keeps container alive; editor will run `go run` as needed.
CMD [ "sleep", "infinity" ]


