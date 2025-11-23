package server

import "context"

// save new url to server

// saveShortenURL(url, shorten) takes a created url and its shorten version and save it to database.
// If successfull it will return ResponseCreatedURLData type to return to user, if not err != nil
func (s *shortenServer) saveShortenURL(url, shorten string) (*ResponseCreatedURLData, error) {
	// save to database
	_, err := s.DB.Exec(
		context.Background(), `
		INSERT INTO shortened (url, shortened_url)
		VALUES ($1, $2);`,
		url, shorten,
	)
	if err != nil {
		return nil, nil
	}

	// Get response format
	var res ResponseCreatedURLData

	row := s.DB.QueryRow(context.Background(), `SELECT * FROM shortened WHERE shortened_url=$1`, shorten)
	row.Scan(&res.ID, &res.URL, &res.ShortCode, &res.CreatedAt, &res.UpdatedAt, &res.Accessed)
	//Scan(&res.ID, &res.URL, &res.ShortCode, &res.CreatedAt, &res.UpdatedAt)
	return &res, nil
}
