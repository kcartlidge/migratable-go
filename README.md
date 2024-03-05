# MIGRATABLE v1.0.0

**Command-line database migrations using up/down SQL files.
Simple, with stand-alone executables small enough to commit
to your own repos alongside your code.**

_The initial release is for PostgreSQL only._

- [AGPL license](./LICENSE.txt)
- [CHANGELOG](./CHANGELOG.md)

Copyright 2024 K Cartlidge.

## Contents

- [About database migrations](#about-database-migrations)
  - [Example migration scripts structure](#example-migration-scripts-structure)
- [Download and installation](#download-and-installation)
- [Usage](#usage)
  - [Example connection string](#example-connection-string)
  - [Command arguments](#command-arguments)
  - [Supported actions](#supported-actions)
  - [Example commands](#example-commands)
- [Generating new builds](#generating-new-builds)
  - [Creating new releases](#creating-new-releases)
  - [One-off builds during development](#one-off-builds-during-development)

## About database migrations

These are sequences of instructions to build up or tear down
things in the database, whether that be structural (eg tables)
or data.

Each migration contains an `up` and a `down`, which allows you
to faithfully recreate the database at a particular point in
time (mainly during development or deployment) whilst also
allowing you to safely apply new changes in a controlled manner.

You can roll forward to the latest, reset back to before you
started, move forward or backward, or target a specific
version.

Once you have reached the point where you are dealing with real
data you use it for structural changes like adding indexes or
extra columns, or for adding things to look-up lists etc.

- Pre-built binaries; not dependent on Go
- It fits any ecosystem (eg Go, Python, Node, .Net etc)
- And it's easy to configure, needing only:
  - A connection string in an environment variable
  - A folder with named migrations using up/down SQL scripts
- Your database structure can be version-controlled
- Roll your database backwards as well as forwards
- Seed/pre-populate, update, or remove data
- Run in transactions for atomic up/down

### Example migration scripts structure

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

There's a set of [example migrations](./cmd/sample-migrations)
in the project.

## Download and installation

_There is no installation needed._

Download the appropriate version from the list below and
run it directly from the command line/terminal.

- [Download for Linux](./builds/linux/)
- [Download for Windows](./builds/windows/)
- [Download for Mac (Apple Silicon)](./builds/macos/)
- [Download for Mac (Intel)](./builds/macos-x64/)

Migratable is standalone and small, so you can place a copy
of the download into any repo/codebase that is using it and
thereby guarantee it will always be available to your code
and any build tool-chain.

## Usage

- Place your connection string into an environment variable
- Specify the environment variable name when calling Migratable
- You can test your connection using the `info` action (see below)

Migratable maintains details of what migrations have been applied
within a `migratable_state` table.  It does _not_ look at your
database and work it out.  It follows then that for consistency
all changes should be done via Migratable and not by hand as that
will cause Migratable's opinion of the database state and the
underlying reality to fall out of sync.

### Example connection string

Mac and Linux

``` shell
export MIGRATABLE="host=127.0.0.1 port=5432 dbname=example search_path=example user=example password=example sslmode=disable"
```

Windows

``` shell
set MIGRATABLE=host=127.0.0.1 port=5432 dbname=example search_path=example user=example password=example sslmode=disable
```

- The login, database, and schema must already exist
- The `search_path` parameter is the database schema
- The Windows version does _not_ include double-quotes

### Command arguments

- `<folder>` - location of your migration scripts
- `<conn-env>` - environment variable holding connection string
- `<action>` - migration action to perform
- `[version]` - migration number (if required)

### Supported actions

- `info` - show migration status
- `reset` - remove all migrations (and tracking table)
- `latest` - apply new migrations
- `next` - roll forward one migration
- `back` - roll backward one migration
- `target` - target specific version

Note that there is a distinction between rolling back all
migrations using the `back` or `target 0` actions, as opposed
to using `reset`.

The former actions will undo all migrations as expected but
will leave the history in the `migratable_state` tracking
table. However using `reset` will drop that table when
the rollbacks complete, removing all trace of _Migratable_
and any migration history.

### Example commands

``` shell
cd cmd
migratable ./sample-migrations MIGRATABLE info      # shows current state
migratable ./sample-migrations MIGRATABLE target 3  # moves forward/back to v3
migratable ./sample-migrations MIGRATABLE latest    # ensures all applied
migratable ./sample-migrations MIGRATABLE reset     # removes all traces
```

In the above, `MIGRATABLE` is the name of an environment
variable containing the database connection string. The
variable should be UPPERCASE in your environment; the name
passed in will be treated as such when checking.

---

## Generating new builds

_This should only be necessary if you are making changes
to Migratable itself.  As a user of Migratable you can
use the existing builds mentioned earlier._

**Note:** if you are unable to run the scripts ensure they
are executable (on Linux or Mac) by using `chmod`.  For example
`chmod +x cmd/scripts/linux.sh` will do the Linux one.

- There are scripts for Linux, Windows, Mac (Intel), and Mac (ARM)
- Run the script that relates to _your_ current platform
- Each script generates the builds for _all_ the platforms
- You _must_ run the scripts from within the `cmd` folder
- Builds are small and should be checked in
  - Only generate new builds when new releases are needed

### Creating new releases

Run the relevant script for _your_ system.

``` shell
cd cmd
./scripts/linux.sh
./scripts/mac.sh
scripts\windows
```

If you make changes to Linux or Mac scripts within Windows
and need to ensure the executable permission is set on them
you can do it on the command line after committing them:

```
git update-index --chmod=+x scripts\linux.sh
git update-index --chmod=+x scripts\mac.sh
```

### One-off builds during development

When development is completed you use the above-mentioned
scripts to generate a full set of releases.

However _during_ your development you probably want a quicker
turnaround so from the following commands you can create a new
version as a one-off by (choose according to _your current_
system).

``` shell
cd cmd
go build -o ../builds/macos/migratable
go build -o ../builds/macos-x64/migratable
go build -o ../builds/linux/migratable
go build -o ..\builds\windows\migratable.exe
```
