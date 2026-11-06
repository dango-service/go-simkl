# Commit Conventions  -  go-simkl

All commits in `go-simkl` must follow the [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) specification.

---

## 1. Format

```text
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

- **Type**: Must be one of the permitted types below.
- **Scope**: Lowercase noun describing the module or domain modified (e.g., `client`, `search`, `anime`, `sync`, `scrobble`, `auth`).
- **Description**: Concise, imperative, present tense, lowercase, no trailing period, <= 72 characters.
- **Body**: Explains *why* the change was made, not just *what*.

---

## 2. Permitted Types

| Type | Purpose | SemVer Impact |
|------|---------|---------------|
| `feat` | New functionality or endpoint exposed | MINOR |
| `fix` | Bug fix or error resolution | PATCH |
| `docs` | Documentation changes only | PATCH |
| `refactor` | Code change that neither fixes a bug nor adds a feature | PATCH |
| `perf` | Performance improvement or allocation reduction | PATCH |
| `test` | Adding or updating unit tests | None |
| `build` | Build system, packaging, or `go.mod` updates | None |
| `ci` | CI workflow changes (`.github/workflows`) | None |
| `chore` | Maintenance tasks | None |

---

## 3. Breaking Changes

Indicated by an exclamation mark `!` before the colon:

```text
feat(client)!: change NewClient signature to require clientID
```

---

## 4. Examples

- `feat(client): implement core client, transport, and rate limit parsing`
- `feat(search,anime): implement text search, id resolution, and anime details`
- `feat(sync): implement watchlist retrieval, activity checks, and history sync`
- `feat(scrobble,user,auth): implement scrobbling sessions, user profile, and oauth2 helpers`
- `ci(github): add matrix test workflow, developer guidelines, and complete phase 5`
