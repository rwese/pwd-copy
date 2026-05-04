# pwd-copy

Copy directory paths to clipboard. Go CLI, supports macOS (pbcopy), Linux (wl-copy, xclip, xsel), and custom clipboard commands.

## Commands

```bash
go run main.go           # Copy current directory path
go run main.go .         # Copy target directory path
go run main.go -r ../foo # Copy relative path to target

PWD_COPY_CLIPBOARD_COMMAND="xclip -selection clipboard" go run main.go
```

## Structure

- `main.go` - CLI entry point, clipboard detection, path resolution

## Code Style

- Standard Go formatting (`gofmt`)
- Error messages to stderr with `ERROR:` prefix
- Exit code 1 on errors, 0 on success

## Git Workflow

- Branches: `feat/<name>`, `fix/<name>` (see skill:git-conventional-commits)
- Commits: Conventional Commits (see skill:git-conventional-commits)

## Boundaries

**ALWAYS**
- Run `go build` to verify compilation before committing

**NEVER**
- Commit secrets or API keys
