# Deployment Guide

This guide explains how to deploy the Tic-Tac-Go server using Docker.

## Prerequisites

- Docker installed on your Linux server
- Docker Compose (optional, for easier deployment)

## Quick Start with Docker Compose

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd tic-tac-go
   ```

2. **Build and run with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

3. **Verify the server is running:**
   ```bash
   curl http://localhost:8080/health
   ```

4. **View logs:**
   ```bash
   docker-compose logs -f
   ```

5. **Stop the server:**
   ```bash
   docker-compose down
   ```

## Manual Docker Deployment

### Build the Docker Image

```bash
# Build the image
docker build -t tic-tac-go-server:latest .

# Verify the image was created
docker images | grep tic-tac-go-server
```

### Run the Container

```bash
# Run on default port 8080
docker run -d \
  --name tic-tac-go-server \
  -p 8080:8080 \
  --restart unless-stopped \
  tic-tac-go-server:latest

# Or run on a custom port
docker run -d \
  --name tic-tac-go-server \
  -p 9090:9090 \
  -e TICTACGO_PORT=9090 \
  --restart unless-stopped \
  tic-tac-go-server:latest
```

### Container Management

```bash
# View logs
docker logs -f tic-tac-go-server

# Stop the container
docker stop tic-tac-go-server

# Start the container
docker start tic-tac-go-server

# Remove the container
docker rm tic-tac-go-server

# Remove the image
docker rmi tic-tac-go-server:latest
```

## Production Deployment

### Using Docker Compose with Custom Configuration

Create a `docker-compose.prod.yml` file:

```yaml
version: '3.8'

services:
  tic-tac-go-server:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: tic-tac-go-server
    ports:
      - "8080:8080"
    environment:
      - TICTACGO_PORT=8080
    restart: always
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
    # Optional: Resource limits
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

Deploy with:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Behind a Reverse Proxy (Nginx)

Example Nginx configuration:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support
        proxy_read_timeout 86400;
    }
}
```

### Using Systemd (Optional)

Create `/etc/systemd/system/tic-tac-go.service`:

```ini
[Unit]
Description=Tic-Tac-Go Server
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/path/to/tic-tac-go
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable tic-tac-go.service
sudo systemctl start tic-tac-go.service
```

## Health Check

The container includes a health check that verifies the `/health` endpoint:

```bash
# Check container health
docker ps  # Look for "healthy" status

# Manual health check
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

## Troubleshooting

### Container won't start

1. **Check logs:**
   ```bash
   docker logs tic-tac-go-server
   ```

2. **Verify port is not in use:**
   ```bash
   netstat -tuln | grep 8080
   # or
   ss -tuln | grep 8080
   ```

3. **Check Docker daemon:**
   ```bash
   docker info
   ```

### Port conflicts

If port 8080 is already in use, change the port:

```bash
docker run -d \
  --name tic-tac-go-server \
  -p 9090:9090 \
  -e TICTACGO_PORT=9090 \
  tic-tac-go-server:latest
```

### Permission issues

The container runs as a non-root user (UID 1000). If you encounter permission issues:

```bash
# Check container user
docker exec tic-tac-go-server id

# If needed, rebuild with different user
# (modify Dockerfile USER directive)
```

## Security Considerations

- ✅ Container runs as non-root user
- ✅ Minimal base image (Alpine Linux)
- ✅ Statically linked binary (no external dependencies)
- ✅ Health check included
- ⚠️ For production, consider:
  - Using HTTPS (via reverse proxy)
  - Implementing rate limiting
  - Adding authentication/authorization
  - Setting up firewall rules
  - Regular security updates

## Updating the Server

1. **Pull latest changes:**
   ```bash
   git pull
   ```

2. **Rebuild and restart:**
   ```bash
   docker-compose build
   docker-compose up -d
   ```

   Or with manual Docker:
   ```bash
   docker build -t tic-tac-go-server:latest .
   docker stop tic-tac-go-server
   docker rm tic-tac-go-server
   docker run -d --name tic-tac-go-server -p 8080:8080 tic-tac-go-server:latest
   ```

## Monitoring

### View resource usage

```bash
docker stats tic-tac-go-server
```

### View logs

```bash
# Follow logs
docker logs -f tic-tac-go-server

# Last 100 lines
docker logs --tail 100 tic-tac-go-server
```

## Backup

Since the server uses in-memory storage, there's no persistent data to backup. However, you should:

- Keep the Docker image in a registry
- Version control your configuration files
- Document any custom configurations

