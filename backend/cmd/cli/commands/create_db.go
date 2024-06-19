package commands

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/spf13/cobra"
)

// createDBCmd is a subcommand to create a database
var createDBCmd = &cobra.Command{
	Use:   "create-db",
	Short: "Runs `CREATE DATABASE <name>;`",
	Long:  "Runs `CREATE DATABASE <name>;`",
	Run:   createDB,
}

func createDB(cobraCmd *cobra.Command, args []string) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslMode := os.Getenv("POSTGRES_SSLMODE")
	user := os.Getenv("POSTGRES_USER")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		host,
		port,
		user,
		password,
		sslMode,
	)

	fmt.Println(connStr)

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Name of the database to create
	newDB := "new_database"

	// SQL statement to create a new database
	createDBSQL := fmt.Sprintf("CREATE DATABASE %s", newDB)

	// Execute the SQL statement to create the database
	_, err = db.Exec(createDBSQL)
	if err != nil {
		log.Fatal("Error creating database:", err)
	}

	fmt.Printf("Database %s created successfully.\n", newDB)
}

func init() {
	rootCmd.AddCommand(createDBCmd)
}
