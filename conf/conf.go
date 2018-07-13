package conf

import (
  "os"
  "gopkg.in/yaml.v2"
  "io/ioutil"
)

var cnfCache *Config

type Config struct {
  Dirs   []string `yaml:",flow"`
  OutDir string   `yaml:"out_dir"`
}


func (c *Config) ValidPaths() []string {
  valid := make([]string, 1)

  for _, path := range c.Dirs {
    if _, err := os.Stat(path); nil != err {
      continue
    }

    valid = append(valid, path)
  }

  return valid
}

func (c *Config) VerifyConfig() bool {
  paths := c.ValidPaths()

  if 0 == len(paths) {
    return false
  }

  if _, err := os.Stat(c.OutDir); nil != err {
    return false
  }

  return true
}


func Load() (*Config, error) {
  if nil != cnfCache {
    return cnfCache, nil
  }

  cnfCache = &Config{}
  yml, err := ioutil.ReadFile("config.yml")

  if nil != err {
    return nil, err
  }

  err = yaml.Unmarshal(yml, cnfCache)

  if nil != err {
    return nil, err
  }

  return cnfCache, nil
}


