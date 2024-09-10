package models

import "time"

type Repository struct {
	ID             		string    `json:"-" gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:id" bson:"_id,omitempty"` // UUID for PostgreSQL, ObjectId for MongoDB
	RepoID         		int       `json:"id" gorm:"type:int;column:repo_id" bson:"repo_id"`
	Name           		string    `json:"name" gorm:"size:255;column:name" bson:"name"`
	FullName       		string    `json:"full_name" gorm:"size:255;column:full_name" bson:"full_name"`
	Description    		string    `json:"description" gorm:"type:text;column:description" bson:"description"`
	URL            		string    `json:"html_url" gorm:"size:255;column:url" bson:"url"`
	Language       		string    `json:"language" gorm:"size:255;column:language" bson:"language"`
	ForksCount     		int       `json:"forks_count" gorm:"type:int;column:forks_count" bson:"forks_count"`
	StarsCount     		int       `json:"stars_count" gorm:"type:int;column:stars_count" bson:"stars_count"`
	OpenIssuesCount		int      `json:"open_issues_count" gorm:"type:int;column:open_issues_count" bson:"open_issues_count"`
	WatchersCount  		int       `json:"watchers_count" gorm:"type:int;column:watchers_count" bson:"watchers_count"`
	CreatedAt      		time.Time `json:"created_at" gorm:"column:created_at" bson:"created_at"`
	UpdatedAt      		time.Time `json:"updated_at" gorm:"column:updated_at" bson:"updated_at"`
}