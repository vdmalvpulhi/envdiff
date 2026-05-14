// Package renamer provides key-renaming functionality for envdiff diff results.
//
// It is useful when variable names are being migrated between environments
// (e.g. renaming DB_HOST to DATABASE_HOST) and you want the diff output to
// reflect the new canonical names without manually editing every .env file.
//
// Basic usage:
//
//	opts := renamer.DefaultOptions()
//	opts.Renames = map[string]string{
//		"DB_HOST": "DATABASE_HOST",
//		"DB_PORT": "DATABASE_PORT",
//	}
//	renamed := renamer.Apply(results, opts)
//
// Pairs can also be built from "OLD=NEW" strings (e.g. from CLI flags):
//
//	m, err := renamer.BuildRenames([]string{"DB_HOST=DATABASE_HOST"})
package renamer
