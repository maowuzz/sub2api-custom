package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type manifestPayload struct {
	SchemaVersion   int    `json:"schema_version"`
	Status          string `json:"status"`
	DisplayVersion  string `json:"display_version"`
	UpstreamVersion string `json:"upstream_version"`
	Image           string `json:"image"`
	ImageDigest     string `json:"image_digest"`
	UpstreamCommit  string `json:"upstream_commit"`
	CustomCommit    string `json:"custom_commit"`
	PublishedAt     string `json:"published_at"`
	ReleaseNotes    string `json:"release_notes,omitempty"`
	HTMLURL         string `json:"html_url,omitempty"`
}

type manifest struct {
	manifestPayload
	Signature string `json:"signature"`
}

func required(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(name + " is required")
	}
	return value
}

func main() {
	payload := manifestPayload{
		SchemaVersion:   1,
		Status:          "ready",
		DisplayVersion:  required("DISPLAY_VERSION"),
		UpstreamVersion: required("UPSTREAM_VERSION"),
		Image:           required("IMAGE"),
		ImageDigest:     required("IMAGE_DIGEST"),
		UpstreamCommit:  required("UPSTREAM_COMMIT"),
		CustomCommit:    required("CUSTOM_COMMIT"),
		PublishedAt:     required("PUBLISHED_AT"),
		ReleaseNotes:    os.Getenv("RELEASE_NOTES"),
		HTMLURL:         os.Getenv("HTML_URL"),
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	mac := hmac.New(sha256.New, []byte(required("CUSTOM_MANIFEST_HMAC_SECRET")))
	_, _ = mac.Write(encoded)
	result := manifest{manifestPayload: payload, Signature: hex.EncodeToString(mac.Sum(nil))}
	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
