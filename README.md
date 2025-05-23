# Go Stack CLI 🚀

A powerful and easy-to-use CLI tool for automated deployment and management of your web projects.

---

## Overview

A powerful, lightweight deployment automation tool written in Go for streamlined Git-based deployments. `go-stack-cli` is perfect for automating server deployments with minimal configuration.

## ✨ Features

- **Git-based Deployments**: Automatically pull latest changes from your Git repository
- **Branch Management**: Switch between branches seamlessly during deployment
- **Smart Updates**: Only deploy when there are actual changes to pull
- **Verbose Logging**: Detailed output for debugging and monitoring
- **Force Reset**: Clean slate deployments with git reset functionality
- **Flexible Configuration**: Customizable working directories, branches, and remotes

## 🛠 Installation

> **Note:** You do **NOT** need to have Go installed to use this tool. Precompiled binaries are available in the [releases](https://github.com/satnamSandhu2001/go-stack-cli/releases) section for easy download and setup.

### Quick Install Script

```bash
curl -sSL https://raw.githubusercontent.com/satnamSandhu2001/go-stack-cli/master/install.sh | bash
```

### From Source

```bash
git clone https://github.com/satnamSandhu2001/go-stack-cli.git
cd go-stack-cli
go build -o go-stack-cli main.go
sudo mv go-stack-cli /usr/local/bin/
```

## 🚦 Usage

### Basic Usage

```bash
go-stack-cli
```

### Advanced Options

```bash
go-stack-cli [OPTIONS]

Options:
  -h                    Show help message
  --verbose             Show detailed output during execution
  --dir string          Root directory of project (default "./")
  --branch string       Git branch name to deploy (default "master")
  --git-remote string   Git remote name (default "origin")
  --git-reset           Force reset Git state before deployment
```

### Examples

**Deploy from a specific directory:**

```bash
go-stack-cli --dir /home/user/myproject
```

**Deploy specific branch with verbose output:**

```bash
go-stack-cli --branch production --verbose
```

**Force clean deployment:**

```bash
go-stack-cli --git-reset --verbose
```

**Deploy from custom remote:**

```bash
go-stack-cli --git-remote upstream --branch develop
```

## 📋 Prerequisites

- Git installed and configured
- Proper Git repository with remote configured
- Appropriate file permissions for the target directory

## 🔧 How It Works

1. **Validation**: Checks if the target directory is a valid Git repository
2. **Branch Management**: Switches to the specified branch if needed
3. **State Management**: Optionally resets Git state for clean deployments
4. **Update Check**: Compares local and remote commits to determine if updates are needed
5. **Deployment**: Pulls latest changes only when necessary
6. **Feedback**: Provides clear status updates throughout the process

## 🎯 Use Cases

- **Web Server Deployments**: Automate updates to web applications
- **CI/CD Integration**: Integrate with continuous deployment pipelines
- **Server Management**: Streamline multiple server deployments
- **Development Workflows**: Quick testing environment updates

## 🔮 Roadmap & Upcoming Features

### 🌐 Web Panel & Management Interface

- **Visual Dashboard**: Web-based interface for managing deployments
- **Real-time Monitoring**: Live deployment status and logs
- **Multi-project Management**: Handle multiple applications from one interface
- **User Authentication**: Secure access control and role management

### ⛓ Webhook Integration

- **GitHub/GitLab Webhooks**: Automatic deployments on push events
- **Custom Webhook Endpoints**: Flexible integration with various platforms
- **Payload Validation**: Secure webhook processing with signature verification
- **Conditional Deployments**: Deploy based on branch, tags, or commit messages

### ⚙️ Infrastructure Automation

- **PM2 Integration**: Automatic process management and restart
- **Nginx Configuration**: Dynamic virtual host and reverse proxy setup
- **Cloudflare Integration**: DNS and CDN configuration automation
- **SSL Certificate Management**: Automatic HTTPS setup with Let's Encrypt

### 🌍 Multi-Language & Framework Support

- **Python/Django**: Complete Django application deployment pipeline
- **Java/Spring Boot**: Spring Boot application deployment and management
- **Node.js/Express**: JavaScript application deployment automation
- **PHP/Laravel**: PHP application deployment with Composer integration
- **Frontend frameworks/ Static Sites**: Static sites and frameworks like React/angular/vue/Nextjs
- **Docker Support**: Containerized application deployment
- **Database Migrations**: Automatic database schema updates

### 🔧 Enhanced Features

- **Configuration Templates**: Predefined deployment configurations
- **Rollback Capabilities**: Quick rollback to previous deployments
- **Health Checks**: Application health monitoring post-deployment
- **Notification System**: Slack, Discord, email notifications
- **Deployment Scheduling**: Cron-like deployment scheduling
- **Environment Management**: Multiple environment support (dev, staging, prod)

## 🤝 Contributing

I welcome contributions! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
git clone https://github.com/satnamSandhu2001/go-stack-cli.git
cd go-stack-cli
go mod tidy
go run main.go -h
```

## 📝 Configuration Examples

## 🐛 Issues & Support

- **Bug Reports**: [GitHub Issues](https://github.com/satnamSandhu2001/go-stack-cli/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/satnamSandhu2001/go-stack-cli/discussions)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🏆 Acknowledgments

- Built with ❤️ using Go
- Inspired by the need for simple, reliable deployment automation
- Community-driven development and feedback

---

**⭐ Star this repository if you find it useful!**

_Go Stack CLI - Making deployments simple, one commit at a time._
