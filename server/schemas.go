package server

// Requests body format shema

type CreateURLData struct {
	Url string `json:"url"`
}

// Response body format schemas

type ResponseCreatedURLData struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	ShortCode string `json:"shortCode"`
	CreatedAt int    `json:"createdAt"`
	UpdatedAt int    `json:"updatedAt"`
}
