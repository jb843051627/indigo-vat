package store

const Schema = `
CREATE TABLE IF NOT EXISTS vats (id TEXT PRIMARY KEY,name TEXT NOT NULL,site TEXT NOT NULL,state TEXT NOT NULL,capacity INTEGER NOT NULL,target_ph REAL NOT NULL,target_temp REAL NOT NULL,timezone TEXT NOT NULL,created_at TEXT NOT NULL,updated_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS recipes (id TEXT PRIMARY KEY,name TEXT NOT NULL,target_hue REAL NOT NULL,hue_tolerance REAL NOT NULL,min_minutes INTEGER NOT NULL,max_minutes INTEGER NOT NULL,state TEXT NOT NULL,created_at TEXT NOT NULL,updated_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS stages (id TEXT PRIMARY KEY,recipe_id TEXT NOT NULL,position INTEGER NOT NULL,name TEXT NOT NULL,minutes INTEGER NOT NULL,ph_min REAL NOT NULL,ph_max REAL NOT NULL,temp_min REAL NOT NULL,temp_max REAL NOT NULL);
CREATE TABLE IF NOT EXISTS cycles (id TEXT PRIMARY KEY,vat_id TEXT NOT NULL,recipe_id TEXT NOT NULL,state TEXT NOT NULL,revision INTEGER NOT NULL,started_at TEXT NOT NULL,matured_at TEXT NOT NULL,released_at TEXT NOT NULL,updated_at TEXT NOT NULL,note TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS samples (id TEXT PRIMARY KEY,cycle_id TEXT NOT NULL,taken_at TEXT NOT NULL,hue REAL NOT NULL,ph REAL NOT NULL,temperature REAL NOT NULL,status TEXT NOT NULL,observer TEXT NOT NULL,note TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS inspections (id TEXT PRIMARY KEY,cycle_id TEXT NOT NULL,kind TEXT NOT NULL,result TEXT NOT NULL,score REAL NOT NULL,inspector TEXT NOT NULL,created_at TEXT NOT NULL,completed_at TEXT NOT NULL,note TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS alerts (id TEXT PRIMARY KEY,cycle_id TEXT NOT NULL,level TEXT NOT NULL,code TEXT NOT NULL,message TEXT NOT NULL,state TEXT NOT NULL,created_at TEXT NOT NULL,acknowledged_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS audit_events (id TEXT PRIMARY KEY,entity_type TEXT NOT NULL,entity_id TEXT NOT NULL,action TEXT NOT NULL,detail TEXT NOT NULL,created_at TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_cycles_state ON cycles(state);
CREATE INDEX IF NOT EXISTS idx_samples_cycle ON samples(cycle_id,taken_at);
CREATE INDEX IF NOT EXISTS idx_alerts_cycle ON alerts(cycle_id,state);
`
