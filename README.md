# Expert System

## Dependencies

* go version 1.13
* github.com/chzyer/readline

## Build

```
git clone https://github.com/acarlson99/expert-system.git
cd expert-system
go build
```

Or

```
go get go get github.com/acarlson99/expert-system
go run github.com/acarlson99/expert-system
```

## Run

```
./expert-system						# Enter interactive mode
./expert-system test/test-all.xs	# Evaluate file
./expert-system -f test/test-all.xs	# Enter interactive mode after file evaluated
```

## Syntax

### Operators

| Operator    | Example   |
| -           | -         |
| Paren       | `(A + B)` |
| Not         | `!A`      |
| And         | `A + B`   |
| Or          | `A \ B`   |
| Xor         | `A ^ B`   |
| Implication | `A => B`  |

### Commands

| Command | Description              |
| :-:     | -                        |
| `exit`  | Exit program             |
| `list`  | List variables and rules |
| `help`  | Display help             |
