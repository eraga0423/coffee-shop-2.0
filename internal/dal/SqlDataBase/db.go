package SqlDataBase

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	Db *sql.DB
}
type Database interface {
	Init() error
}

func (d *DB) Init() error {
	var err error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)
	fmt.Println("Waiting for database to be ready...")
	for i := 0; i < 5; i++ {
		d.Db, err = sql.Open("postgres", info)
		if err == nil {
			if err = d.Db.Ping(); err == nil {
				log.Println("Successfully connected to the database!")
				err := d.InsertInto()
				if err != nil {
					return err
				}
				return nil
			}
		}

		log.Printf("Database not ready yet... retrying (%d/5)\n", i+1)
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("failed to connect to database after multiple attempts: %v", err)
}

func NewDB() DB {
	return DB{
		Db: new(sql.DB),
	}
}
