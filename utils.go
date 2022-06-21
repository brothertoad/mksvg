package main

import (
  "io"
  "log"
  "os"
)

func fileExists(path string) bool {
  fileInfo, err := os.Stat(path)
  if err != nil {
    return false
  }
  if !fileInfo.Mode().IsRegular() {
    log.Fatal("%s exists, but is not a file\n", path)
  }
  return true
}

func dirExists(dir string) bool {
  fileInfo, err := os.Stat(dir)
  if err != nil {
    return false
  }
  if !fileInfo.IsDir() {
    log.Fatal("%s exists, but is not a directory\n", dir)
  }
  return true
}

func dirMustExist(dir string) {
  if !dirExists(dir) {
    log.Fatalf("%s does not exist\n", dir)
  }
}

func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

// Got this from https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file
// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }
    return out.Close()
}
