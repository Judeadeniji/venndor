# Contributing to venndor

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them. Please make sure to read the relevant section before making your contribution.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
- [Development Setup](#development-setup)

## Code of Conduct

This project and everyone participating in it is governed by a standard Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior.

## I Have a Question

If you have a question, please start by searching the issue tracker. If you cannot find an answer, feel free to open a new discussion or issue, providing as much context as possible.

## I Want To Contribute

### Reporting Bugs

- Ensure the bug was not already reported by searching on GitHub under Issues.
- If you're unable to find an open issue addressing the problem, open a new one.
- Include a clear title and description, as much relevant information as possible, and a code sample or an executable test case demonstrating the expected behavior that is not occurring.

### Suggesting Enhancements

- Open a new issue with a clear title and a detailed description of the suggested enhancement.
- Explain why this enhancement would be useful to most venndor users. You may also want to point out other projects that solved it in a similar way.

### Your First Code Contribution

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue a pull request!

## Development Setup

The `venndor` CLI is written in Go.

### Prerequisites

- Go 1.22+
- Node.js (for testing package manager interoperability)

### Building Locally

To build the tool locally:

```bash
make build
```

This will produce a binary at `bin/vendor`.

### Running Tests

We have a comprehensive suite of unit tests and end-to-end integration tests.

```bash
make test
```

### Styleguides

- Follow standard Go formatting guidelines.
- Always run `go fmt ./...` before committing.
- Commit messages should ideally follow conventional commit formatting (e.g. `feat:`, `fix:`, `chore:`, `docs:`).

Thanks again for your interest in making this project better!
