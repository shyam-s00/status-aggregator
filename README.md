# Status Aggregator
This is WIP project, more changes are incoming.

A service for aggregating status checks and health metrics from multiple sources into a unified dashboard or response.

## Description

This project collects status information from configured services and provides an aggregated view of system health. It is designed to help operations teams and developers quickly assess the overall state of their infrastructure.

## Prerequisites

- Go (1.21 or higher)

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd status-aggregator
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

## Configuration

Set up your environment variables by creating a `.env` file based on `.env.example` (if provided) or configuring the `config/` directory.

## Usage

To run the application directly:

```bash
go run main.go
```

To build and run the binary:

```bash
go build -o status-aggregator
./status-aggregator
```

## Testing

Run the test suite with:

```bash
go test ./...
```
