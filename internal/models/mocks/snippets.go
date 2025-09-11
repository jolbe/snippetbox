package mocks

import (
	"time"

	"snippetbox.gregor-pifko/internal/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "Mock Snippet",
	Content: "This is a mock snippet",
	Created: time.Now(),
	Expires: time.Now().Add(5 * 24 * time.Hour),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
