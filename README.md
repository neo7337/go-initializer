# go-initializer

Go Initializer is a modern web-based tool to quickly scaffold Go projects with your preferred project type, Go version, and framework/dependency. It features a beautiful React frontend and is designed for speed and ease of use.

## Features

- Select **Project Type**: Microservice, CLI App, API Server, or Simple Project
- Choose **Go Version**: 1.22.0 (latest stable), 1.21.7, 1.20.14
- Pick a **Framework/Dependency** based on project type (e.g., Gin, Echo, Fiber, Cobra, urfave/cli, etc.)
- Enter project metadata: module name, project name, and description
- One-click project generation
- Light/Dark theme toggle

## Usage

### Docker Setup (Recommended)

**Production Mode:**
```sh
docker-compose up --build
```

**Development Mode:**
```sh
docker-compose -f docker-compose.dev.yml up --build
```

Access the application at [http://localhost:3000](http://localhost:3000)

### Manual Setup

1. **Start the Backend**
	- Navigate to the `backend` directory:
	  ```sh
	  cd backend
	  go run .
	  ```

2. **Start the Frontend**
	- Navigate to the `frontend` directory:
	  ```sh
	  cd frontend
	  npm install
	  npm start
	  ```
	- Open [http://localhost:3000](http://localhost:3000) in your browser.

3. **Fill in the Form**
	- Select your desired project type, Go version, and framework/dependency.
	- Enter your module name, project name, and description.
	- Click **GENERATE** to scaffold your Go project.

## Project Structure

- `frontend/` – React app for the UI (see its README for dev scripts)
- `backend/` – (To be implemented) Go backend for project generation

## Framework/Dependency Options

Framework options are context-aware and change based on the selected project type:

| Project Type    | Framework/Dependency Options         |
|-----------------|--------------------------------------|
| Microservice    | Golly, Gin, Echo, Fiber, Go kit      |
| CLI App         | Golly, Cobra, urfave/cli, Kingpin    |
| API Server      | Golly, Gin, Echo, Fiber, Chi         |
| Simple Project  | Golly                                |

Support various addons to enhance your project:

| Addon           | Description                          |
|-----------------|--------------------------------------|
| Logging         | Zap, Logrus, Zerolog                 |
| Database        | GORM, SQLX, Ent                      |
| Cache           | Redis                                |
| Testing         | Testify, Ginkgo                      |
| Docker          | Dockerfile, Docker Compose           |
| CI/CD           | GitHub Actions, GitLab CI            |

## Contributing

Contributions are welcome! Please open issues or pull requests for suggestions and improvements.

## License

MIT
