# Docker Setup Guide

## Prerequisites
- Docker (Docker Desktop, Colima, Podman, etc.)
- Docker Compose

## Docker Runtime Compatibility

### Docker Desktop ✅
Works out of the box on Windows, macOS, and Linux.

### Colima (macOS) ✅
Requires network address flag:
```bash
colima start --network-address
```

### Podman ⚠️
Ensure rootless mode and port forwarding are configured.

### Rancher Desktop ✅
Works similarly to Docker Desktop.

## Common Issues

[Include the troubleshooting section from above]