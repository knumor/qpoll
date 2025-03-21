package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"runtime"
	"time"

	"github.com/knumor/qpoll/handlers"
	"github.com/knumor/qpoll/models"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct {
	readDB  *sql.DB
	writeDB *sql.DB
}

// NewSQLiteStore creates a new sqlite storage.
func NewSQLiteStore() handlers.Storage {
	s := &sqlite{}
	s.readDB, s.writeDB = initDatabases()
	verifyOrCreateTables(s.writeDB)
	return s
}

func initDatabases() (*sql.DB, *sql.DB) {
	connParams := make(url.Values)
	connParams.Add("_txlock", "immediate")
	connParams.Add("_journal_mode", "WAL")
	connParams.Add("_busy_timeout", "5000")
	connParams.Add("_synchronous", "NORMAL")
	connParams.Add("_cache_size", "1000000000")
	connParams.Add("_foreign_keys", "true")
	connURL := "file:db/qpoll.db?" + connParams.Encode()

	readDB, err := sql.Open("sqlite3", connURL)
	if err != nil {
		panic("Failed to open read database")
	}
	readDB.SetMaxOpenConns(max(4, runtime.NumCPU()))
	writeDB, err := sql.Open("sqlite3", connURL)
	if err != nil {
		panic("Failed to open write database")
	}
	writeDB.SetMaxOpenConns(1)

	for _, pragma := range getPragmas() {
		writeDB.Exec("PRAGMA " + pragma)
		readDB.Exec("PRAGMA " + pragma)
	}

	return readDB, writeDB
}

func getPragmas() []string {
	return []string{
		"temp_store=memory",
	}
}

func verifyOrCreateTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS polls (id TEXT PRIMARY KEY, code TEXT UNIQUE, type INT, data TEXT) STRICT")
	if err != nil {
		slog.Error("Failed to create polls table", "error", err)
		panic("Failed to create polls table")
	}
}

// Save saves the poll to the storage.
func (s *sqlite) Save(p models.Poll) error {
	start := time.Now()
	slog.Info("Save", "id", p.ID())
	slog.Info("Save", "code", p.Code())
	data, _ := p.MarshalJSON()
	_, err := s.writeDB.Exec(
		"INSERT OR REPLACE INTO polls (id, code, type, data) VALUES (?, ?, ?, ?)",
		p.ID(),
		p.Code(),
		p.Type(),
		string(data),
	)
	if err != nil {
		slog.Error("Failed to save poll", "error", err)
		return fmt.Errorf("failed to save poll: %w", err)
	}
	elapsed := time.Since(start)
	slog.Info("Save done", "elapsed", elapsed)
	return nil
}

// Load loads the poll from the storage.
func (s *sqlite) Load(id string) (models.Poll, error) {
	start := time.Now()
	slog.Info("Load", "id", id)
	row := s.readDB.QueryRow("SELECT code, type, data FROM polls WHERE id = ?", id)
	var code string
	var data string
	var pollType models.PollType
	err := row.Scan(&code, &pollType, &data)
	if err != nil {
		return nil, fmt.Errorf("poll with id %s not found", id)
	}
	elapsed := time.Since(start)
	p, err := createPollObject(id, pollType, data)
	slog.Info("Load done", "code", code, "type", pollType, "data", data, "elapsed", elapsed)
	return p, err
}

// LoadByCode loads the poll by code from the storage.
func (s *sqlite) LoadByCode(code string) (models.Poll, error) {
	start := time.Now()
	slog.Info("LoadByCode", "code", code)
	row := s.readDB.QueryRow("SELECT id, type, data FROM polls WHERE code = ?", code)
	var id string
	var data string
	var pollType models.PollType
	err := row.Scan(&id, &pollType, &data)
	if err != nil {
		return nil, fmt.Errorf("poll with code %s not found", code)
	}
	elapsed := time.Since(start)
	p, err := createPollObject(id, pollType, data)
	slog.Info("Load done", "code", code, "type", pollType, "data", data, "elapsed", elapsed)
	return p, err
}

// LoadAllByUser loads all polls by user from the storage.
func (s *sqlite) LoadAllByUser(username string) ([]models.Poll, error) {
	start := time.Now()
	slog.Info("LoadAllByUser", "username", username)
	rows, err := s.readDB.Query("SELECT id, type, data FROM polls")
	if err != nil {
		return nil, fmt.Errorf("failed to load polls for user %s: %w", username, err)
	}
	defer rows.Close()
	var polls []models.Poll
	for rows.Next() {
		var id string
		var data string
		var pollType models.PollType
		err = rows.Scan(&id, &pollType, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		p, err := createPollObject(id, pollType, data)
		if err != nil {
			return nil, fmt.Errorf("failed to create poll object: %w", err)
		}
		if p.Owner() != username {
			continue
		}
		polls = append(polls, p)
	}
	elapsed := time.Since(start)
	slog.Info("LoadAllByUser done", "username", username, "elapsed", elapsed)
	return polls, nil
}

func createPollObject(id string, pollType models.PollType, data string) (models.Poll, error) {
	switch pollType {
	case models.MultipleChoicePoll:
		return models.MultipleChoiceFromJSON(id, []byte(data))
	case models.WordCloudPoll:
		return models.WordCloudFromJSON(id, []byte(data))
	default:
		return nil, fmt.Errorf("invalid poll type %d", pollType)
	}
}

func (s *sqlite) Close() {
	s.readDB.Close()
	s.writeDB.Close()
}
