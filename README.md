# ezazssh

A beautiful Terminal User Interface (TUI) for connecting to Azure VMs via SSH using passwordless authentication.

## Overview

`ezazssh` is a wrapper around the Azure CLI that provides an intuitive, interactive interface for establishing SSH connections to your Azure virtual machines. Built with [Charm Bracelet's Bubbletea](https://github.com/charmbracelet/bubbletea), it offers a smooth and elegant terminal experience.

## Features

- 🎨 **Beautiful TUI** - Clean, modern interface built with Bubbletea
- 🔐 **Passwordless SSH** - Leverages Azure CLI's passwordless authentication
- ⚡ **Fast & Intuitive** - Quickly browse and connect to your Azure VMs
- 🚀 **Portable** - Single binary with no dependencies

## Tips
- Press `W` to set your windows user if you are connecting to a windows VM
- While selecting your subscription, use `D` to save it as the default the next time you load

## Prerequisites

- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) installed and configured
- Active Azure subscription
- Authenticated Azure CLI session (`az login`)
- SSH extension installed `az extension add --name ssh`

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/burneystarke/ezazssh/releases):

```bash
# Linux (amd64)
curl -LO https://github.com/burneystarke/ezazssh/releases/latest/download/ezazssh_linux_amd64
chmod +x ezazssh_linux_amd64
sudo mv ezazssh_linux_amd64 /usr/local/bin/ezazssh

# Linux (arm64)
curl -LO https://github.com/burneystarke/ezazssh/releases/latest/download/ezazssh_linux_arm64
chmod +x ezazssh_linux_arm64
sudo mv ezazssh_linux_arm64 /usr/local/bin/ezazssh

# macOS (amd64)
curl -LO https://github.com/burneystarke/ezazssh/releases/latest/download/ezazssh_darwin_amd64
chmod +x ezazssh_darwin_amd64
sudo mv ezazssh_darwin_amd64 /usr/local/bin/ezazssh

# macOS (arm64/Apple Silicon)
curl -LO https://github.com/burneystarke/ezazssh/releases/latest/download/ezazssh_darwin_arm64
chmod +x ezazssh_darwin_arm64
sudo mv ezazssh_darwin_arm64 /usr/local/bin/ezazssh

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/burneystarke/ezazssh/releases/latest/download/ezazssh_windows_amd64.exe" -OutFile "ezazssh.exe"
```

### Build from Source

```bash
git clone https://github.com/burneystarke/ezazssh.git
cd ezazssh
go build -o ezazssh
```

## Requirements for Azure VMs

To use passwordless SSH with your Azure VMs, ensure:

- VM has WindowsOpenSSH or AADSSHLoginForLinux extension installed
- Your Azure AD account has appropriate VM access roles (e.g., "Virtual Machine Administrator Login" or "Virtual Machine User Login")
- Network Security Group allows SSH traffic (port 22)
- Windows VMs require a local username

## Troubleshooting

### SSH connection fails?
-Ensure you have installed the Azure CLI from [Microsoft's official guide](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli).
-Ensure you have installed the ssh extension `az extension add --name ssh`
-Try connecting via the cli `az ssh vm --help` for more info

### No items.
You may need to select a different subscription, or login to az cli with `az login`
