package store

import (
	"database/sql"
	"log"
	"time"

	"goparking/models"

	_ "github.com/mattn/go-sqlite3"
)

// Store wraps the SQL DB and exposes helper methods.
type Store struct {
	DB *sql.DB
}

// NewStore opens the DB and creates tables if needed.
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


CREATE TABLE IF NOT EXISTS history (
  time INTEGER PRIMARY KEY,
  slot_id TEXT NOT NULL,
  occupied INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS commands (
  key TEXT PRIMARY KEY,
  action TEXT NOT NULL,
  threshold INTEGER NOT NULL,
  last_update INTEGER NOT NULL
);

INSERT OR IGNORE INTO commands (key, action, threshold, last_update)
VALUES ('main', 'none', 30, strftime('%s','now')*1000);
`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal("init schema:", err)
	}

	return &Store{DB: db}
}

// ----------------- SLOT CRUD ---------------------

func (s *Store) UpsertSlot(slot models.ParkingSlot) error {
	_, err := s.DB.Exec(`
        INSERT INTO slots (id, distance, status, last_update)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
          distance=excluded.distance,
          status=excluded.status,
          last_update=excluded.last_update
    `, slot.ID, slot.Distance, slot.Status, slot.LastUpdate)
	return err
}

func (s *Store) GetSlot(id string) (*models.ParkingSlot, error) {
	row := s.DB.QueryRow(`SELECT id, distance, status, last_update FROM slots WHERE id=?`, id)
	var slot models.ParkingSlot
	if err := row.Scan(&slot.ID, &slot.Distance, &slot.Status, &slot.LastUpdate); err != nil {
		return nil, err
	}
	return &slot, nil
}

func (s *Store) ListSlots() ([]models.ParkingSlot, error) {
	rows, err := s.DB.Query(`SELECT id, distance, status, last_update FROM slots ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.ParkingSlot
	for rows.Next() {
		var slot models.ParkingSlot
		if err := rows.Scan(&slot.ID, &slot.Distance, &slot.Status, &slot.LastUpdate); err != nil {
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

// ----------------- HISTORY ---------------------

func (s *Store) AddHistoryPoint(slotID string, occupied int) error {
	ts := time.Now().UnixMilli()
	_, err := s.DB.Exec(`
        INSERT INTO history (time, slot_id, occupied)
        VALUES (?, ?, ?)
    `, ts, slotID, occupied)
	return err
}

// ----------------- COMMANDS ---------------------

func (s *Store) GetCommand() (*models.Command, error) {
	row := s.DB.QueryRow(`SELECT action, threshold FROM commands WHERE key='main'`)
	var cmd models.Command
	if err := row.Scan(&cmd.Action, &cmd.Threshold); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (s *Store) SetCommand(cmd models.Command) error {
	ts := time.Now().UnixMilli()
	_, err := s.DB.Exec(`
        UPDATE commands
           SET action=?, threshold=?, last_update=?
         WHERE key='main'
    `, cmd.Action, cmd.Threshold, ts)
	return err
}

// Update only the action (keep existing threshold)
func (s *Store) SetCommandAction(action string) error {
	ts := time.Now().UnixMilli()
	_, err := s.DB.Exec(`
        UPDATE commands
           SET action=?, last_update=?
         WHERE key='main'
    `, action, ts)
	return err
}
