@"
# Supply Chain Intelligence Platform (ScIntel) - Backend

> **AI-Powered Demand Forecasting & Warehouse Operations Management**

## 📖 Overview
ScIntel is a production-grade **Golang backend system** designed for multi-warehouse inventory management, demand forecasting (ETS & ARIMA), real-time anomaly detection, and event-driven webhooks. Built with Clean Architecture principles, it serves as a portfolio-grade enterprise backend.

## 🚀 Quick Start
### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make (optional but recommended)

### Running Locally
1. Clone the repository
2. Copy environment variables:
   `cp .env.example .env.local`
3. Start infrastructure:
   `make docker-up`
4. Run database migrations:
   `make db-migrate`
5. Run the application:
   `make run`

## 🛠 Technology Stack
| Component | Technology |
|-----------|------------|
| **Language** | Go 1.21+ |
| **Framework** | Standard Library + net/http |
| **Database** | PostgreSQL 15+ |
| **Cache** | Redis 7+ |
| **Containerization** | Docker & Docker Compose |
| **CI/CD** | GitHub Actions |

## 🧬 Architecture
- **Clean Architecture**: Handler → Service → Repository
- **Authentication**: JWT (Access: 1h, Refresh: 7d)
- **Authorization**: RBAC (Admin, Manager, Analyst, Viewer)
- **Forecasting**: ETS & ARIMA time-series algorithms
- **Caching**: Redis (30s inventory TTL, 24h forecast TTL)

## 📁 Project Structure
\`\`\`
scinel-backend/
├── cmd/api/          # Entry point
├── internal/         # Private business logic
│   ├── handler/      # HTTP layer
│   ├── service/      # Business logic
│   └── repository/   # Database layer
├── migrations/       # SQL migrations
└── pkg/              # Reusable libraries
\`\`\`

## 🤝 Contributing
Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 🔒 Security
Please read [SECURITY.md](SECURITY.md) for reporting vulnerabilities.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author
**M. Arif Aulia** - Backend Engineer

---
*Built for portfolio excellence. 10/10 recruiter rating target.*
"@ | Out-File -FilePath README.md -Encoding utf8