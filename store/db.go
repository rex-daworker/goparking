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
CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    message TEXT NOT NULL,
    timestamp INTEGER NOT NULL
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

// In store/db.go
func (s *Store) GetCommand() (*models.Command, error) {
	row := s.DB.QueryRow(`SELECT action, threshold FROM commands WHERE key='main'`)
	var cmd models.Command
	if err := row.Scan(&cmd.Action, &cmd.Threshold); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (s *Store) SetCommand(cmd models.Command) error {
	_, err := s.DB.Exec(`
        INSERT INTO commands (key, action, threshold)
        VALUES ('main', ?, ?)
        ON CONFLICT(key) DO UPDATE SET
            action=excluded.action,
            threshold=excluded.threshold
    `, cmd.Action, cmd.Threshold)
	return err
}

func (s *Store) AddEvent(event models.Event) error {
    _, err := s.DB.Exec(`
        INSERT INTO events (id, type, message, timestamp)
        VALUES (?, ?, ?, ?)
    `, event.ID, event.Type, event.Message, event.Timestamp)
    return err
}

func (s *Store) ListEvents() ([]models.Event, error) {
    rows, err := s.DB.Query(`SELECT id, type, message, timestamp FROM events`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var events []models.Event
    for rows.Next() {
        var e models.Event
        if err := rows.Scan(&e.ID, &e.Type, &e.Message, &e.Timestamp); err != nil {
            return nil, err
        }
        events = append(events, e)
    }
    return events, nil
}
