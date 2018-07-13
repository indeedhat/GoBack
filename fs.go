package main

type BackupEntry struct {
  Path    string
  Name    string
  Date    int64
  ModTime int64
  Hash    string
  Size    int64
}