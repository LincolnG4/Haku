package store

type Task struct {
	ID          int64  `json:"id"`
	PipelineID  int64  `json:pipeline_id`
	Name        string `json:"name"`
	Description string `json:"description"`
	UiDisplay   int    `json:"ui_display"`
	Type        string `json:"type"`
	Config      Config `json:config`
	Status      string `json:"status"`
	Error       string `json:error`
	CreatedAt   string `json:"create_at"`
	UpdatedAt   string `json:"update_at"`
}

type Config struct {
	SourcePath string `json: source_path`
	TargetPath string `json: target_path`
}
