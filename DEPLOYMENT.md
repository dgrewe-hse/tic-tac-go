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

2. **Clean up any existing containers (if you encounter errors):**
   ```bash
   # Remove existing container if it exists
   docker-compose down
   docker rm -f tic-tac-go-server 2>/dev/null || true
   ```

3. **Build and run with Docker Compose:**
   ```bash
   docker-compose up -d --build
   ```

   **Alternative:** If you're using newer Docker versions, use `docker compose` (without hyphen):
   ```bash
   docker compose up -d --build
   ```

   **Für externen Port 8081:** Bearbeiten Sie `docker-compose.yml` und ändern Sie:
   ```yaml
   ports:
     - "8081:8080"  # Extern 8081, intern 8080
   ```

4. **Verify the server is running:**
   ```bash
   curl http://localhost:8080/health
   # Oder bei Port 8081:
   curl http://localhost:8081/health
   ```

5. **View logs:**
   ```bash
   docker-compose logs -f
   # or
   docker compose logs -f
   ```

6. **Stop the server:**
   ```bash
   docker-compose down
   # or
   docker compose down
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
# Run on default port 8080 (extern und intern)
docker run -d \
  --name tic-tac-go-server \
  -p 8080:8080 \
  --restart unless-stopped \
  tic-tac-go-server:latest

# Extern auf Port 8081, intern weiterhin auf 8080
# Syntax: -p EXTERNER_PORT:INTERNER_PORT
docker run -d \
  --name tic-tac-go-server \
  -p 8081:8080 \
  --restart unless-stopped \
  tic-tac-go-server:latest

# Oder extern auf Port 9090, intern auf 8080
docker run -d \
  --name tic-tac-go-server \
  -p 9090:8080 \
  --restart unless-stopped \
  tic-tac-go-server:latest

# Wenn Sie den internen Port auch ändern möchten (nicht empfohlen)
docker run -d \
  --name tic-tac-go-server \
  -p 9090:9090 \
  -e TICTACGO_PORT=9090 \
  --restart unless-stopped \
  tic-tac-go-server:latest
```

**Wichtig:** 
- Die Syntax ist `-p EXTERNER_PORT:INTERNER_PORT`
- Der Container läuft **intern** immer auf Port 8080 (oder dem Wert von `TICTACGO_PORT`)
- Der **externe Port** kann beliebig gewählt werden (z.B. 8081, 9090, etc.)
- Von außen greifen Sie dann über `http://localhost:EXTERNER_PORT` zu

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
      - "8080:8080"  # Format: "EXTERNAL:INTERNAL" - z.B. "8081:8080" für externen Port 8081
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

**Für externen Port 8081 ändern Sie die ports-Zeile:**
```yaml
ports:
  - "8081:8080"  # Extern 8081, intern 8080
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

### Docker Compose Errors

**If you encounter `KeyError: 'ContainerConfig'` or similar errors:**

1. **Clean up existing containers and images:**
   ```bash
   docker-compose down
   docker rm -f tic-tac-go-server 2>/dev/null || true
   docker rmi tic-tac-go-server:latest 2>/dev/null || true
   docker rmi tic-tac-go_tic-tac-go-server:latest 2>/dev/null || true
   ```

2. **Rebuild from scratch:**
   ```bash
   docker-compose build --no-cache
   docker-compose up -d
   ```

3. **Alternative: Use Docker directly instead of docker-compose:**
   ```bash
   docker build -t tic-tac-go-server:latest .
   docker run -d --name tic-tac-go-server -p 8080:8080 tic-tac-go-server:latest
   ```

4. **If using newer Docker (v20.10+), use `docker compose` (without hyphen):**
   ```bash
   docker compose down
   docker compose up -d --build
   ```

### Port conflicts

If port 8080 is already in use, you have two options:

**Option 1: Map external port to internal port 8080 (empfohlen)**
```bash
# Extern auf Port 8081, intern auf 8080
docker run -d \
  --name tic-tac-go-server \
  -p 8081:8080 \
  --restart unless-stopped \
  tic-tac-go-server:latest

# Dann über http://localhost:8081 erreichbar
curl http://localhost:8081/health
```

**Option 2: Ändern Sie auch den internen Port (nicht empfohlen)**
```bash
# Extern und intern auf Port 9090
docker run -d \
  --name tic-tac-go-server \
  -p 9090:9090 \
  -e TICTACGO_PORT=9090 \
  --restart unless-stopped \
  tic-tac-go-server:latest
```

**Empfehlung:** Verwenden Sie Option 1, da der Container standardmäßig auf Port 8080 konfiguriert ist.

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

