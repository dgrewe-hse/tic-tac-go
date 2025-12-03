## Tic-Tac-Go Game Server 

<p align="center">
  <a href="LICENSE"><img alt="Educational" src="https://img.shields.io/badge/status-educational-blue"></a>
  <a href="docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-available-brightgreen"></a>
  <a href="LICENSE"><img alt="License: Apache-2.0" src="https://img.shields.io/badge/license-Apache--2.0-blue"></a>
  <img alt="Labs" src="https://img.shields.io/badge/labs-12_planned-informational">
  <img alt="Made at HSE Esslingen" src="https://img.shields.io/badge/made%20at-HSE%20Esslingen-0a7ea4">
</p>

> IMPORTANT: This repository is for educational purposes only. It may contain unfinished, faulty, or even non-executable code used solely for teaching during the Software Testing course. Do not use this repository for any production system.

---

This repository contains a minimal Go backend for a Tic-Tac-Toe game server, designed to expose REST and WebSocket APIs for frontends such as Angular or Vue. This project is guided by the design document in `docs/spec.md`.

### Running the server

From the project root:

```bash
go run ./cmd/server
```

By default the server listens on port `8080`. You can override this via:

```bash
TICTACGO_PORT=9090 go run ./cmd/server
```

Once running, you can verify the basic health endpoint:

```bash
curl http://localhost:8080/health
```

You should receive:

```json
{"status":"ok"}
```

Further endpoints and game logic will be implemented in subsequent steps.

### Development

If you prefer to develop or run the project inside a preconfigured container, this repository includes a Dev Container setup under the `.devcontainer` directory.

- Using VS Code with Dev Containers extension installed:
  - Open the project folder in your editor.
  - When prompted, choose "Reopen in Container" (or use the Command Palette: `Dev Containers: Reopen in Container`).
  - The container image defined in `.devcontainer/Dockerfile` will be built, dependencies will be prepared, and you will be attached to a shell with Go tools preinstalled.
  - From inside the container, you can run the server the same way:

```bash
go run ./cmd/server
```

## License
This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

