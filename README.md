# Fluxion

![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-MVP-orange)

**Fluxion** is an intelligent CLI tool designed to modernize legacy Java applications. It analyzes your source code, detects architecture patterns (Monolith vs Microservices), and automatically generates production-ready **Dockerfiles** and **Docker Compose** environments.

> Stop writing YAML manually. Let Fluxion handle the infrastructure.

## ⚡ Key Features

* **Smart Detection**: Automatically identifies Spring Boot versions, Java versions (8-21), and build tools.
* **Microservices Ready**: Detects Spring Cloud Gateway, Eureka, and Config Server relationships to generate a unified system config.
* **Infrastructure Analysis**: Scans dependencies to auto-configure required services:
    * 🐘 PostgreSQL / MySQL
    * 🧠 Redis / MongoDB
    * broker Kafka / RabbitMQ
    * 🔍 Elasticsearch / Consul
* **Production-Grade Dockerfiles**: Generates optimized Multi-Stage builds using Spring Boot **Layered Jars** for faster deployments.
* **Health-Checked Compose**: Generates `docker-compose.yml` with proper `depends_on` conditions and healthchecks, so your app waits for the DB to start.

## 🚀 Installation

### From Source
```bash
# Clone the repository
git clone [https://github.com/YOUR_USERNAME/fluxion.git](https://github.com/YOUR_USERNAME/fluxion.git)

# Build the binary
cd fluxion
go build -o fluxion main.go

# Run
./fluxion -path /path/to/your/java-project
```
(Binary releases coming soon)

## 📖 Usage
Simply point Fluxion to your project root (or a folder containing multiple microservices):

```bash
fluxion -path ./my-legacy-project
```
**What happens next?**

1. Fluxion scans the directory tree.
2. Identifies if it's a standalone app or a microservice mesh.
3. Gezerates Dockerfile in each service folder.
4. Generates a global docker-compose.yml with all dependencies wired up.
5. You run docker-compose up and enjoy.

## 🗺 Roadmap
- [x] MVP: Spring Boot & Docker Compose generation
- [ ] Support for Gradle projects
- [ ] Kubernetes manifests generation (Helm/Kustomize)
- [ ] CI/CD pipeline generation (GitHub Actions/GitLab CI)
- [ ] Interactive TUI (Terminal UI)

## 🤝 Contributing
We are looking for Go developers to join the team! If you are interested in DevOps, AST parsing, or CLI tools:
1. Fork the repository.
2. Create your feature branch (git checkout -b feature/AmazingFeature).
3. Commit your changes.
4. Open a Pull Request.

## 📄 License
Distributed under the MIT License. See LICENSE for more information.
