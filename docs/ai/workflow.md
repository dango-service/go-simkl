# AI Agent Workflow  -  go-simkl

Methodology and rigorous workflow that all AI agents must follow when modifying this repository.

---

## 1. Core Principle: "Read Before Writing"

1. **Consult the Sources of Truth**:
   - `AGENTS.md` for project rules and invariants.
   - `docs/reference/simkl-api.md` to validate Simkl official schema types and fields.
   - `docs/design/public-api.md` to verify public API contracts.
2. **Do Not Invent APIs or Fields**:
   - Every field in a struct must match official Simkl OpenAPI 3.1 specifications.
   - Public signatures must follow established design decisions.

---

## 2. Managing Dependencies and Decisions

- **No Arbitrary Dependencies**: Any proposed addition to `go.mod` must be evaluated and formally justified in `docs/design/`.
- **Clarifying Ambiguities**: When facing ambiguities, consult official API specs and repository invariants before making assumptions.

---

## 3. Feature Implementation Cycle

1. **Define the HTTP Endpoint and Payload**: Document exact path, headers, query parameters, and JSON models.
2. **Define Go Types**: Create compact, strongly typed Go structs.
3. **Write Tests with Fixtures**:
   - Add mock JSON responses in `testdata/`.
   - Write unit tests using `httptest.Server` verifying sent variables, headers, and output.
4. **Implement Service Methods**: Connect transport layer to the public service.
5. **Verify**: Run `go test -v ./...` and `go vet ./...`.
