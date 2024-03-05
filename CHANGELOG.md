# CHANGELOG

- 2024-03-05 **`v1.0.0`**
  - Perform migrations
    - PostgreSQL only currently
    - Roll forward or backward
    - Remove migration table on `reset`
  - Update `migratable_state`
    - Within a transaction with the migration
    - Columns changed for clarity

- 2024-03-04
  - Connect to the database
    - PostgreSQL only currently
  - Create the `migratable_state` table
  - Return the current migration version

- 2024-03-03
  - Initial project creation
  - Standard repo files added
  - Command argument parsing
  - Loading migrations
