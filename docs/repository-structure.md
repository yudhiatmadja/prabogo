# Repository Structure

![Alt text](./images/internal.jpg "Internal")

## Overview

The Go Boilerplate project follows hexagonal architecture (also known as ports and adapters) principles with a carefully designed internal repository structure. This document explains how the project's internal repository is organized and how the components interact with each other.

## Repository Structure

### Directory Organization

```
internal/
├── adapter/              # Implementations of ports (adapters)
│   ├── inbound/          # Adapters receiving requests into the application
│   │   ├── command/      # CLI command adapters
│   │   ├── fiber/        # HTTP Fiber framework adapters
│   │   ├── mcp/          # Model Context Protocol adapters
│   │   ├── rabbitmq/     # RabbitMQ consumer adapters
│   │   └── temporal/     # Temporal workflow adapters
│   └── outbound/         # Adapters sending requests to external systems
│       ├── http/         # HTTP client adapters
│       ├── postgres/     # PostgreSQL database adapters
│       ├── rabbitmq/     # RabbitMQ producer adapters
│       └── redis/        # Redis cache adapters
├── domain/               # Core business logic
│   └── client/           # Client domain logic example
├── migration/            # Database migration scripts
│   └── postgres/         # PostgreSQL specific migrations
├── model/                # Data structures/entities
└── port/                 # Interface definitions
    ├── inbound/          # Inbound ports (for receiving requests)
    └── outbound/         # Outbound ports (for sending requests)
```

## Key Components

### Ports (Interfaces)

Located in the `internal/port/` directory, ports define the contracts for how the application communicates with the external world:

1. **Inbound Ports (`internal/port/inbound/`)**: Define how external systems communicate with the application.
   - HTTP interfaces (`registry_http.go`)
   - Message consumer interfaces (`registry_message.go`)
   - Command interfaces (`registry_command.go`)
   - Workflow interfaces (`registry_workflow.go`)
   - MCP interfaces (`registry_mcp.go`)

2. **Outbound Ports (`internal/port/outbound/`)**: Define how the application communicates with external systems.
   - Database interfaces (`registry_database.go`)
   - HTTP client interfaces (`registry_http.go`)
   - Message producer interfaces (`registry_message.go`)
   - Cache interfaces (`registry_cache.go`)

### Adapters (Implementations)

Located in the `internal/adapter/` directory, adapters implement the ports:

1. **Inbound Adapters (`internal/adapter/inbound/`)**: 
   - `fiber/`: HTTP handlers using the Fiber framework
   - `rabbitmq/`: Message consumers using RabbitMQ
   - `command/`: CLI command handlers
   - `temporal/`: Workflow adapters using Temporal
   - `mcp/`: MCP tool adapters for AI clients

2. **Outbound Adapters (`internal/adapter/outbound/`)**: 
   - `postgres/`: PostgreSQL database adapters
   - `http/`: HTTP client adapters for external APIs
   - `rabbitmq/`: Message producers using RabbitMQ
   - `redis/`: Cache adapters using Redis

### Domain Logic

Located in `internal/domain/`, this contains the core business logic of the application, independent from external systems. The domain layer uses ports to interact with external systems but isn't aware of the actual implementations (adapters).

### Models

Located in `internal/model/`, these are the data structures and entities used throughout the application. They represent the business objects and their attributes.

## How Components Interact

### Dependency Flow

The dependency flow follows the hexagonal architecture principle:

1. **Domain logic** depends on **ports** (interfaces)
2. **Adapters** implement these **ports**
3. The application wires everything together at startup

This ensures that the domain logic remains independent of the implementation details of external systems.

### Registry Pattern

The project uses a registry pattern to organize and access adapters:

1. Registry interfaces in port directory (`registry_*.go` files) define groups of related ports
2. Registry implementations in adapter directories handle the creation and management of adapter instances
3. The application uses these registries to access the appropriate adapters at runtime

## Component Creation Process

The project uses `prabogo-cli` to automate the creation of new components:

### Core Usage Modes

- `prabogo-cli`: Starts interactive mode. If `prabogo.lock.json` exists, it shows `run` and `make`; otherwise it shows `install`.
- `prabogo-cli install <project-name>`: Creates a new Prabogo project and generates `prabogo.lock.json`.
- `prabogo-cli make`: Opens interactive code generation (model, migration, inbound, outbound).
- `prabogo-cli run`: Opens interactive run/test target selection.

### Inbound Adapters

- `prabogo-cli inbound-http-fiber <name>`: Creates HTTP handler interfaces, adapters, and registry updates
- `prabogo-cli inbound-message-rabbitmq <name>`: Creates RabbitMQ consumer interfaces, adapters, and registry updates
- `prabogo-cli inbound-command <name>`: Creates command handler interfaces, adapters, and registry updates
- `prabogo-cli inbound-workflow-temporal <name>`: Creates Temporal workflow interfaces, adapters, and registry updates

### Outbound Adapters

- `prabogo-cli outbound-database-postgres <name>`: Creates database interfaces, adapters, and registry updates
- `prabogo-cli outbound-http <name>`: Creates HTTP client interfaces, adapters, and registry updates
- `prabogo-cli outbound-message-rabbitmq <name>`: Creates RabbitMQ producer interfaces, adapters, and registry updates
- `prabogo-cli outbound-cache-redis <name>`: Creates Redis cache interfaces, adapters, and registry updates
- `prabogo-cli outbound-workflow-temporal <name>`: Creates Temporal workflow starter interfaces, adapters, and registry updates

### Models and Migrations

- `prabogo-cli model <name>`: Creates model structure with basic fields
- `prabogo-cli migration-postgres <name>`: Creates PostgreSQL migration files
- `<name>` supports `snake_case` (for example: `user_profile`, `client`).

### Run and Validation Targets

- `prabogo-cli build`: Builds the Docker image
- `prabogo-cli http`: Runs HTTP mode in Docker
- `prabogo-cli message <subscriber>`: Runs message consumer mode in Docker
- `prabogo-cli command <cmd> <val>`: Runs command mode in Docker
- `prabogo-cli workflow <name>`: Runs workflow worker mode in Docker
- `prabogo-cli generate-mocks`: Generates mocks from `go:generate` declarations
- `prabogo-cli lint`: Runs lint checks
- `prabogo-cli test`: Runs unit tests
- `prabogo-cli test-coverage`: Runs tests and writes coverage output
- `prabogo-cli test-integration`: Runs integration tests

## Testing Strategy

The repository structure is designed to facilitate testing:

1. **Mocks**: Generated automatically for each port interface using `mockgen` (`prabogo-cli generate-mocks`)
2. **Unit Tests**: Test domain logic in isolation using mocks of the ports (`prabogo-cli test`)
3. **Coverage**: Validate test coverage for critical paths (`prabogo-cli test-coverage`)
4. **Integration Tests**: Validate infrastructure integration (`prabogo-cli test-integration`)

## Benefits of This Structure

1. **Clear Separation of Concerns**: Domain logic is separate from infrastructure details
2. **Flexibility**: Easy to swap out adapters (e.g., switch from Postgres to MySQL)
3. **Testability**: Domain logic can be tested in isolation
4. **Maintainability**: Changes to one component don't affect others
5. **Consistency**: Standardized approach to adding new features
6. **Code Generation**: Repetitive boilerplate code can be generated using `prabogo-cli make` or direct generator commands

## Best Practices

1. **Keep Domain Logic Pure**: Domain logic should not depend on specific technologies or frameworks
2. **Use Interfaces Wisely**: Define clear interfaces between components
3. **Follow the Dependency Rule**: Dependencies should point inward toward the domain
4. **Generate Mocks**: Use `prabogo-cli generate-mocks` after creating new ports
5. **Automate Creation**: Use the provided `prabogo-cli` generators to create new components consistently
