package main

import (
  "database/sql"
  "fmt"
)
type Database struct {
  Connection *sql.DB
}
var db Database

func dbOpen() Database {
  if (Database{}) != db {
    return db
  }

  con, err := sql.Open("sqlite3", "file:history.sqlite?cache=shared")

  if nil != err {
    return nil
  }

  con.SetMaxOpenConns(1)
  db = Database{con}

  return db
}

func (db *Database) Close() error {
  return db.Connection.Close()
}

func (db *Database) GetLatestBackupTime(fileHash string) (int64, error) {
  var t int64

  err := db.Connection.QueryRow(`
    SELECT
      mod_time
    FROM
      history
    WHERE
      hash = ?
    ORDER BY
      id DESC
    `, fileHash).Scan(&t)

  if nil != err {
    return nil, err
  }

  return t, nil
}

func (db *Database) AddFile(entry *BackupEntry) bool {
  stmt, err := db.Connection.Prepare(`
    INSERT INTO
      history
    (
      date,
      hash,
      mod_time,
      name,
      path,
      file_size
    ) VALUES (
      ?,
      ?,
      ?,
      ?,
      ?,
      ?
    )
  `)

  if nil != err {
    return false
  }

  rows, err := stmt.Exec(entry.Date, entry.Hash, entry.ModTime, entry.Name, entry.Path, entry.Size)

  if nil != err {
    return false
  }

  if count, err := rows.RowsAffected(); nil != err || 0 == count {
    return false
  }

  return true
}

func (db *Database) GetLatestFileEntry(hash string) (*BackupEntry, error) {
  entry := BackupEntry{}

  err := db.Connection.QueryRow(`
    SELECT
      date,
      hash,
      mod_time,
      name,
      path,
      file_size
    FROM
      history
    WHERE
      hash = ?
    ORDER BY
      id DESC
  `, hash).Scan(&entry.Date, &entry.Hash, &entry.ModTime, &entry.Name, &entry.Path, &entry.Size)

  if nil != err {
    return nil, err
  }

  return &entry, nil
}

func (db *Database) GetAllFileEntries(hash string) ([]*BackupEntry, error) {
  var entries []*BackupEntry

  rows, err := db.Connection.Query(`
    SELECT
      date,
      hash,
      mod_time,
      name,
      path,
      file_size
    FROM
      history
    WHERE
      hash = ?
    ORDER BY
      id DESC
  `, hash)

  if nil != err {
    return nil, err
  }

  for rows.Next() {
    entry := BackupEntry{}
    err := rows.Scan(&entry.Date, &entry.Hash, &entry.ModTime, &entry.Name, &entry.Path, &entry.Size)
    if nil != err {
      continue
    }

    entries = append(entries, &entry)
  }

  return entries, nil
}

func (db *Database) ListFiles(dir string) ([]*BackupEntry, error) {
  var entries []*BackupEntry

  rows, err := db.Connection.Query(`
    SELECT
      date,
      hash,
      mod_time,
      name,
      path
    FROM
      history
    WHERE
      path LIKE ?
    ORDER BY
      name ASC,
      id DESC
    GROUP BY
      id
  `, fmt.Sprintf("%s%%", dir))

  if nil != err {
    return nil, err
  }

  for rows.Next() {
    entry := BackupEntry{}
    err := rows.Scan(&entry.Date, &entry.Hash, &entry.ModTime, &entry.Name, &entry.Path)
    if nil != err {
      continue
    }

    entries = append(entries, &entry)
  }

  return entries, nil
}