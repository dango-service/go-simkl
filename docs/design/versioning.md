# SemVer and Release Lifecycle  -  `go-simkl`

This document outlines the versioning strategy and release policies for `go-simkl`.

---

## 1. Semantic Versioning (SemVer 2.0.0)

Version numbers follow `vMAJOR.MINOR.PATCH`:

- **MAJOR**: Breaking changes to the public Go API or method signatures.
- **MINOR**: Backward-compatible new endpoints, options, or data fields.
- **PATCH**: Backward-compatible bug fixes, performance improvements, or documentation updates.

---

## 2. Pre-v1.0.0 Lifecycle

While in `v0.x.y`:
- Minor releases may introduce refined API signatures as real usage in `Dango` informs improvements.
- Stable release `v1.0.0` will be cut when all core tracking scenarios are verified in production.
