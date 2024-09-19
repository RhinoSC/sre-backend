package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Conectar a la base de datos SQLite
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Leer el archivo de migración
	migration, err := os.ReadFile("./migrations/db_tables.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Ejecutar las migraciones
	_, err = db.Exec(string(migration))
	if err != nil {
		log.Fatal(err)
	}

	var test_data = func() {
		// Leer el archivo de migración
		migration, err := os.ReadFile("./migrations/seed_srex.sql")
		if err != nil {
			log.Fatal(err)
		}

		// Ejecutar las migraciones
		_, err = db.Exec(string(migration))
		if err != nil {
			log.Fatal(err)
		}
	}

	test_data()

	fmt.Println("Migraciones ejecutadas con éxito")
}
