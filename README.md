# MIGRATABLE

Command-line database migrations using up/down SQL files.
Simple, with stand-alone executables small enough to commit
to your own repos alongside your code.
The initial release is for _PostgreSQL_ only.

- [AGPL license](./LICENSE.txt)
- [CHANGELOG](./CHANGELOG.md)

```
/migrations
  /001_create_accounts_table
    down.sql
    up.sql
  /002_add_Sample_data
    down.sql
    up.sql
```

The migration version (the leading digits in the folder name)
can be separated from the migration name that follows by any
non-digit character. Leading zeroes are optional.

## Usage

- Place your connection string into an environment variable
- Run the appropriate version from the [`builds``](./builds) folder
- Specify the environment variable name when calling Migratable

## Command arguments

- `<folder>` - location of your migration scripts
- `<conn-env>` - environment variable holding connection string
- `<action>` - migration action to perform
- `[version]` - migration number (if required)

## Supported actions

- `info` - show migration status
- `list` - list known migrations
- `reset` - remove all migrations
- `latest` - apply new migrations
- `next` - roll forward one migration
- `back` - roll backward one migration
- `target` - target specific version
