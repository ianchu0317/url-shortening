package server

import (
	"context"
)

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

// retrieveOriginalURL(shortCode) returns the original URL given a shortCode. IF not found then will return ""
func (s *shortenServer) retrieveOriginalURL(shortCode string) (*ResponseCreatedURLData, error) {
	// Check short code on server
	var res ResponseCreatedURLData
	err := s.DB.QueryRow(
		context.Background(),
		`UPDATE shortened SET accessed=(accessed+1) WHERE short_code=$1
		RETURNING id, url, short_code, created_at, updated_at, accessed`,
		shortCode,
	).Scan(&res.ID, &res.URL, &res.ShortCode, &res.CreatedAt, &res.UpdatedAt, &res.Accessed)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// updateAccessCount updates de access counter given a shortCode
func (s *shortenServer) updateAccessCount(shortCode string) error {
	_, err := s.DB.Exec(
		context.Background(),
		`UPDATE shortened SET accessed=(accessed+1) WHERE short_code=$1`,
		shortCode,
	)
	if err != nil {
		return err
	}
	return nil
}

// updateOriginalURL takes new url and a shortCode. After created will return the updated data.
func (s *shortenServer) updateOriginalURL(newURL, shortCode string) (*ResponseCreatedURLData, error) {
	var res ResponseCreatedURLData
	err := s.DB.QueryRow(
		context.Background(),
		`UPDATE shortened
		SET url=$1, updated_at=NOW()
		WHERE short_code=$2
		RETURNING id, url, short_code, created_at, updated_at, accessed`,
		newURL,
		shortCode,
	).Scan(&res.ID, &res.URL, &res.ShortCode, &res.CreatedAt, &res.UpdatedAt, &res.Accessed)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *shortenServer) deleteShortCode(shortCode string) error {
	_, err := s.DB.Exec(
		context.Background(),
		`DELETE FROM shortened 
		WHERE short_code=$1`,
		shortCode,
	)
	if err != nil {
		return err
	}
	return nil
}

// Auxiliar functions check if

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

// isShortCodeInDB(shortCode) takes a short code and check if it is already in DB
// If is in server it returns true, otherwise false
func (s *shortenServer) isShortCodeInDB(shortCode string) (bool, error) {
	commandTag, err := s.DB.Exec(
		context.Background(),
		`SELECT * FROM shortened WHERE short_code=$1;`,
		shortCode,
	)
	if err != nil {
		return false, err
	}
	return commandTag.RowsAffected() >= 1, nil
}
