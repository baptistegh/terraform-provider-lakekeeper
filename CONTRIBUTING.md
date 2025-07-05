# How to Contribute

We want to make contributing to this project as easy and smooth as possible.

## Reporting Issues

If you encounter a bug or have a feature request, please report it on the
[issue tracker](https://github.com/baptistegh/terraform-provider-lakekeeper/issues).

If you already plan to fix the issue yourself, it's not necessary to open a separate issue beforehand — you can directly open a merge request (MR) with a description of the problem you're addressing.

## Contributing Code

Merge requests are always welcome. If you're unsure whether your contribution fits well within the project scope, feel free to open an issue first to discuss your idea.

If you'd like to work on an issue that is already open, feel free to leave a comment saying you're interested.
We'll assign the issue to you to avoid duplicate work.

### Preparing your Merge Request

Before submitting your merge request, please run:

```sh
make reviewable
```

This command will lint and format your code, run tests, and generate any necessary documentation files to ensure your changes meet the project standards.
Running it locally helps reduce CI failures and keeps the codebase consistent.

### Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to keep a consistent and meaningful commit history.
Please write your commit messages accordingly, for example:

- `feat: add support for new resource`
- `fix: resolve issue with provider initialization`
- `docs: update README with new examples`

This helps with automated changelog generation and improves collaboration.

## Coding Style

We follow standard Go best practices and use [`golangci-lint`](https://golangci-lint.run/) to enforce code quality and formatting.

Before submitting a MR, please review the codebase and ensure your changes are consistent with the existing style.

## Setting Up Your Local Development Environment

To install dependencies and ensure the project builds:

```sh
make build
```

### Running Tests

This project includes both unit tests and acceptance tests.

#### Unit Tests

```sh
make test
```

#### Acceptance Tests

1. Start the test environment:

   ```sh
   make testacc-up
   ```

2. Run the acceptance tests:

   ```sh
   make testacc
   ```

3. Tear everything down (volumes will not be preserved):

   ```sh
   make testacc-down
   ```

## Licensing

By contributing code to this project, you agree to license your contribution under the terms of the MPL-2.0 license.
