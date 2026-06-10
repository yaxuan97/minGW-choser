package main

import (
	_ "embed"
	"encoding/json"
	"mingw-chooser/match"
)

//go:embed builds.json
var buildsJSON []byte

type configFile struct {
	Version        int            `json:"version"`
	Sources        []sourceConfig `json:"sources"`
	Rules          match.Rules    `json:"rules"`
	FallbackBuilds []match.Build  `json:"fallback_builds"`
}

type sourceConfig struct {
	Name        string `json:"name"`
	API         string `json:"api"`
	FallbackURL string `json:"fallback_url"`
}

func loadConfig() (configFile, error) {
	var cfg configFile
	if err := json.Unmarshal(buildsJSON, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
