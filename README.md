# StackJet 🚀

A powerful and easy-to-use CLI + UI tool for multistack automated deployment and management of your web projects.

---

## Overview

A powerful, lightweight deployment automation tool written in Go for streamlined Git-based deployments. `stackjet` is perfect for automating server deployments with minimal configuration and aims to support multiple technology stacks including Node.js, Golang, Java, Python and PHP.

## ✨ Features

- **Git-based Deployments**: Automatically pull latest changes from your Git repository
- **Branch Management**: Switch between branches seamlessly during deployment
- **Smart Updates**: Only deploy when there are actual changes to pull
- **Custom Commands**: Configurable build, start, and post-deployment commands
- **Verbose Logging**: Detailed output for debugging and monitoring
- **Force Reset**: Clean slate deployments with git reset functionality
- **Flexible Configuration**: Customize StackJet's behavior based on your needs and preferences

## 🛠 Installation

> **Note:** You do **NOT** need to have Go installed to use this tool. Precompiled binaries are available in the [releases](https://github.com/satnamSandhu2001/StackJet/releases) section for easy download and setup.

### Quick Install Script

```bash
curl -sSL https://raw.githubusercontent.com/satnamSandhu2001/StackJet/master/install.sh | bash
```

## 🚦 Usage

### Initialize StackJet

Before using StackJet, initialize it in your environment:

```bash
stackjet init

Options:
  -f, --force    Force recreate config (Use with caution! This will overwrite any existing config)
  -h, --help     Show help message
```

### Add New Application

Add a new application to StackJet for deployment management:

```bash
stackjet add [OPTIONS]

Required Options:
  -t, --tech string       App's Technology Stack Type (currently supports: nodejs)
  -r, --repo string       Git repository URL
  -p, --port int          Port number for the application

Optional Options:
  --branch string         Git branch name (default from config)
  --git-remote string     Git remote name (default from config)
  --build string          Build commands (e.g., 'npm install && npm run build')
  --start string          App start commands (e.g., 'npm start') default is 'npm start'
  --post string           Post deployment commands (e.g., 'npm run post-deploy')
  -h, --help              Show help message
```

#### Examples for Add Command

**Add a Node.js application:**

```bash
stackjet add --tech nodejs -p 3000 --repo https://github.com/username/my-app.git
```

**Add with custom build and start commands:**

```bash
stackjet add --tech nodejs -p 8080 --repo https://github.com/username/my-app.git \
  --build "npm install && npm run build" \
  --start "npm run prod" \
  --post "npm run migrate"
```

**Add with specific branch:**

```bash
stackjet add --tech nodejs -p 3000 --repo https://github.com/username/my-app.git \
  --branch production --git-remote upstream
```

### Deploy Application

Deploy your application with Git sync, process management, and more:

```bash
stackjet deploy [OPTIONS]

Options:
  -d, --dir string        Root directory of project (default "./")
  --branch string         Git branch name to deploy
  --git-remote string     Git remote name
  --git-hash string       Rollback to specific commit hash
  --git-reset             Force reset Git state before deployment (default true)
  -h, --help              Show help message
```

#### Examples for Deploy Command

**Deploy from current directory:**

```bash
stackjet deploy
```

**Deploy from specific directory:**

```bash
stackjet deploy --dir "/var/www/sites/<You will recieve directory path when you add new app to stackjet>"
```

**Deploy specific branch:**

```bash
stackjet deploy --branch "production"
```

**Rollback to specific commit:**

```bash
stackjet deploy --git-hash "abc123def456"
```

## 🔧 Technology Stack Support

### Node.js Applications

StackJet provides comprehensive support for Node.js applications with PM2 integration:

- **Automatic PM2 Setup**: Process management with PM2 for production deployments
- **Custom Start Commands**: Support for various Node.js start commands
- **Build Process**: Configurable build commands for compilation and optimization
- **Post-Deployment Hooks**: Execute custom commands after deployment

**Supported Node.js Commands:**
Start commands must be in the format: `[npm|yarn|pnpm] [<script name>]` without any extra args. args must be embedded in package file.

For example:

- `npm start`
- `npm run prod`
- `yarn dev:server`
- `npm start:prod`

**Default Behavior:**

- If no start command is specified, defaults to `npm start`
- Build commands are optional and executed before starting the application
- Post commands run after successful deployment

### Upcoming Stack Support

- **Java/Spring Boot**: Spring Boot application deployment and management
- **Python/Django**: Complete Django application deployment pipeline
- **PHP/Laravel**: PHP application deployment with Composer integration
- **Static Sites**: Static site deployments
- **Docker Support**: Containerized application deployment

## 📋 Prerequisites

- Git installed and configured
- (For Node.js applications only) Node.js, [npm|yarn|pnpm] and [PM2](https://pm2.keymetrics.io/docs/usage/quick-start) installed

## 🔧 How StackJet Works

### Application Addition Flow

1. **Stack Validation**: Validates the specified technology stack
2. **Repository Setup**: Validates Git repository URL and accessibility
3. **Port Management**: Checks port availability and validates port numbers
4. **Command Validation**: Validates start commands for the specified stack
5. **Database Storage**: Stores application configuration for future deployments

### Deployment Flow

1. **Configuration Loading**: Loads application configuration from database
2. **Git Operations**: Handles branch switching, pulling, and reset operations
3. **Build Process**: Executes build commands if specified
4. **Process Management**: Manages application processes (PM2 for Node.js)
5. **Health Checks**: Verifies successful deployment
6. **Post-Deployment**: Executes post-deployment commands

## 🎯 Use Cases

- **Node.js Web Applications**: Deploy Express.js, React.js, Angular.js, Vue.js Koa.js, or any Node.js based applications
- **API Services**: Deploy REST APIs and microservices
- **Full-Stack Applications**: Deploy complete web applications with frontend and backend
- **CI/CD Integration**: Integrate with continuous deployment pipelines
- **Multi-Environment Deployments**: Manage development, staging, and production deployments

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

- **Enhanced PM2 Integration**: Advanced PM2 configuration and monitoring
- **Nginx Configuration**: Dynamic virtual host and reverse proxy setup
- **Cloudflare Integration**: DNS and CDN configuration automation
- **SSL Certificate Management**: Automatic HTTPS setup with Let's Encrypt

### 🌍 Enhanced Multi-Language Support

- **Python/Django**: Complete Django application deployment pipeline
- **Java/Spring Boot**: Spring Boot application deployment and management
- **PHP/Laravel**: PHP application deployment with Composer integration
- **Frontend Frameworks**: Static sites and SPAs (React, Angular, Vue, Next.js)
- **Docker Support**: Containerized application deployment
- **Database Migrations**: Automatic database schema updates

### 🔧 Enhanced Features

- **Configuration Templates**: Predefined deployment configurations
- **Advanced Rollback**: Quick rollback to previous deployments with history
- **Health Checks**: Application health monitoring post-deployment
- **Notification System**: Slack, Discord, email notifications
- **Deployment Scheduling**: Cron-like deployment scheduling
- **Environment Management**: Multiple environment support (dev, staging, prod)
- **Load Balancing**: Automatic load balancer configuration

## 🤝 Contributing

I welcome contributions! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
git clone https://github.com/satnamSandhu2001/StackJet.git
cd stackjet
go mod tidy
go run main.go -h
```

## 🐛 Issues & Support

- **Bug Reports**: [GitHub Issues](https://github.com/satnamSandhu2001/StackJet/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/satnamSandhu2001/StackJet/discussions)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🏆 Acknowledgments

- Built with ❤️ using Go
- Inspired by the need for simple, reliable deployment automation
- Community-driven development and feedback

---

**⭐ Star this repository if you find it useful!**

_StackJet - Deploy Fast, Fly High!_
