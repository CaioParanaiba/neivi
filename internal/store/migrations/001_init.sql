-- disciplina: @BD (RN-21: ativa ou encerrada)
CREATE TABLE IF NOT EXISTS subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'closed')) DEFAULT ('active')
);

-- projeto: +neivi (D-06, D-36); tipo academico/pessoal eh deduzido do subject_id
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    github TEXT,
    dir TEXT,
    status TEXT NOT NULL CHECK (status IN ('active', 'paused', 'done', 'canceled')) DEFAULT ('active'),
    continuous INTEGER NOT NULL CHECK (continuous IN (0, 1)) DEFAULT (0),
    subject_id INTEGER REFERENCES subjects(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime'))
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

-- habito (RN-23): todo dia, dias especificos ou N por semana.
-- goal eh a meta do dia (daily, weekdays) ou da semana (times_per_week).
-- freq_value: so no weekdays, os dias como numeros do time.Weekday ("1,3,5" = seg, qua, sex)
CREATE TABLE IF NOT EXISTS habits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    goal INTEGER NOT NULL CHECK (goal > 0) DEFAULT (1),
    freq_kind TEXT NOT NULL CHECK (freq_kind IN ('daily', 'weekdays', 'times_per_week')) DEFAULT ('daily'),
    freq_value TEXT,
    status TEXT NOT NULL CHECK (status IN ('active', 'archived')) DEFAULT ('active'),
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- registro: um por habito e dia; a quantidade do dia soma (RN-24)
CREATE TABLE IF NOT EXISTS habit_logs (
    habit_id INTEGER NOT NULL REFERENCES habits(id),
    date TEXT NOT NULL,
    value INTEGER NOT NULL CHECK (value > 0),
    PRIMARY KEY (habit_id, date)
);
