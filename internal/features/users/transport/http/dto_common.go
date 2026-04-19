package users_transport_http

import "github.com/artlink52/go-todoapp/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	userDTOs := make([]UserDTOResponse, len(users))

	for i, user := range users {
		userDTOs[i] = userDTOFromDomain(user)
	}

	return userDTOs
}
