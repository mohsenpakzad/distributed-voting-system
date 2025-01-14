# Distributed Voting System

This is a **Distributed Voting System** project built using **Golang** with **Gin** for RESTful services and **Kafka** for message queuing. The system is designed to handle voting operations in a distributed environment, ensuring scalability, fault tolerance, and consistency. The system consists of three main modules:

## Modules Overview

### 1. Vote Submitter

- **Description:**
  The `vote-submitter` module is a RESTful application built with **Gin** that receives votes from users.
- **Responsibilities:**
  - Handles user authentication.
  - Receives vote submissions.
  - Pushes the submitted votes into a **Kafka** queue.

### 2. Vote Validator

- **Description:**
  The `vote-validator` module consumes votes from the Kafka queue, validates them, and pushes the validated votes into another Kafka topic.
- **Responsibilities:**
  - Consumes votes from the queue.
  - Validates the votes based on predefined criteria (e.g., user eligibility, duplicate votes).
  - Pushes validated votes into a new Kafka topic for further processing.

### 3. Vote Counter

- **Description:**
  The `vote-counter` module is responsible for counting the validated votes using the **Raft consensus algorithm**.
- **Responsibilities:**
  - Consumes validated votes from the Kafka topic.
  - Ensures consistency in the vote counting process by leveraging Raft for distributed consensus.
  - Produces final vote counts.

## Features

- **Distributed Architecture:**
  - The system is composed of independent services that communicate asynchronously through Kafka.
- **Scalability:**
  - Each module can be scaled independently to handle varying loads.
- **Fault Tolerance:**
  - Kafka ensures message durability, and the Raft algorithm guarantees data consistency.
- **Modular Design:**
  - Each module is a standalone service, making it easier to develop, deploy, and maintain.

## Technologies Used

- **Golang:** For implementing the microservices.
- **Gin:** A high-performance HTTP web framework for building RESTful APIs.
- **Kafka:** For message queuing and asynchronous communication between services.
- **Raft Algorithm:** For distributed consensus and ensuring consistent vote counting.

## Getting Started

### Prerequisites

- **Golang** (v1.19+)
- **Kafka** (v2.8.0+)
- **Docker** (optional, for containerized deployment)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/distributed-voting-system.git
   cd distributed-voting-system
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Services

1. **Start Kafka**
   Ensure Kafka is running locally or in a Docker container.

2. **Run the Vote Submitter**

   ```bash
   cd vote-submitter
   go run main.go
   ```

3. **Run the Vote Validator**

   ```bash
   cd vote-validator
   go run main.go
   ```

4. **Run the Vote Counter**
   ```bash
   cd vote-counter
   go run main.go
   ```

## Deployment

- **Docker Compose:** A `docker-compose.yml` file can be used to orchestrate the deployment of all services.
- **Kubernetes:** For production-grade deployment, Kubernetes manifests can be created to manage the services and Kafka.
