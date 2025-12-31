# UniFi Protect MCP

Model Context Protocol (MCP) server for Ubiquiti UniFi Protect management. Monitor and manage your UniFi security devices through an AI-powered interface.

**Focused on:** Cameras, motion sensors, smart lights, chimes, live view configurations, and security events.

## Features

- **14 management tools** for complete Protect operations
- **Camera Management**: Monitor and control UniFi cameras with PTZ and RTSPS support
- **Motion & Environment Sensors**: Real-time sensor data and detailed queries
- **Smart Lighting**: Control UniFi smart lights
- **Audio Alerts**: Manage door chimes and notifications
- **Live Views**: Configure custom camera view layouts
- **Advanced Camera Controls**: PTZ patrols, presets, talkback sessions
- **Event Monitoring**: Query security events and webhook integration
- **Viewer Management**: Manage Protect viewers and NVR systems
- **Stdio Transport**: MCP protocol over standard input/output for seamless integration

## Quick Start

### Installation

```bash
# Clone and build
git clone https://github.com/surrealwolf/unifi-protect-mcp.git
cd unifi-protect-mcp
go build -o bin/unifi-protect-mcp ./cmd
```

### Configuration

Create a `.env` file:

```bash
UNIFI_BASE_URL=https://your-unifi-controller.com:7442
UNIFI_API_KEY=your-api-key-here
UNIFI_SKIP_SSL_VERIFY=false
LOG_LEVEL=info
```

### Running the Server

```bash
./bin/unifi-protect-mcp
```

The server listens on stdio and is ready for MCP protocol messages.

## Available Tools (14 Total)

### Device Queries (6 tools)
- `get_protect_devices` - List all Protect devices
- `get_protect_cameras` - List all cameras
- `get_protect_sensors` - List all motion/environmental sensors
- `get_protect_lights` - List all smart lights
- `get_protect_chimes` - List all audio chimes
- `get_protect_liveviews` - List all configured live views

### Detailed Information (5 tools)
- `get_camera_detailed` - Get detailed camera information and settings
- `get_sensor_detailed` - Get detailed sensor information
- `get_light_detailed` - Get detailed light information and status
- `get_chime_detailed` - Get detailed chime information
- `get_liveview_detailed` - Get detailed live view configuration

### Events & Activity (1 tool)
- `get_protect_events` - Query security events with pagination

### Camera Controls (2 tools)
- `camera_create_rtsps_stream` - Create RTSPS video stream
- `camera_create_talkback_session` - Start two-way audio session

### PTZ Camera Control (Optional)
- `ptz_move_to_preset` - Move PTZ camera to saved preset
- `ptz_start_patrol` - Start automatic patrol sequence
- `ptz_stop_patrol` - Stop patrol sequence

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `UNIFI_BASE_URL` | UniFi Protect controller URL | Required |
| `UNIFI_API_KEY` | API key from UniFi controller | Required |
| `UNIFI_SKIP_SSL_VERIFY` | Skip SSL certificate verification | false |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | info |

## Usage with Claude/Copilot

When using this MCP with Claude or GitHub Copilot:

```
<mcp_server>
  <name>unifi-protect-mcp</name>
  <command>./path/to/bin/unifi-protect-mcp</command>
</mcp_server>
```

Then request Protect management tasks:
- "Show me all cameras and their current status"
- "Get recent security events from the last hour"
- "What motion sensors are currently active?"
- "List all configured live views"

For detailed AI assistant guidance, see [.github/copilot-instructions.md](.github/copilot-instructions.md).

## Docker Support

### Build Docker Image

```bash
docker build -t unifi-protect-mcp:latest .
```

### Run with Docker

```bash
docker run -e UNIFI_API_KEY=your-key -e UNIFI_BASE_URL=https://your-url unifi-protect-mcp:latest
```

### Docker Compose

```bash
# Create .env file with your configuration
cp .env.example .env
# Edit .env with your UniFi credentials

# Start the service
docker-compose up -d

# View logs
docker-compose logs -f
```

## GitHub Actions & CI/CD

This project includes automated workflows:

- **Tests**: Runs on every push and pull request
- **Docker Build**: Validates Dockerfile builds
- **Release**: Creates multi-platform binaries (Linux, macOS, Windows)
- **Auto-assign**: Assigns PRs to authors

See [.github/workflows](.github/workflows) for details.

## Skills & Capabilities

This MCP implements the following domain-specific skills:

- **Surveillance Monitoring**: Camera status, event review, coverage verification
- **Device Management**: Camera and sensor inventory, device monitoring
- **Security Management**: Access controls, system security configuration

See [.github/skills](.github/skills) for detailed skill documentation.

## Project Structure

```
unifi-protect-mcp/
├── cmd/
│   └── main.go              # Entry point and signal handling
├── internal/
│   ├── mcp/
│   │   └── server.go        # 14 MCP tool definitions and handlers
│   └── unifi/
│       ├── network.go       # Network API client (shared package)
│       ├── protect.go       # Protect API client
│       ├── doc.go           # Package documentation
│       └── client_test.go   # Integration tests
├── docs/
│   ├── API_REFERENCE.md     # Detailed API documentation
│   ├── GETTING_STARTED.md   # Setup guide
│   └── EXAMPLES.md          # Tool usage examples
├── bin/
│   └── unifi-protect-mcp    # Compiled binary
├── go.mod                   # Go module definition
├── go.sum                   # Dependency lock file
├── Makefile                 # Build and development tasks
└── .env.example             # Configuration template
```

## Development

### Building from Source

```bash
make build
```

### Running Tests

```bash
make test
```

### Cleaning Build Artifacts

```bash
make clean
```

## API Reference

For detailed API documentation, see [docs/API_REFERENCE.md](docs/API_REFERENCE.md).

For usage examples with specific tools, see [docs/EXAMPLES.md](docs/EXAMPLES.md).

## Troubleshooting

See [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) for common issues and solutions.

## Requirements

- Go 1.23.2 or later
- UniFi Protect system 3.x or later with API access enabled
- UniFi Controller 7.x or later

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please ensure:
- Code follows Go conventions
- All tests pass (`make test`)
- Changes are documented

## Related Projects

- **unifi-network-mcp**: MCP server for UniFi Network (WiFi, firewall, VPN, DPI)
- **UniFi Controller**: Official UniFi management software

---

**Version:** 0.1.0 | **Last Updated:** December 2025
