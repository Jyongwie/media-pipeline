# 🎬 Distributed Media Rendering Pipeline

A full-stack, event-driven architecture designed to manage and simulate high-intensity background rendering jobs. Built with a globally distributed frontend, a strictly-typed concurrent backend, and real-time WebSocket telemetry.

![Architecture: Event-Driven](https://img.shields.io/badge/Architecture-Event--Driven-blue)
![Frontend: Angular 21](https://img.shields.io/badge/Frontend-Angular_21-dd0031?logo=angular)
![Backend: Go](https://img.shields.io/badge/Backend-Go_1.26-00add8?logo=go)
![Database: PostgreSQL](https://img.shields.io/badge/Database-Neon_Serverless-336791?logo=postgresql)
![Deployment: Edge/Containerized](https://img.shields.io/badge/Deployment-Vercel_%7C_Render-black)

## 🏗 System Architecture

This system decouples the user interface from heavy processing tasks using a background worker pool and real-time event streaming.

* **The Edge UI (Frontend):** Built with modern **Angular 21** (Standalone Components) and deployed to **Vercel's** global edge network for sub-second load times.
* **The Engine (Backend):** A strictly-typed **Golang** REST API running securely in a **Render** Docker container. It handles HTTP routing, custom CORS middleware, and database connection pooling.
* **The State (Database):** A **Neon** Serverless PostgreSQL cluster that manages job states (`queued` -> `processing` -> `completed`) and scales to zero when inactive.
* **The Worker Pool:** A concurrent Go routine (`goroutine`) that safely locks database rows using `FOR UPDATE SKIP LOCKED` to simulate intense GPU workloads without blocking the main web server.
* **The Telemetry (Real-Time):** A persistent **Gorilla WebSocket** tunnel that broadcasts state changes directly from the Go worker to the Angular UI, entirely eliminating the need for client-side polling.

## ✨ Key Features

* **Zero-Downtime Telemetry:** The UI updates instantly via WebSockets when a background worker finishes a job.
* **Concurrency Safe:** Implements strict row-level database locking to ensure multiple background workers never process the same job twice.
* **Environment Agnostic:** Smart API routing (`isDevMode`) and robust Environment Variable management for seamless local development and cloud production deployments.
* **Type-Safe Contract:** Unified data structures across the stack, ensuring the JSON payload from the Angular UI perfectly maps to the Go Structs and PostgreSQL schema.

---

## 🚀 Live Demo

**[View the Live Application Here](https://media-frontend-liart.vercel.app)** *(Note: The backend runs on a free Render tier and may take 30-50 seconds to spin up on the very first request).*

---

## 💻 Local Development Setup

To run this architecture on your local machine, follow these steps:

### Prerequisites
* [Node.js](https://nodejs.org/) (v24+)
* [Angular CLI](https://angular.dev/tools/cli) (`npm install -g @angular/cli`)
* [Go](https://go.dev/) (v1.26+)
* [Docker](https://www.docker.com/) (For local PostgreSQL)

### 1. Start the Local Database
Navigate to the root directory and spin up the Dockerized PostgreSQL instance:
```bash
docker-compose up -d