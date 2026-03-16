# Contributing

## Getting Started

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes and test locally: `go run ./cmd/agent`
4. Commit with a clear message
5. Open a pull request against `main`

## Code Style

- Follow standard Go formatting: `gofmt -w .`
- Keep functions small and focused
- Add comments for exported functions

## Testing

Run tests before submitting:
```
go test ./...
```

## Reporting Issues

Open a GitHub issue with:
- What you expected to happen
- What actually happened
- Steps to reproduce
