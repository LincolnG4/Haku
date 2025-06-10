package store

import "time"

type Task struct {
	ID      int64         `json:"id"`
	Name    string        `json:"name"`
	Edges   []*Task       `json:"edges"`
	Job     string        `json:"job"`
	Config  *Config       `json:config`
	Timeout time.Duration `json:"timeout"`
}

type Config struct {
	SourcePath string `json: source_path`
	TargetPath string `json: target_path`
}
