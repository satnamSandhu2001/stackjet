package models

type Stack struct {
	ID        int64  `json:"id"`
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Type      string `json:"type"`
	RepoUrl   string `json:"repo_url"`
	Branch    string `json:"branch"`
	Remote    string `json:"remote"`
	Port      uint16 `json:"port"` // 1024–65535
	CreatedAt string `json:"created_at"`
}
