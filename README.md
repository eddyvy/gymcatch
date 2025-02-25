# My App

This project is a full-stack application with a React frontend and a Go backend using the Fiber framework. The frontend is built with TypeScript and Vite, while the backend serves the static files and provides an API.

## Project Structure

```
my-app
├── backend          # Go backend application
│   ├── main.go     # Entry point for the Go application
│   ├── go.mod      # Go module definition
│   └── go.sum      # Dependency checksums
├── frontend         # React frontend application
│   ├── src
│   │   ├── App.tsx # Main React component
│   │   ├── main.tsx # Entry point for the React application
│   │   └── vite-env.d.ts # TypeScript definitions for Vite
│   ├── public
│   │   └── index.html # Main HTML file for the React app
│   ├── tsconfig.json # TypeScript configuration
│   ├── package.json  # npm configuration
│   └── vite.config.ts # Vite configuration
└── README.md        # Project documentation
```

## Getting Started

### Prerequisites

- Go (version 1.16 or later)
- Node.js (version 14 or later)
- npm (comes with Node.js)

### Setup Instructions

1. **Clone the repository:**

   ```
   git clone <repository-url>
   cd my-app
   ```

2. **Set up the backend:**

   Navigate to the `backend` directory and install the dependencies:

   ```
   cd backend
   go mod tidy
   ```

3. **Set up the frontend:**

   Navigate to the `frontend` directory and install the dependencies:

   ```
   cd frontend
   npm install
   ```

4. **Run the backend:**

   In the `backend` directory, run the Go application:

   ```
   go run main.go
   ```

5. **Run the frontend:**

   In the `frontend` directory, start the Vite development server:

   ```
   npm run dev
   ```

### Usage

- The frontend will be available at `http://localhost:3000` (or the port specified by Vite).
- The backend will serve the static files and can be accessed at `http://localhost:8080` (or the port specified in the Go application).

### Contributing

Feel free to submit issues or pull requests for any improvements or bug fixes. 

### License

This project is licensed under the MIT License.