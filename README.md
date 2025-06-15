# Courses Service
[![Coverage Status](https://coveralls.io/repos/github/classconnect-grupo3/courses-service/badge.svg?branch=develop)](https://coveralls.io/github/classconnect-grupo3/courses-service?branch=develop)
## Overview
This is a Go-based microservice for managing courses. It provides endpoints to create, retrieve, and delete courses, storing data in a MongoDB database.

## Project Structure
- `src/`: Contains the source code of the project.
  - `model/`: Defines the data models (e.g., `Course`).
  - `schemas/`: Defines the request and response schemas (e.g., `CourseRequest` and `CourseResponse`).
  - `router/`: Contains the HTTP router (e.g., `router.go`).
  - `controller/`: Contains the HTTP controller (e.g., `CourseController`).
  - `service/`: Contains the business logic (e.g., `CourseService`).
  - `repository/`: Implements data access logic (e.g., `CourseRepository`).
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
- `GET /courses/title/{title}`: Retrieve courses by title
