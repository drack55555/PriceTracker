package main

import (
	"database/sql"
	"log"
	"time"
)

type StorageService struct {
	db *sql.DB
}

func NewStorageService(path string) (*StorageService, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("No Sql driver to Open")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	s := &StorageService{db: db}

	if err = s.initDb(); err != nil {
		return nil, err
	}

	log.Printf("Database connection successfull and Tables Initialized")
	return s, nil
}

func (s *StorageService) initDb() error {
	createSQLTable := `
		CREATE TABLE IF NOT EXISTS products	(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL UNIQUE,
			target_price REAL NOT NULL,
			last_price REAL NOT NULL,
			last_checked DATETIME,
			alert_sent BOOLEAN DEFAULT FALSE
		);
	`

	_, err := s.db.Exec(createSQLTable)
	if err != nil {
		log.Printf("Error creating database: %v", err)
		return err
	}
	return nil
}

func (s *StorageService) CreateNewTrackRequest(url string, targetPrice float64) error {

	// using insert ignore because url is unique and if already present url so it will ignore
	insertSQL := `INSERT OR IGNORE INTO products (url, target_price) VALUES (?, ?)`
	_, err := s.db.Exec(insertSQL, url, targetPrice)
	if err != nil {
		log.Printf("Failed to insert track request: %v", err)
		return err
	}

	// Step 5c: Return 'nil' (no error) on success.
	return nil
}

// this function retrieves every Product we are tracking
func (s *StorageService) GetAllProducts() ([]*Product, error) {
	query := `SELECT id, url, targetPrice, lastPrice, lastChecked FROM products`
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("Failed to query Product: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []*Product

	for rows.Next() {
		p := &Product{}

		//fields which can be bank needs 0 for safe keeping
		var lastChecked sql.NullTime

		err := rows.Scan(&p.ID, &p.URL, &p.TargetPrice, &p.LastPrice, &lastChecked)
		if err != nil {
			log.Printf("Failed to scan product row, continung to next one: %v", err)
			continue
		}

		if lastChecked.Valid {
			p.LastChecked = lastChecked.Time
		}

		products = append(products, p)
	}
	//checks for eeror if happened during iteration..
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *StorageService) UpdatePrice(id int, newPrice float64, alertSent bool) error {
	updateQuery := `
		UPDATE products
		SET last_price = ?, last_checked = ?, alert_sent = ?
		WHERE id = ?
	`
	_, err := s.db.Exec(updateQuery, newPrice, time.Now(), alertSent, id)
	if err != nil {
		log.Printf("Failed to update price for product %d: %v", id, err)
		return err
	}

	return nil
}
