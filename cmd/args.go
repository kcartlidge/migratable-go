// Argument parsing.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// config holds the user request.
type config struct {
	folder, connEnv, action string
	target                  int
}

// getConfig fetches the user request from the command arguments.
// It also displays the options.
func getConfig() (*config, error) {
	c := config{
		folder:  ".",
		connEnv: "MIGRATABLE",
		action:  "info",
		target:  0,
	}
	c.usage()

	// Fetch the arguments, with a basic sanity check.
	if len(os.Args) < 4 || len(os.Args) > 5 {
		return nil, errors.New("unexpected number of command arguments")
	}
	c.folder = strings.TrimSpace(os.Args[1])
	c.connEnv = strings.ToUpper(strings.TrimSpace(os.Args[2]))
	c.action = strings.ToLower(strings.TrimSpace(os.Args[3]))

	// Check the requested action.
	switch c.action {
	case "info":
		fallthrough
	case "reset":
		fallthrough
	case "latest":
		fallthrough
	case "next":
		fallthrough
	case "back":
		// Further sanity checks as should not have a target version.
		if len(os.Args) > 4 {
			return nil, errors.New("unexpected extra argument(s)")
		}
		break
	case "target":
		// Further sanity checks as need a target version.
		if len(os.Args) != 5 {
			return nil, errors.New("expected a target migration version")
		}
		t, err := strconv.Atoi(os.Args[4])
		if err != nil {
			return nil, errors.New("target version is not an integer")
		}
		if t < 0 {
			return nil, errors.New("target version cannot be negative")
		}
		c.target = t
		break
	default:
		return nil, errors.New("unknown action requested")
	}
	return &c, nil
}

// usage displays the command line instructions.
func (c *config) usage() {
	fmt.Println("USAGE:")
	fmt.Println("  <folder>   Location of your migration scripts")
	fmt.Println("  <conn-env> Environment variable holding connection string")
	fmt.Println("  <action>   Migration action to perform")
	fmt.Println("  [version]  Migration number (if required)")
	fmt.Println()
	fmt.Println("ACTIONS:")
	fmt.Println("  info       Show migration status")
	fmt.Println("  reset      Remove all migrations")
	fmt.Println("  latest     Apply new migrations")
	fmt.Println("  next       Roll forward one migration")
	fmt.Println("  back       Roll backward one migration")
	fmt.Println("  target     Target specific [version]")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  migratable my-migrations MY_CONNSTR info")
	fmt.Println("  migratable my-migrations MY_CONNSTR latest")
	fmt.Println("  migratable my-migrations MY_CONNSTR target 4")
	fmt.Println("  migratable my-migrations MY_CONNSTR reset")
}

// describe details what has been requested.
func (c *config) describe() {
	section("REQUESTED")
	fmt.Println("  Folder   : ", c.folder)
	fmt.Println("  Conn Env : ", c.connEnv)
	fmt.Println("  Action   : ", c.action)
	if c.action == "target" {
		fmt.Println("  Version:  ", c.target)
	}
	fmt.Println()
}
