package server

import "context"

// save new url to server

// saveShortenURL(url, shorten) takes a created url and its shorten version and save it to database.
// If successfull it will return ResponseCreatedURLData type to return to user, if not err != nil
func (s *shortenServer) saveShortenURL(url, shortCode string) (*ResponseCreatedURLData, error) {
	// Save to database and get response info
	var res ResponseCreatedURLData
	err := s.DB.QueryRow(
		context.Background(), `
		INSERT INTO shortened (url, short_code)
		VALUES ($1, $2)
		RETURNING id, url, short_code, created_at, updated_at, accessed`,
		url, shortCode,
	).Scan(&res.ID, &res.URL, &res.ShortCode, &res.CreatedAt, &res.UpdatedAt, &res.Accessed)
	if err != nil {
		return nil, err
	}
	//Scan(&res.ID, &res.URL, &res.ShortCode, &res.CreatedAt, &res.UpdatedAt)
	return &res, nil
}

// isURLInDB(url) takes an URL and check if it is already in server.
// If is in server it returns true, otherwise false
func (s *shortenServer) isUrlInDB(url string) (bool, error) {
	commandTag, err := s.DB.Exec(
		context.Background(),
		`SELECT * FROM shortened WHERE url=$1;`,
		url,
	)
	if err != nil {
		return false, err
	}
	return commandTag.RowsAffected() >= 1, nil
}
