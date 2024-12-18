package db

import "github.com/jackc/pgx/v5/pgxpool"

func SetupSupavaultSchemaIfNotExists(pool *pgxpool.Pool) error {
	statements := []string{
		"CREATE SCHEMA IF NOT EXISTS supavault",
		"CREATE TABLE IF NOT EXISTS supavault.version_history ()",
		"ALTER TABLE supavault.version_history ADD COLUMN IF NOT EXISTS version text NOT NULL PRIMARY KEY",
		"ALTER TABLE supavault.version_history ADD COLUMN IF NOT EXISTS at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP",
		"ALTER TABLE supavault.version_history DROP CONSTRAINT IF EXISTS valid_semver",
		"ALTER TABLE supavault.version_history ADD CONSTRAINT valid_semver CHECK (version ~ '^(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(-(0|[1-9A-Za-z-][0-9A-Za-z-]*)(\\.[0-9A-Za-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$')",
	}

	return Batch(pool, statements)
}

func SetupSupabaseMigrationsSchemaIfNotExists(pool *pgxpool.Pool) error {
	statements := []string{
		"CREATE SCHEMA IF NOT EXISTS supabase_migrations",
		"CREATE TABLE IF NOT EXISTS supabase_migrations.schema_migrations ()",
		"ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS version text NOT NULL PRIMARY KEY",
		"ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS statements text[]",
		"ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS name text",
		"CREATE TABLE IF NOT EXISTS supabase_migrations.seed_files ()",
		"ALTER TABLE supabase_migrations.seed_files ADD COLUMN IF NOT EXISTS path text NOT NULL PRIMARY KEY",
		"ALTER TABLE supabase_migrations.seed_files ADD COLUMN IF NOT EXISTS hash text NOT NULL",
	}

	return Batch(pool, statements)
}
