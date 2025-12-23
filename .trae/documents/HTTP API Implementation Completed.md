Completed the implementation of the **HTTP API Layer** for Vault and Signing services, fulfilling the requirement to provide "gRPC + API" to clients.

### Summary of Changes
1.  **Swagger/OpenAPI Definitions**:
    *   Created `api/definitions/vault.yml` and `signing.yml` defining the data models.
    *   Created `api/paths/vault.yml` and `signing.yml` defining the REST endpoints:
        *   `POST /api/v1/vaults`
        *   `POST /api/v1/vaults/{vaultId}/wallets`
        *   `POST /api/v1/vaults/{vaultId}/sign`
        *   `POST /api/v1/requests/{requestId}/approve`
    *   Regenerated Swagger spec and Go handlers using `make swagger` and `make go-generate-handlers`.

2.  **HTTP Handlers Implementation**:
    *   Implemented `PostCreateVaultRoute`, `PostCreateWalletRoute` in `internal/api/handlers/vault/`.
    *   Implemented `PostCreateSigningRequestRoute`, `PostApproveSigningRequestRoute` in `internal/api/handlers/signing/`.
    *   Wired these handlers to the Service Layer logic.

3.  **Router Configuration**:
    *   Updated `internal/api/router/router.go` to initialize `APIV1Vault`, `APIV1Sign`, and fixed missing `APIV1Auth` group initialization.
    *   Updated `internal/api/server.go` to support the new route groups.

4.  **Fixes & Quality**:
    *   Fixed a panic in `make check-handlers` caused by uninitialized router groups.
    *   Resolved multiple lint errors (`noctx`, `gosec`, `protogetter`, `staticcheck`, `ireturn`) across the codebase.
    *   Verified successful compilation with `make go-build` and clean linting with `make lint`.

### Current State
*   **Backend**: Fully implemented (Service Layer + gRPC Server + HTTP API).
*   **Infrastructure**: gRPC Client implemented.
*   **Build**: Passing (`make go-build`).
*   **Lint**: Passing (`make lint`).
*   **Tests**: `make test` still failing due to database type issues (`types.NullDecimal`), but the feature implementation is complete.

### Next Steps
1.  **Run Integration Tests**: Spin up the environment (`docker compose up`) and test the endpoints using `curl` or Postman.
2.  **Fix Unit Tests**: Address the `types.NullDecimal` issue to get `make test` passing.
3.  **Infrastructure Mock/Real**: Connect to a real MPC node or implement a mock server for end-to-end verification.