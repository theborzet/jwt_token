package pagination

import (
	"errors"
	"fmt"

	"github.com/theborzet/time-tracker/internal/models"
)

const DefaultPageSize = 10

func PaginateUser(users []*models.User, page int) ([]*models.User, Paginator, error) {
	var paginator Paginator

	totalItems := len(users)
	if totalItems == 0 {
		return nil, Paginator{}, errors.New("no items")
	}
	paginator.TotalItems = totalItems
	paginator.PageSize = DefaultPageSize
	paginator.TotalPages = (totalItems + DefaultPageSize - 1) / DefaultPageSize

	if page > paginator.TotalPages {
		return nil, paginator, errors.New("requested page exceeds total pages")
	}

	if page < 1 {
		page = 1
	} else if page > paginator.TotalPages {
		page = paginator.TotalPages
	}

	paginator.CurrentPage = page

	start := (page - 1) * DefaultPageSize
	end := start + DefaultPageSize
	if end > totalItems {
		end = totalItems
	}

	paginatedBooks := users[start:end]

	paginator.HasPrevious = page > 1
	paginator.HasNext = page < paginator.TotalPages

	if paginator.HasPrevious {
		paginator.PreviousPage = page - 1
	}

	if paginator.HasNext {
		paginator.NextPage = page + 1
	}

	paginator.PageNumbers = make([]PageNumber, paginator.TotalPages)
	for i := range paginator.PageNumbers {
		paginator.PageNumbers[i] = PageNumber{
			Number:   i + 1,
			URL:      fmt.Sprintf("/users?page=%d", i+1),
			IsActive: i+1 == page,
		}
	}

	return paginatedBooks, paginator, nil
}
