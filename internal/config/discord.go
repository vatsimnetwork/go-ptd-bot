package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	ID                  string             `yaml:"id"`
	Name                string             `yaml:"name"`
	RatingRoles         []RatingRoleConfig `yaml:"rating_roles"`
	PilotRatingRoles    []RatingRoleConfig `yaml:"pilot_rating_roles"`
	MilitaryRatingRoles []RatingRoleConfig `yaml:"military_rating_roles"`
}

type RatingRoleConfig struct {
	ShortName string `yaml:"name"`
	RoleID    string `yaml:"id"`
	CertValue int    `yaml:"cert_value"`
}

var (
	ConfigPath   = os.Getenv("CONFIG_PATH")
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	Env          = os.Getenv("GO_ENV")
	SentryDSN    = os.Getenv("SENTRY_DSN")
)

func LoadAllServerConfigOrPanic(configPath string) map[string]ServerConfig {
	cfgs, err := LoadAllServerConfig(configPath)
	if err != nil {
		log.Printf(err.Error())
	}

	return cfgs
}

func LoadAllServerConfig(configPath string) (map[string]ServerConfig, error) {
	files, err := os.ReadDir(configPath)
	if err != nil {
		return nil, errors.New("failed to load server configs")
	}

	cfgs := make(map[string]ServerConfig, len(files))
	for _, f := range files {
		if !f.IsDir() {
			cfg, err := LoadServerConfig(fmt.Sprintf("%s/%s", configPath, f.Name()))
			if err != nil {
				sentry.CaptureException(err)
				log.Printf("failed to load config for %s: %v", f.Name(), err)
				continue
			}
			cfgs[cfg.ID] = *cfg
		}
	}

	return cfgs, nil
}

func LoadServerConfig(configPath string) (*ServerConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	// TODO: Validate that roles aren't duplicated
	// TODO: Validate role criteria
	log.Printf("Loaded Config for %s (%s)\n", cfg.Name, cfg.ID)

	return &cfg, nil
}

var configs = LoadAllServerConfigOrPanic(ConfigPath)

func IntervalReloadConfigs() {
	for {
		time.Sleep(5 * time.Minute)
		log.Print("Reloading server configs")
		configs = LoadAllServerConfigOrPanic(ConfigPath)
	}
}

func GetRatingsRoles(id string) []RatingRoleConfig {
	cfg, _ := configs[id]
	return cfg.RatingRoles
}

func GetPilotRatingRoles(id string) []RatingRoleConfig {
	cfg, _ := configs[id]
	return cfg.PilotRatingRoles
}

func GetMilitaryRatingRoles(id string) []RatingRoleConfig {
	cfg, _ := configs[id]
	return cfg.MilitaryRatingRoles
}
