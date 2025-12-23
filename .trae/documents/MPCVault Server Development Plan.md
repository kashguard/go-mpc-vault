Thanks for the clarification. The plan is now refined to strictly separate the protocols: **gRPC Only** for Backend-to-Infra communication, and **Dual Protocol (gRPC + HTTP)** for Client-to-Backend communication.

# MPCVault Server Development Plan (v3)

## Phase 1: Proto Setup & Generation
**Goal**: Import Infra protos (for backend client) and define App protos (for client interface).

1.  **Import Infra Protos** (`proto/infra/v1/`):
    *   Copy `*.proto` from `/Users/caimin/Desktop/kms/go-mpc-wallet/proto/infra/v1/`.
    *   *Purpose*: Used ONLY to generate the **gRPC Client** to talk to the Infra layer.
2.  **Define App Protos** (`proto/api/v1/`):
    *   Create `auth.proto`, `vault.proto`, `signing.proto`.
    *   Add `google.api.http` annotations to these protos to enable HTTP transcoding (optional, but good practice if we use gRPC-Gateway later; for now, we will maintain parallel Echo handlers as per project standard).
    *   *Purpose*: Used to generate the **gRPC Server** interface for mobile/web clients.
3.  **Generate Code**:
    *   Generate `internal/infra/grpc` (Client stubs).
    *   Generate `internal/api/grpc` (Server stubs).

## Phase 2: Database & Models
**Goal**: Set up the persistent storage.

1.  **Migrations**: Create SQL files for `organizations`, `vaults`, `wallets`, `credentials` (WebAuthn), `requests`, `approvals`, `risk_rules`.
2.  **Models**: Run `make sql` to generate SQLBoiler models.

## Phase 3: Infrastructure Client (gRPC Only)
**Goal**: Implement the client that talks to the Infra layer.

1.  **MPC Client** (`internal/infrastructure/mpc/`):
    *   Establish mTLS gRPC connection to Infra nodes.
    *   **Strictly gRPC**: No HTTP calls to Infra.
    *   Implement methods: `CreateKey` (DKG), `Sign` (Threshold).

## Phase 4: Service Layer (Business Logic)
**Goal**: The core logic that bridges the Database, Infra Client, and Client Interfaces.

1.  **Services**:
    *   `AuthService`: WebAuthn verification.
    *   `VaultService`: Calls `MpcClient.CreateKey`, saves to DB.
    *   `SigningService`: Validates policies, calls `MpcClient.Sign`.
    *   *Design*: These services return Go structs/errors, unaware of HTTP vs gRPC.

## Phase 5: Client Interfaces (Dual Protocol)
**Goal**: Expose the Service Layer to clients via both gRPC and HTTP.

1.  **gRPC Server** (`internal/api/grpc_server/`):
    *   Implement `api/v1` proto interfaces.
    *   Convert Proto messages <-> Service Models.
    *   Call Service Layer.
2.  **HTTP Server** (`internal/api/handlers/`):
    *   Use existing Echo framework.
    *   Bind JSON <-> Service Models.
    *   Call Service Layer.
    *   *Note*: This ensures existing tooling (Swagger/OpenAPI) remains valid while adding gRPC support.

## Phase 6: Entrypoint & Config
1.  **Server Startup**:
    *   Modify `cmd/server/main.go` to start **both** the Echo Server (e.g., :8080) and gRPC Server (e.g., :9090).
2.  **Config**: Add `GrpcPort` and Infra Connection settings.

---

**Immediate Next Step**:
I will proceed with **Phase 1**: Copying the Infra protos and creating the new App API protos.
