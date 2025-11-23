package server

import "context"

// save new url to server

// saveShortenURL(url, shorten) takes a created url and its shorten version and save it to database.
// If successfull it will return ResponseCreatedURLData type to return to user, if not err != nil
func (s *shortenServer) saveShortenURL(url, shorten string) (*ResponseCreatedURLData, error) {
	// Save to database and get response info
	var res ResponseCreatedURLData
	err := s.DB.QueryRow(
		context.Background(), `
		INSERT INTO shortened (url, shortened_url)
		VALUES ($1, $2)
		RETURNING id, url, shortened_url, created_at, updated_at, accessed`,
		url, shorten,
	).Scan(&res.ID, &res.URL, &res.ShortCode, &res.CreatedAt, &res.UpdatedAt, &res.Accessed)
	if err != nil {
		return nil, err
	}
	//Scan(&res.ID, &res.URL, &res.ShortCode, &res.CreatedAt, &res.UpdatedAt)
	return &res, nil
}
