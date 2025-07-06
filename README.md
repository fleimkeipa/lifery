# Lifery - Digital Life Diary

Lifery is a comprehensive full-stack web platform that enables users to record, organize, and share life memories and milestones in a rich multimedia diary format. Users can attach photos, videos, audio, and text to each memory and organize them chronologically on a timeline, grouped by life periods (e.g., "University Years," "First Job"). Each entry supports granular visibility settings (public, friends-only, or private), providing users with full control over their digital narrative.

<p align="center">
<img width="466" alt="Lifery Screenshot" src="https://github.com/user-attachments/assets/4ee1a138-0c73-4cf3-8180-93db9b412401" />
</p>

## 🌟 Key Features

- **Multimedia Memory Storage**: Support for photos, videos, audio, and text content
- **Timeline Organization**: Chronological organization with life period grouping
- **Granular Privacy Controls**: Public, friends-only, and private visibility settings
- **Social Connectivity**: Friend connections and social sharing capabilities
- **Multi-language Support**: Turkish and English interface
- **OAuth Integration**: Google and LinkedIn authentication
- **Real-time Notifications**: User activity and connection notifications
- **Responsive Design**: Modern, mobile-friendly interface

## 🏗️ Technical Architecture

### Backend (Go/Golang)

- **Framework**: Echo v4 for HTTP routing and middleware
- **Architecture**: Clean Architecture with clear separation of concerns
  - **Controllers**: HTTP request handling and response formatting
  - **Use Cases**: Business logic and application rules
  - **Repositories**: Data access layer with interface abstractions
  - **Models**: Data structures and validation
- **Database**: PostgreSQL with go-pg ORM
- **Caching**: Redis for performance optimization
- **Authentication**: JWT-based authentication with role-based access control
- **Documentation**: Swagger/OpenAPI 3.0 integration
- **Testing**: Comprehensive test suite with testcontainers for integration tests
- **Logging**: Structured logging with Zap
- **Validation**: Request validation with go-playground/validator

### Frontend (Vue.js/Nuxt.js)

- **Framework**: Nuxt 3 with Vue 3 Composition API
- **UI Library**: Nuxt UI with Tailwind CSS
- **State Management**: Vue composables for reactive state
- **Form Validation**: Vee-validate with Yup schemas
- **Internationalization**: Nuxt i18n with Turkish and English support
- **Authentication**: JWT token management with automatic refresh
- **File Upload**: Cloudinary integration for media storage
- **SSR/SSG**: Server-side rendering for improved SEO and performance

### DevOps & Deployment

- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose for local development
- **Cloud Platform**: Railway for production deployment
- **CI/CD**: GitHub integration for automated deployments
- **Environment Management**: Environment-specific configuration
- **Database**: PostgreSQL with connection pooling
- **Caching**: Redis for session and data caching

## 📁 Project Structure

```
lifery/
├── api/                    # Backend Go application
│   ├── controller/         # HTTP handlers and request/response logic
│   ├── model/             # Data models and validation schemas
│   ├── repositories/      # Data access layer with interfaces
│   ├── uc/               # Use cases (business logic)
│   ├── pkg/              # Shared utilities and configurations
│   ├── util/             # Helper functions and utilities
│   ├── tests/            # Test suites and test utilities
│   ├── docs/             # Swagger documentation
│   └── docker-compose.yaml # Local development environment
├── ui/                    # Frontend Nuxt.js application
│   ├── components/        # Vue components
│   ├── pages/            # Application pages and routing
│   ├── composables/      # Vue composables for state management
│   ├── middleware/       # Nuxt middleware
│   ├── layouts/          # Page layouts
│   └── i18n/             # Internationalization files
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.23+
- Node.js 18+
- Docker and Docker Compose
- PostgreSQL
- Redis

### Backend Setup

```bash
cd api
# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run with Docker Compose
docker-compose up -d

# Or run locally
go run main.go
```

### Frontend Setup

```bash
cd ui
# Install dependencies
npm install

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run development server
npm run dev
```

### API Documentation

Once the backend is running, access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## 🔧 Configuration

### Environment Variables

#### Backend (.env)

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=lifery
DB_USER=postgres
DB_PASSWORD=password

# JWT
JWT_SECRET=your-jwt-secret

# OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
LINKEDIN_CLIENT_ID=your-linkedin-client-id
LINKEDIN_CLIENT_SECRET=your-linkedin-client-secret

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

#### Frontend (.env)

```env
NUXT_API_BASE_URL=http://localhost:8080
NUXT_PUBLIC_CLOUDINARY_CLOUD_NAME=your-cloudinary-cloud-name
NUXT_PUBLIC_CLOUDINARY_UPLOAD_PRESET=your-upload-preset
```

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Role-based Access Control**: Admin, Editor, and Viewer roles
- **Password Hashing**: Secure password storage with bcrypt
- **CORS Configuration**: Proper cross-origin resource sharing
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries with go-pg
- **Rate Limiting**: API rate limiting (configurable)

## 🔐 OAuth Integration

Lifery supports multiple OAuth providers for seamless authentication:

### Supported Providers
- **Google OAuth 2.0**: Sign in with Google account
- **LinkedIn OAuth 2.0**: Sign in with LinkedIn account

### OAuth Flow
1. **Authorization Request**: User clicks "Sign in with [Provider]"
2. **Provider Redirect**: User is redirected to the OAuth provider
3. **User Consent**: User grants permission to Lifery
4. **Callback Handling**: Provider redirects back with authorization code
5. **Token Exchange**: Backend exchanges code for access token
6. **User Creation/Login**: User account is created or existing user is logged in
7. **JWT Token**: User receives JWT token for subsequent API calls

## 🌐 Internationalization

The application supports multiple languages:

- **Turkish (tr)**: Default language
- **English (en)**: Secondary language

Language switching is available through the UI, and all user-facing content is localized.

## 🚀 Deployment

### Production Deployment

The application is deployed on Railway with automated CI/CD from GitHub:

1. **Backend**: Containerized Go application with PostgreSQL and Redis
2. **Frontend**: Static site generation with Nuxt.js
3. **Database**: Managed PostgreSQL instance
4. **Caching**: Redis for session and data caching
5. **CDN**: Cloudinary for media file delivery

### Live Demo

- **API Documentation**: https://lifery-production.up.railway.app/swagger/index.html
- **Frontend Application**: https://lifery.bio

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:

- Create an issue in the GitHub repository
- Check the API documentation for technical details
- Review the test files for usage examples

---
