package domainmodel

import "time"

type AppselfAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AppselfUserInfo struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Role      UserRole   `json:"role"`
	Bio       string     `json:"bio"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type AppselfUsersListItem struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

type AppselfPagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type AppselfUsersList struct {
	Data       []AppselfUsersListItem `json:"data"`
	Pagination AppselfPagination      `json:"pagination"`
}

type AppselfNodeInfo struct {
	Id          int64      `json:"id"`
	Mid         string     `json:"mid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsValidated bool       `json:"is_validated"`
	ValidatedBy *int64     `json:"validated_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type AppselfNodesListItem struct {
	Id          int64  `json:"id"`
	Mid         string `json:"mid"`
	Name        string `json:"name"`
	IsValidated bool   `json:"is_validated"`
}

type AppselfNodesList struct {
	Data       []AppselfNodesListItem `json:"data"`
	Pagination AppselfPagination      `json:"pagination"`
}

type AppselfSessionInfo struct {
	Id        int64      `json:"id"`
	UserId    *int64     `json:"user_id"`
	NodeId    *int64     `json:"node_id"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
}

type AppselfSessionsList struct {
	Data       []AppselfSessionInfo `json:"data"`
	Pagination AppselfPagination    `json:"pagination"`
}

type AppselfPostSessionResponse struct {
	Id int64 `json:"id"`
}

type AppselfRecordItem struct {
	Id              int64     `json:"id"`
	SessionId       *int64    `json:"session_id"`
	Time            time.Time `json:"time"`
	HeartRate       *float64  `json:"heart_rate"`
	BodyTemperature *float64  `json:"body_temperature"`
	Latitude        *float64  `json:"latitude"`
	Longitude       *float64  `json:"longitude"`
}

type AppselfRecordsList struct {
	Data []AppselfRecordItem `json:"data"`
}
