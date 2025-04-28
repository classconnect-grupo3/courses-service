# Courses Service

## Overview
This is a Go-based microservice for managing courses. It provides endpoints to create, retrieve, and delete courses, storing data in a MongoDB database.

## Project Structure
- `src/`: Contains the source code of the project.
  - `model/`: Defines the data models (e.g., `Course`).
  - `repository/`: Implements data access logic (e.g., `CourseRepository`).
  - `handler/`: Contains HTTP handlers for the API endpoints.
  - `main.go`: Entry point of the application.

## Setup with Docker
1. Ensure you have Docker and Docker Compose installed on your system.
2. Clone the repository and navigate to the project directory.
3. Build and start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```
   This will start both the Go application and a MongoDB instance.

## API Endpoints
- `GET /courses`: Retrieve all courses.
- `POST /courses`: Create a new course.
- `GET /courses/{id}`: Retrieve a specific course by ID.
- `DELETE /courses/{id}`: Delete a specific course by ID.
- `GET /courses/teacher/{teacherId}`: Retrieve courses by teacher ID.
- `GET /courses/title/{title}`: Retrieve courses by title.
