CREATE TABLE subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT UNIQUE
);

CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT UNIQUE
);

CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    priority TEXT CHECK ( priority IN ('A', 'B', 'C', 'D') ),
    due TEXT,
    estimate_min INTEGER CHECK (estimate_min > 0),
    subject_id INTEGER REFERENCES subjects(id),
    project_id INTEGER REFERENCES projects(id),
    parent_id INTEGER REFERENCES tasks(id),
    extras TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime')),
    done_at TEXT
);

CREATE TABLE checklist_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    done INTEGER NOT NULL CHECK (done IN (0,1)) DEFAULT (0)
);

CREATE TABLE habits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL CHECK (kind IN ('bool','count')),
    goal INTEGER CHECK (goal > 0),
    freq_kind TEXT NOT NULL CHECK (freq_kind IN ('daily','weekdays','times_per_week')),
    freq_value TEXT,
    slack INTEGER NOT NULL CHECK (slack >= 0) DEFAULT (0),
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    archived INTEGER NOT NULL CHECK (archived IN (0,1)) DEFAULT (0)
);

CREATE TABLE habit_logs (
    habit_id INTEGER NOT NULL REFERENCES habits(id),
    date TEXT NOT NULL DEFAULT (date('now','localtime')),
    value INTEGER NOT NULL,
    PRIMARY KEY (habit_id,date)
);


