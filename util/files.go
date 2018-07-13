package util

import (
  "os"
  "fmt"
)

// this is not my code i pulled it from stackoverflow, i may rewrite it at some point if i feel like it
// credit goes to:
// markc : https://stackoverflow.com/a/21067803
func CopyFile(src, dst string) (err error) {
  sfi, err := os.Stat(src)
  if err != nil {
    return
  }
  if !sfi.Mode().IsRegular() {
    // cannot copy non-regular files (e.g., directories,
    // symlinks, devices, etc.)
    return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
  }
  dfi, err := os.Stat(dst)
  if err != nil {
    if !os.IsNotExist(err) {
      return
    }
  } else {
    if !(dfi.Mode().IsRegular()) {
      return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
    }
    if os.SameFile(sfi, dfi) {
      return
    }
  }
  if err = os.Link(src, dst); err == nil {
    return
  }
  err = copyFileContents(src, dst)
  return
}