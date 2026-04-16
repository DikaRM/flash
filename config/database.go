package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Conn() *sql.DB {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		panic(err)
	}
	log.Println("Database Berhasil Dibuat")
	DB = db
	query := `
-- Table modules
CREATE TABLE IF NOT EXISTS modules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama_modules TEXT,
    deskripsi TEXT,
    file TEXT,
    badge TEXT,
    total_rating INTEGER DEFAULT 0,
    avg_rating REAL DEFAULT 0,
    author INTEGER,
    view INTEGER DEFAULT 0,
    thumbnail TEXT,
    kategori TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table users
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE,
    password TEXT DEFAULT NULL,
    bio TEXT DEFAULT Programmer,
    level INTEGER DEFAULT 1,
    badge TEXT DEFAULT stator,
    avatar TEXT,
    email TEXT UNIQUE NOT NULL,
    google_id TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table ratings
CREATE TABLE IF NOT EXISTS ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    module_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    rating INTEGER CHECK(rating >= 1 AND rating <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(module_id) REFERENCES modules(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, module_id)
);

-- Table comments
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    module_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(module_id) REFERENCES modules(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index untuk comments (dibuat terpisah)
CREATE INDEX IF NOT EXISTS idx_comments_module_id ON comments(module_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);

-- Index untuk ratings
CREATE INDEX IF NOT EXISTS idx_ratings_module_id ON ratings(module_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user_id ON ratings(user_id);
`
	_, eror := db.Exec(query)
	if eror != nil {
		log.Fatal(eror)
	}
	log.Println("Tabel Modules Berhasil Di BUAT ")
	return db
}
