CREATE TABLE IF NOT EXISTS subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    priority TEXT CHECK ( priority IN ('A', 'B', 'C', 'D') ),
    due TEXT,
    estimate_min INTEGER CHECK (estimate_min > 0),
    subject_id INTEGER REFERENCES subjects(id),
    project_id INTEGER REFERENCES projects(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    done_at TEXT
);

CREATE TABLE IF NOT EXISTS habits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    goal INTEGER NOT NULL CHECK (goal > 0) DEFAULT (1),
    freq_kind TEXT NOT NULL CHECK (freq_kind IN ('daily','weekdays','times_per_week')),
    freq_value TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    archived INTEGER NOT NULL CHECK (archived IN (0,1)) DEFAULT (0)
);

CREATE TABLE IF NOT EXISTS habit_logs (
    habit_id INTEGER NOT NULL REFERENCES habits(id),
    date TEXT NOT NULL DEFAULT (date('now','localtime')),
    value INTEGER NOT NULL,
    PRIMARY KEY (habit_id,date)
);


