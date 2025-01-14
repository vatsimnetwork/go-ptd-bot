package config

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)
import _ "github.com/joho/godotenv/autoload"

type ServerConfig struct {
	Id                  string             `yaml:"id"`
	Name                string             `yaml:"name"`
	RatingRoles         []RatingRoleConfig `yaml:"rating_roles"`
	PilotRatingRoles    []RatingRoleConfig `yaml:"pilot_rating_roles"`
	MilitaryRatingRoles []RatingRoleConfig `yaml:"military_rating_roles"`
}
type RatingRoleConfig struct {
	ShortName     string `yaml:"name"`
	DiscordRoleId string `yaml:"id"`
	CertValue     int    `yaml:"cert_value"`
}

var DiscordToken = os.Getenv("DISCORD_TOKEN")
var SentryDSN = os.Getenv("SENTRY_DSN")
var Env = os.Getenv("GO_ENV")
var ConfigPath = os.Getenv("CONFIG_PATH")

func LoadAllServerConfigOrPanic(configPath string) map[string]ServerConfig {
	cfgs, err := LoadAllServerConfig(configPath)
	if err != nil {
		log.Printf(err.Error())
	}
	return cfgs
}

func LoadAllServerConfig(configPath string) (map[string]ServerConfig, error) {
	cfgs := make(map[string]ServerConfig, 0)
	files, err := os.ReadDir(configPath)
	if err != nil {
		return nil, errors.New("failed to load server configs")
	}
	for _, f := range files {
		if !f.IsDir() {
			cfg, err := LoadServerConfig(fmt.Sprintf("%s/%s", configPath, f.Name()))
			if err != nil {
				log.Printf(err.Error())
				return nil, nil
			}
			cfgs[cfg.Id] = *cfg
		}
	}
	return cfgs, nil
}

func LoadServerConfig(configPath string) (*ServerConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	var cfg ServerConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	// TODO: Validate that roles aren't duplicated
	// TODO: Validate role criteria
	log.Printf("Loaded Config for %s (%s)\n", cfg.Name, cfg.Id)
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

	Cfg, _ := configs[id]
	return Cfg.RatingRoles
}

func GetPilotRatingRoles(id string) []RatingRoleConfig {
	Cfg, _ := configs[id]
	return Cfg.PilotRatingRoles
}

func GetMilitaryRatingRoles(id string) []RatingRoleConfig {
	Cfg, _ := configs[id]
	return Cfg.MilitaryRatingRoles
}
