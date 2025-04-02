# Testing Guide for Task Management System

This document provides instructions for running the unit tests, generating test coverage reports, and interpreting the results for the Task Management System.

---

## Prerequisites

1. Ensure you have Go installed (version 1.20 or later).
2. Install the required dependencies by running:
   ```bash
   go mod tidy
   ```

---

## Running Unit Tests

To run all unit tests in the project, execute the following command in the root directory of the project:

```bash
go test -v ./...
```

- The `-v` flag enables verbose output, showing detailed information about each test.
- The `./...` ensures that tests in all subdirectories are executed.

---

## Running Specific Test Suites

To run tests for a specific package, navigate to the package directory or specify the package path. For example:

- To run tests for the `Usecases` package:
  ```bash
  go test -v ./Usecases
  ```

- To run tests for the `Infrastructure` package:
  ```bash
  go test -v ./Infrastructure
  ```

---

## Generating Test Coverage Report

To generate a test coverage report, use the following command:

```bash
go test -coverprofile=coverage.out ./...
```

This will create a file named `coverage.out` in the root directory.

### Viewing the Coverage Report

1. To view the coverage summary in the terminal:
   ```bash
   go tool cover -func=coverage.out
   ```

2. To view the coverage report in a browser:
   ```bash
   go tool cover -html=coverage.out
   ```

   This will open an interactive HTML report showing which lines of code are covered by tests.

---

## Interpreting Test Coverage Metrics

- **Coverage Percentage**: Indicates the percentage of code covered by tests. Aim for a high percentage (e.g., 80% or above) to ensure good test coverage.
- **Uncovered Lines**: The HTML report highlights uncovered lines in red. These lines should be reviewed to determine if additional tests are needed.

---

## Continuous Integration (CI) Testing

The project includes a GitHub Actions CI pipeline that automatically runs tests and uploads the coverage report. The CI configuration is located in `.github/workflows/ci.yml`.

### Key Steps in CI Pipeline:
1. **Run Tests**: Executes all unit tests using `go test`.
2. **Generate Coverage Report**: Creates a `coverage.out` file.
3. **Upload Coverage Report**: Saves the coverage report as an artifact for review.

---

## Troubleshooting

- **Test Failures**: If a test fails, review the error message and the corresponding test case to identify the issue.
- **Dependency Issues**: Run `go mod tidy` to ensure all dependencies are installed.
- **Coverage Report Issues**: Ensure the `coverage.out` file is generated before running `go tool cover`.

---

## Example Commands

- Run all tests:
  ```bash
  go test -v ./...
  ```

- Generate and view coverage report:
  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```

- Run tests for a specific file:
  ```bash
  go test -v ./Usecases/task_usecases_test.go
  ```

---

By following this guide, you can ensure the quality and reliability of the Task Management System through comprehensive testing.
