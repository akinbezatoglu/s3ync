package config

import (
	"strconv"

	"github.com/akinbezatoglu/s3ync/pkg/config"
)

type Config interface {
	GetProfileNames() ([]string, error)
	GetSyncListFromProfile(profile string) ([]string, error)
	AddSync(profile, path, bucketname, bucketregion string) error
	RemoveSync(profile, path string) error
	Write() error
	Set(keys []string, value string)
	GetConfigFilePath() string
	GetRegionFromProfile(profile string) (string, error)
	IsProfileExistInConfigFile(p string) bool
}

type cfg struct {
	cfg *config.Config
}

func NewConfig() (Config, error) {
	c, err := config.Read(fallbackConfig())
	if err != nil {
		return nil, err
	}
	return &cfg{c}, nil
}

func (c *cfg) GetConfigFilePath() string {
	return config.GeneralConfigFile()
}

func (c *cfg) Write() error {
	return config.Write(c.cfg)
}

func (c *cfg) IsProfileExistInConfigFile(p string) bool {
	profiles, _ := c.GetProfileNames()
	for _, profile := range profiles {
		if p == profile {
			return true
		}
	}
	return false
}

func (c *cfg) GetProfileNames() ([]string, error) {
	return c.cfg.Keys([]string{"s3", "profiles"})
}

func (c *cfg) GetRegionFromProfile(profile string) (string, error) {
	return c.cfg.Get([]string{"s3", "profiles", profile, "region"})
}

func (c *cfg) GetSyncListFromProfile(profile string) ([]string, error) {
	idx, err := c.cfg.Keys([]string{"s3", "profiles", profile, "syncs"})
	if err != nil {
		return nil, err
	}
	if len(idx) == 0 {
		return nil, nil
	}
	var x []string
	for _, i := range idx {
		local, err := c.cfg.Get([]string{"s3", "profiles", profile, "syncs", i, "local"})
		if err != nil {
			return nil, err
		}
		bucket, err := c.cfg.Get([]string{"s3", "profiles", profile, "syncs", i, "bucket", "name"})
		if err != nil {
			return nil, err
		}
		x = append(x, local, bucket)
	}
	return x, nil
}

func (c *cfg) AddSync(profile, path, bucketname, bucketregion string) error {
	idx, err := c.cfg.Keys([]string{"s3", "profiles", profile, "syncs"})
	if err != nil {
		return err
	}
	c.cfg.Set([]string{"s3", "profiles", profile, "syncs", strconv.Itoa(len(idx) + 1), "local"}, path)
	c.cfg.Set([]string{"s3", "profiles", profile, "syncs", strconv.Itoa(len(idx) + 1), "bucket", "name"}, bucketname)
	c.cfg.Set([]string{"s3", "profiles", profile, "syncs", strconv.Itoa(len(idx) + 1), "bucket", "region"}, bucketregion)
	c.Write()
	return nil
}

func (c *cfg) RemoveSync(profile, path string) error {
	idx, err := c.cfg.Keys([]string{"s3", "profiles", profile, "syncs"})
	if err != nil {
		return err
	}
	flg := false
	for l, id := range idx {
		p, err := c.cfg.Get([]string{"s3", "profiles", profile, "syncs", id, "local"})
		if err != nil {
			return err
		}
		if !flg {
			if p == path {
				err := c.cfg.Remove([]string{"s3", "profiles", profile, "syncs", id})
				if err != nil {
					return err
				}
				if len(idx) == 1 {
					c.cfg.Set([]string{"s3", "profiles", profile, "syncs"}, "")
				}
			}
			flg = true
		} else {
			local, err := c.cfg.Get([]string{"s3", "profiles", profile, "syncs", id, "local"})
			if err != nil {
				return err
			}
			bucketname, err := c.cfg.Get([]string{"s3", "profiles", profile, "syncs", id, "bucket", "name"})
			if err != nil {
				return err
			}
			bucketregion, err := c.cfg.Get([]string{"s3", "profiles", profile, "syncs", id, "bucket", "region"})
			if err != nil {
				return err
			}
			c.cfg.Set([]string{"s3", "profiles", profile, "syncs", strconv.Itoa(l), "local"}, local)
			c.cfg.Set([]string{"s3", "profiles", profile, "syncs", strconv.Itoa(l), "bucket", "name"}, bucketname)
			c.cfg.Set([]string{"s3", "profiles", profile, "syncs", strconv.Itoa(l), "bucket", "region"}, bucketregion)
			err = c.cfg.Remove([]string{"s3", "profiles", profile, "syncs", id})
			if err != nil {
				return err
			}
		}
	}
	c.Write()
	return nil
}

func (c *cfg) Set(keys []string, value string) {
	c.cfg.Set(keys, value)
}

func fallbackConfig() *config.Config {
	return config.ReadFromString(defaultConfigStr)
}

const defaultConfigStr = `
s3:
  profiles:
`
