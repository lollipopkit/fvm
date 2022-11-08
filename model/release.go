package model

import "time"

type AllReleases struct {
	BaseURL        string `json:"base_url"`
	CurrentRelease struct {
		Beta   string `json:"beta"`
		Dev    string `json:"dev"`
		Stable string `json:"stable"`
	} `json:"current_release"`
	Releases []Release `json:"releases"`
}

type Release struct {
	Hash           string    `json:"hash"`
	Channel        string    `json:"channel"`
	Version        string    `json:"version"`
	DartSdkVersion string    `json:"dart_sdk_version,omitempty"`
	DartSdkArch    string    `json:"dart_sdk_arch,omitempty"`
	ReleaseDate    time.Time `json:"release_date"`
	Archive        string    `json:"archive"`
	Sha256         string    `json:"sha256"`
}
