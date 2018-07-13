package main

import (
  "sync"
  "path/filepath"
  "os"
  "time"
  "crypto/md5"
  "fmt"
  "./util"
  "./conf"
)

var scanGroup   sync.WaitGroup
var checkGroup  sync.WaitGroup
var backupGroup sync.WaitGroup

var checkQueue  chan *BackupEntry
var backupQueue chan *BackupEntry

func workerInit() {
  checkQueue  = make(chan *BackupEntry, 50)
  backupQueue = make(chan *BackupEntry, 50)
}

func scanWorker(dir string) {
  defer scanGroup.Done()
  scanGroup.Add(1)

  filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
    byts := []byte(fmt.Sprintf("%s/%s_::_%d", path, f.Name(), f.ModTime().Unix()))


    checkQueue <- &BackupEntry{
      Date: time.Now().Unix(),
      Hash: string(md5.Sum(byts)),
      ModTime: f.ModTime().Unix(),
      Name: f.Name(),
      Path: path,
      Size: f.Size(),
    }

    return nil
  })

  // TODO: find a way of allowing for multiple scan workers
  close(checkQueue)
}

func checkWorker() {
  defer checkGroup.Done()
  checkGroup.Add(1)

  db := dbOpen()
  for entry := range checkQueue {
    modTime, err := db.GetLatestBackupTime(entry.Hash)

    if nil != err {
      // TODO: do some logging or something here
      continue
    }

    if modTime != entry.ModTime {
      backupQueue <- entry
    }
  }

  // TODO: find a way of allowing for multiple check workers
  close(backupQueue)
}

func backupWorker() {
  defer wg.Done()
  wg.Add(1)

  db := dbOpen()
  cnf, err := conf.Load()
  if nil != err {
    return
  }

  for entry := range backupQueue {
    // ensure the dir exists
    os.MkdirAll( fmt.Sprintf("%s/%s", cnf.OutDir, entry.Path), os.ModePerm)

    // do the copy for now i dont care about errors
    util.CopyFile(
      fmt.Sprintf("%s/%s.%d", entry.Path, entry.Name, entry.Date),
      fmt.Sprintf("%s/%s/%s.%d", cnf.OutDir, entry.Path, entry.Name, entry.Date),
    )

    // add entry to the history database
    db.AddFile(entry)
  }
}