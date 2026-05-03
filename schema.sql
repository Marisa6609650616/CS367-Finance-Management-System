-- ============================================
-- CS367 Finance API — Database Schema
-- ออกแบบโดย: จันทร์พงศ์
-- Database: SQLite
-- ============================================

PRAGMA foreign_keys = ON;

-- ============================================
-- ตาราง: users
-- เก็บข้อมูลผู้ใช้งานระบบ (ทั้ง user และ admin)
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    email       TEXT    NOT NULL UNIQUE,
    password    TEXT    NOT NULL,               -- bcrypt hashed
    name        TEXT    NOT NULL,
    role        TEXT    NOT NULL DEFAULT 'user' CHECK(role IN ('user', 'admin')),
    created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

-- ============================================
-- ตาราง: categories
-- หมวดหมู่รายรับ-รายจ่าย (preset + custom)
-- ============================================
CREATE TABLE IF NOT EXISTS categories (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL,
    type        TEXT    NOT NULL CHECK(type IN ('income', 'expense')),
    user_id     INTEGER,                        -- NULL = preset ของระบบ, มีค่า = custom ของ user นั้น
    created_at  TEXT    NOT NULL DEFAULT (datetime('now')),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================
-- ตาราง: transactions
-- บันทึกรายรับ-รายจ่ายของผู้ใช้แต่ละคน
-- ============================================
CREATE TABLE IF NOT EXISTS transactions (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL,
    category_id     INTEGER NOT NULL,
    type            TEXT    NOT NULL CHECK(type IN ('income', 'expense')),
    amount          REAL    NOT NULL CHECK(amount > 0),
    note            TEXT,                       -- หมายเหตุ (optional)
    date            TEXT    NOT NULL,           -- วันที่เกิดรายการ เช่น '2025-04-28'
    created_at      TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT    NOT NULL DEFAULT (datetime('now')),

    FOREIGN KEY (user_id)     REFERENCES users(id)      ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

-- ============================================
-- Indexes: เพิ่มประสิทธิภาพ query
-- ============================================
CREATE INDEX IF NOT EXISTS idx_transactions_user_id   ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date       ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_type       ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_categories_user_id      ON categories(user_id);

-- ============================================
-- Seed Data: หมวดหมู่ preset ของระบบ
-- (user_id = NULL หมายความว่าเป็นของระบบ ทุกคนใช้ได้)
-- ============================================

-- หมวดหมู่รายจ่าย
INSERT OR IGNORE INTO categories (name, type, user_id) VALUES
    ('อาหารและเครื่องดื่ม', 'expense', NULL),
    ('ค่าเดินทาง',           'expense', NULL),
    ('ที่พักอาศัย',           'expense', NULL),
    ('สุขภาพและยา',          'expense', NULL),
    ('ช้อปปิ้ง',              'expense', NULL),
    ('บันเทิง',               'expense', NULL),
    ('การศึกษา',              'expense', NULL),
    ('อื่นๆ (รายจ่าย)',       'expense', NULL);

-- หมวดหมู่รายรับ
INSERT OR IGNORE INTO categories (name, type, user_id) VALUES
    ('เงินเดือน',             'income', NULL),
    ('รายได้เสริม',           'income', NULL),
    ('โบนัส',                 'income', NULL),
    ('อื่นๆ (รายรับ)',        'income', NULL);
