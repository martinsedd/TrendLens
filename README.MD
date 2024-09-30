# TrendLens

## Overview
TrendLens is a project designed to analyze and visualize trending topics across various platforms. The application aims to provide insights into popular discussions, helping users stay informed and engaged with current trends.

## Features
- Real-time trending topic analysis
- Data visualization for easy understanding
- User-friendly interface
- Integration with various APIs for data fetching

## Technologies Used
- ![Go](https://img.icons8.com/color/48/000000/golang.png) **Backend:** Go
- ![TypeScript](https://img.icons8.com/color/48/000000/typescript.png) **TypeScript**
- ![React](https://img.icons8.com/color/48/000000/react-native.png) **React**
- ![TailwindCSS](https://img.icons8.com/color/48/000000/tailwindcss.png) **TailwindCSS**
- ![MongoDB](https://img.icons8.com/color/48/000000/mongodb.png) **Database:** MongoDB
- ![Redis](https://img.icons8.com/color/48/000000/redis.png) **Caching:** Redis
- ![Docker](https://img.icons8.com/color/48/000000/docker.png) **Deployment:** Docker


## Getting Started

### Prerequisites
- Go 1.22+
- Node.js
- MongoDB
- Redis

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/martinsedd/TrendLens.git
   cd TrendLens
   ```

2. Install backend dependencies:
   ```bash
   cd api
   go mod tidy
   ```

3. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   ```

4. Configure environment variables:
    - Create a `.env` file in the `api` directory and specify your configuration settings (e.g., MongoDB URI, Redis settings).

### Running the Project

1. Start the backend:
   ```bash
   cd api
   go run main.go
   ```

2. Start the frontend:
   ```bash
   cd frontend
   npm start
   ```

Visit `http://localhost:3000` in your browser to see the application in action.


## Acknowledgments
- Thank you to the following Go libraries that are integral to this project:
   - ![Redis](https://img.icons8.com/color/48/000000/redis.png) `github.com/go-redis/redis/v8 v8.11.5`
   - ![Gorilla Mux](https://img.icons8.com/color/48/000000/gorilla.png) `github.com/gorilla/mux v1.8.1`
   - ![Dotenv](https://img.icons8.com/color/48/000000/environment.png) `github.com/joho/godotenv v1.5.1`
   - ![CORS](https://img.icons8.com/color/48/000000/cors.png) `github.com/rs/cors v1.11.1`
   - ![MongoDB](https://img.icons8.com/color/48/000000/mongodb.png) `go.mongodb.org/mongo-driver v1.17.0`

## Contact
For any inquiries, please contact:
- Email: martins.edd04@gmail.com
- GitHub: [martinsedd](https://github.com/martinsedd)
