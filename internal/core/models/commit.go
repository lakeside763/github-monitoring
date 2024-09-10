package models

import "time"

type Commit struct {
	ID      		string      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" bson:"_id,omitempty"` // UUID for PostgreSQL, ObjectID for MongoDB
	RepoId	int							`json:"repo_id" gorm:"type:int;column:repo_id"`
	SHA     		string      `gorm:"size:40;column:sha" bson:"sha"`
	Message 		string      `gorm:"type:text;column:message" bson:"message"`
	Author  		string      `gorm:"size:255;column:author" bson:"author"`
	Date    		time.Time   `gorm:"column:date" bson:"date"`
	URL     		string      `gorm:"size:255;column:url" bson:"url"`
	RepoID  		string      `gorm:"type:uuid;column:repo_id" bson:"repoId,omitempty"`
	CreatedAt   time.Time 	`json:"created_at" gorm:"column:created_at" bson:"created_at"`
	UpdatedAt   time.Time 	`json:"updated_at" gorm:"column:updated_at" bson:"updated_at"`
}

type ResponseCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	HTMLURL string `json:"html_url"`
}

type FetchCommitInput struct {
	Repo				    string
	Since           time.Time
	Until           time.Time
	Limit           int
	Page            int
	NextWindowSince time.Time
}

type Pagination struct {
	Since           time.Time
	Until           time.Time
	Page            int
	NextWindowSince time.Time
}

type CommitParams struct {
	Since  string `schema:"since"`
	Until  string `schema:"until"`
	Limit  int    `schema:"limit"`
	Page   int    `schema:"page"`
}