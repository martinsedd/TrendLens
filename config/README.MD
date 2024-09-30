## 2. `config/`

#### `config.go`

- **Purpose**:
  This file is responsible for managing global configuration settings and environment variables for the project. It centralizes access to key configuration data used throughout the application.

- **Current Role**:
    - Loads configuration values such as API keys, ports, and other service-related settings, potentially from a `.env` file or environment variables.
    - Acts as a centralized point of configuration, ensuring the entire application has consistent access to necessary settings.

- **Best Practice**:
    - Safely load configuration from environment variables, ensuring any sensitive data is stored securely (e.g., through a `.env` file or secret manager).
    - Implement helper functions like `GetEnv` to provide controlled access to environment variables.
    - Ensure the configuration file only exposes necessary values to other parts of the app, keeping sensitive data hidden.
    - Consider validating configuration settings at the startup to catch misconfigurations early.

#### `mongo.go`

- **Purpose**:
  This file handles the setup and management of the MongoDB database connection.

- **Current Role**:
    - Initializes the MongoDB client and handles connections to the database.
    - Sets up any necessary indexes for MongoDB collections.
    - Ensures that the application successfully connects to MongoDB and manages connection retries and error handling.

- **Best Practice**:
    - Focus solely on database connection logic in this file. Include retry mechanisms and robust error handling for scenarios where the database might not be immediately available.
    - Avoid mixing unrelated configuration or business logic into this file. All MongoDB-related settings (e.g., connection strings) should come from `config.go`, keeping responsibilities clearly separated.
    - Consider logging connection status and errors to track potential issues during runtime.
