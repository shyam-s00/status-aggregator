# Status Aggregator

[![Build Status](https://github.com/shyam-s00/status-aggregator/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/shyam-s00/status-aggregator/actions/workflows/go.yml)

A service for aggregating status checks and health metrics from multiple sources into a unified dashboard.

## Description

This project collects status information from configured services and provides an aggregated view of system health. It utilizes a hybrid approach to ensure data accuracy:

*   **Current Status**: Scrapes the official HTML status pages of providers to determine if a system is currently operational or has an active incident.
*   **Incident History**: Consumes RSS feeds to build a historical timeline of past incidents.

## Prerequisites

*   Go (1.23 or higher)

## Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/shyam-s00/status-aggregator.git
    cd status-aggregator
    ```

2.  Download dependencies:
    ```bash
    go mod download
    ```

## Configuration

The application is driven by configuration files (e.g., `config.json`) where you can define:
*   System providers (RSS feed URLs, HTML status page URLs).
*   HTML selectors for scraping status text.
*   History limits for incident logs.

*Note: Ensure your configuration file is present in the working directory or the specific config path.*

## Usage

### Running Locally

To run the application directly from source:

```bash
go run ./cmd/server
