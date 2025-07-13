# 🚀 Helm Chart Browser

A beautiful, interactive terminal UI for browsing and downloading Helm chart values. Navigate through repositories, charts, and versions with ease!

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)

## ✨ Features

- 🎯 **Interactive Navigation** - Use arrow keys, vim keys (j/k), or number shortcuts
- 📊 **Beautiful Table Layout** - Clean, aligned columns for easy scanning
- 📄 **Smart Pagination** - Browse large lists with 10 items per page
- 🎨 **Color-coded Interface** - Visual hierarchy with syntax highlighting
- ⚡ **Fast & Responsive** - Async operations with loading states
- 🏷️ **Latest Version Badge** - Clearly identifies the newest chart version
- 💾 **Auto File Naming** - Downloads as `chartname-version-default-values.yaml`
- ⌨️ **Keyboard Shortcuts** - Full keyboard navigation support

## 🎬 Demo

```
🚀 Helm Chart Browser

🚀 Select a Helm repository:

     REPOSITORY           URL
──── ──────────────────── ───────────────────────────────────
► 1. argo                 https://argoproj.github.io/argo-helm
  2. external-secrets     https://charts.external-secrets.io
  3. apisix               https://charts.apiseven.com

📄 3 repositories available

⌨️  Navigate: ↑/↓ arrows or j/k • Select: Enter/Space or number (1-9,0) • Back: Backspace/Esc • Quit: q/Ctrl+C
💡 Tip: Use arrow keys to navigate through pages of results
```

## 📋 Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Helm CLI** - [Install Helm](https://helm.sh/docs/intro/install/)
- **Configured Helm Repositories** - Add repos with `helm repo add`

## 🚀 Quick Start

### Option 1: Download Release (Recommended)

```bash
# Download the latest release for your platform
curl -L https://github.com/tankibaj/helm-browser/releases/latest/download/helm-browser-linux -o helm-browser
chmod +x helm-browser
./helm-browser
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/tankibaj/helm-browser.git
cd helm-browser

# Build the application
go build -o helm-browser .

# Run it
./helm-browser
```

### Option 3: Install with Go

```bash
go install github.com/tankibaj/helm-browser@latest
helm-browser
```

## 🛠️ Development Setup

### 1. Clone and Setup

```bash
git clone https://github.com/tankibaj/helm-browser.git
cd helm-browser
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Run in Development Mode

```bash
go run .
```

### 4. Build for Production

```bash
# Build for current platform
go build -o helm-browser .

# Build for multiple platforms
make build-all
```

## 📦 Build Instructions

### Single Platform Build

```bash
go build -ldflags="-s -w" -o helm-browser .
```

### Cross-Platform Builds

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o helm-browser-linux .

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o helm-browser-darwin-amd64 .

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o helm-browser-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o helm-browser-windows.exe .
```

### Using Makefile

```bash
# Build for all platforms
make build-all

# Clean build artifacts
make clean

# Run tests
make test
```

## 🎮 Usage

### Navigation Controls

|Key                 |Action                      |
|--------------------|----------------------------|
|`↑/↓` or `j/k`      |Navigate up/down            |
|`Enter` or `Space`  |Select item                 |
|`1-9`, `0`          |Quick select (items 1-9, 10)|
|`Backspace` or `Esc`|Go back                     |
|`q` or `Ctrl+C`     |Quit application            |

### Workflow

1. **Start the application** - Automatically updates Helm repositories
1. **Select a repository** - Browse your configured Helm repos
1. **Choose a chart** - View all charts in the selected repository
1. **Pick a version** - See all available versions with app versions
1. **Download values** - Automatically saves `chartname-version-default-values.yaml`

### Example Session

```bash
$ ./helm-browser

# Navigate through:
# Repositories → Charts → Versions → Download

# Result:
# ✅ Successfully downloaded: argo-cd-5.46.8-default-values.yaml
```

## 🏗️ Architecture

### Key Components

- **Bubble Tea TUI** - Terminal user interface framework
- **Lipgloss Styling** - Beautiful colors and layouts
- **Helm CLI Integration** - Executes helm commands under the hood
- **Async Operations** - Non-blocking UI with loading states
- **State Management** - Clean state machine pattern

### Project Structure

```
helm-browser/
├── main.go           # Main application code
├── go.mod           # Go module dependencies
├── go.sum           # Dependency checksums  
├── README.md        # This file
├── LICENSE          # MIT license
├── Makefile         # Build automation
└── .github/
    └── workflows/   # CI/CD workflows
```

## 🧪 Testing

### Run Tests

```bash
go test ./...
```

### Test with Different Helm Setups

```bash
# Add test repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add stable https://charts.helm.sh/stable
helm repo update

# Run the application
./helm-browser
```

## 📊 Performance

- **Startup Time**: ~2-5 seconds (includes `helm repo update`)
- **Memory Usage**: ~10-20MB
- **Chart Search**: ~200-500ms per repository
- **Version Loading**: ~1-3s for charts with 1000+ versions
- **UI Responsiveness**: 60 FPS with async operations

## 🤝 Contributing

We welcome contributions! Here’s how to get started:

### 1. Fork & Clone

```bash
git clone https://github.com/tankibaj/helm-browser.git
cd helm-browser
```

### 2. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 3. Make Changes

- Follow Go best practices
- Add tests for new features
- Update documentation

### 4. Test Your Changes

```bash
go test ./...
go build .
./helm-browser
```

### 5. Submit Pull Request

- Write clear commit messages
- Include description of changes
- Reference any related issues

## 🐛 Troubleshooting

### Common Issues

**“helm command not found”**

```bash
# Install Helm first
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

**“Failed to list repos”**

```bash
# Add some repositories first
helm repo add stable https://charts.helm.sh/stable
helm repo update
```

**“No charts found”**

```bash
# Verify repositories are working
helm search repo --max-col-width=0
```

### Debug Mode

```bash
# Run with verbose output
HELM_DEBUG=true ./helm-browser
```

## 📝 License

This project is licensed under the MIT License - see the <LICENSE> file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Amazing TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Beautiful terminal styling
- [Helm](https://helm.sh/) - The package manager for Kubernetes
- [Charm](https://charm.sh/) - For creating delightful CLI tools

-----

**Made with ❤️ and Go**

*If you find this tool useful, please give it a ⭐ on GitHub!*