package store

import (
	"database/sql"
	"goparking/models"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	DB *sql.DB
}

func NewStore(path string) *Store {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("open db:", err)
	}

	schema := `
    CREATE TABLE IF NOT EXISTS slots (
        id TEXT PRIMARY KEY,
        distance INTEGER NOT NULL,
        status TEXT NOT NULL,
        last_update INTEGER NOT NULL,
        device_id TEXT,
        device_name TEXT,
        sensor_status TEXT
    );
    `
	if _, err := db.Exec(schema); err != nil {
		log.Fatal("create schema:", err)
	}

	return &Store{DB: db}
}

func (s *Store) UpsertSlot(slot models.ParkingSlot) error {
	_, err := s.DB.Exec(`
        INSERT INTO slots (id, distance, status, last_update, device_id, device_name, sensor_status)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            distance=excluded.distance,
            status=excluded.status,
            last_update=excluded.last_update,
            device_id=excluded.device_id,
            device_name=excluded.device_name,
            sensor_status=excluded.sensor_status
    `, slot.ID, slot.Distance, slot.Status, slot.LastUpdate, slot.DeviceID, slot.DeviceName, slot.SensorStatus)
	return err
}

func (s *Store) GetSlot(id string) (*models.ParkingSlot, error) {
	row := s.DB.QueryRow(`SELECT id, distance, status, last_update, device_id, device_name, sensor_status FROM slots WHERE id=?`, id)
	var slot models.ParkingSlot
	if err := row.Scan(&slot.ID, &slot.Distance, &slot.Status, &slot.LastUpdate, &slot.DeviceID, &slot.DeviceName, &slot.SensorStatus); err != nil {
		return nil, err
	}
	return &slot, nil
}

func (s *Store) ListSlots() ([]models.ParkingSlot, error) {
	rows, err := s.DB.Query(`SELECT id, distance, status, last_update, device_id, device_name, sensor_status FROM slots ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.ParkingSlot
	for rows.Next() {
		var slot models.ParkingSlot
		if err := rows.Scan(&slot.ID, &slot.Distance, &slot.Status, &slot.LastUpdate, &slot.DeviceID, &slot.DeviceName, &slot.SensorStatus); err != nil {
			return nil, err
		}
		res = append(res, slot)
	}
	return res, nil
}

func (s *Store) DeleteSlot(id string) error {
	_, err := s.DB.Exec(`DELETE FROM slots WHERE id=?`, id)
	return err
}
