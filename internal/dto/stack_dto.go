package dto

type Stack_CreateRequest struct {
	ID      int64  `json:"id"`
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	RepoUrl string `json:"repo_url"`
	Branch  string `json:"branch"`
	Remote  string `json:"remote"`
	Port    uint16 `json:"port" binding:"required"`
}
