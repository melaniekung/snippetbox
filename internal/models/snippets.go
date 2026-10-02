package models

import (
	"database/sql"
	"errors"
	"time"
)

// define type to hold data for individual snippet
type Snippet struct {
	ID	int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// define type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// insert new snippet into database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// get ID of new record in snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// return specific snippet based on id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// inititalize new zeroed Snippet struct
	var s Snippet

	// ccopy values from each field in sql.Row to corresponding field in the Snippet struct
	// NOTE: row.Scan arguments are pointers
	// 		 number of rows must match number of columns returned by statement
	err:= m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		switch {
		// query returns no rows -> error
		case errors.Is(err, sql.ErrNoRows):
			return Snippet{}, ErrNoRecord
		default:
			return Snippet{}, err
		}
	}

	return s, nil
}

// return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// IMPORTANT: defer rows.Close() to ensure sql.Rows resultset is properly closed before return
	// NOTE: should come after error check from Query()
	defer rows.Close()

	// initialize empty slice
	var snippets []Snippet

	// iterate through rows in resultset
	// resultset automatically closes itself once iteration is complete
	for rows.Next() {
		// create zero value struct
		var s Snippet
		// copy values from each field in row to new Snippet struct
		// NOTE: arguments must be pointers
		//       number of arguments must match number of columns returned by statement
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// append to slice of snippets
		snippets = append(snippets, s)
	}

	// retreieve any errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}