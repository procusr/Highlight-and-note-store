package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Test database connection
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/mh_go?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test if we can ping the database
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("✅ Database connection successful")

	// Test if auth schema exists
	var schemaExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = 'auth')").Scan(&schemaExists)
	if err != nil {
		log.Fatal("Failed to check auth schema:", err)
	}

	if schemaExists {
		fmt.Println("✅ Auth schema exists")
	} else {
		fmt.Println("❌ Auth schema does not exist")
	}

	// Test if users table exists
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'auth' AND table_name = 'users')").Scan(&tableExists)
	if err != nil {
		log.Fatal("Failed to check users table:", err)
	}

	if tableExists {
		fmt.Println("✅ Users table exists")
	} else {
		fmt.Println("❌ Users table does not exist")

		// Try to create the schema and table
		fmt.Println("Creating schema and table...")

		_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		if err != nil {
			log.Fatal("Failed to create uuid extension:", err)
		}

		_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS auth")
		if err != nil {
			log.Fatal("Failed to create auth schema:", err)
		}

		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS auth.users (
				id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
				email text UNIQUE NOT NULL,
				encrypted_password text NOT NULL,
				email_confirmed_at timestamptz,
				created_at timestamptz NOT NULL DEFAULT timezone('utc'::text, now()),
				updated_at timestamptz NOT NULL DEFAULT timezone('utc'::text, now())
			)`)
		if err != nil {
			log.Fatal("Failed to create users table:", err)
		}

		fmt.Println("✅ Schema and table created successfully")
	}

	fmt.Println("🎉 Database setup is ready!")
}
