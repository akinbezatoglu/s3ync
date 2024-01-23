package s3

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/ini.v1"
)

// return: [user1, region1, user2, region2, ...]
func GetLocalAwsProfilesWithDefaultRegion() []string {
	fname := config.DefaultSharedConfigFilename()
	f, err := ini.Load(fname)
	if err != nil {
		return nil
	}
	profiles := f.SectionStrings()

	// There is no profile in config file
	if len(profiles) == 1 && profiles[0] == "DEFAULT" {
		return nil
	}

	profiles = profiles[1:] // remove DEFAULT
	var x []string
	for _, profile := range profiles {
		region := f.Section(profile).Key("region").String()
		x = append(x, profile, region)
	}
	return x
}

func IsAwsCliConfigured() bool {
	fnames := config.DefaultSharedConfigFiles
	if len(fnames) != 0 {
		for _, fname := range fnames {
			f, err := ini.Load(fname)
			if err != nil {
				return false
			}
			profiles := f.SectionStrings()
			if len(profiles) == 1 && profiles[0] == "DEFAULT" {
				return false
			}
		}
		return true
	}
	return false
}

func IsProfileExistInAwsConfig(p string) bool {
	fnames := config.DefaultSharedConfigFiles
	if len(fnames) != 0 {
		for _, fname := range fnames {
			f, err := ini.Load(fname)
			if err != nil {
				return false
			}
			profiles := f.SectionStrings()
			if len(profiles) == 1 && profiles[0] == "DEFAULT" {
				return false
			}

			profiles = profiles[1:] // remove DEFAULT
			for _, profile := range profiles {
				if p == profile {
					// found
					return true
				}
			}
		}
		return false
	}
	return false
}
