package usecase

import (
	"github.com/mike-neck/go-hakone-qualification/hakone"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ListTeamsTestRepository struct {
	teams []hakone.Team
}

func (r *ListTeamsTestRepository) ListAllTeams() []hakone.Team {
	teams := r.teams
	result := make([]hakone.Team, len(teams))
	for i, t := range teams {
		result[i] = t
	}
	return result
}

var listTeamsTestRepository = &ListTeamsTestRepository{
	teams: []hakone.Team{
		{Id: 3, Name: "早稲田大"},
		{Id: 1, Name: "東海大"},
		{Id: 2, Name: "東洋大"},
		{Id: 4, Name: "日本大"},
		{Id: 5, Name: "日本体育大"},
	},
}

func TestTeamService_ListAllTeamsOrderById(t *testing.T) {
	service := TeamService{Repository: listTeamsTestRepository}

	teams := service.ListAllTeamsOrderById()

	assert.Equal(t, 1, teams[0].Id)
	assert.Equal(t, 2, teams[1].Id)
	assert.Equal(t, 3, teams[2].Id)
	assert.Equal(t, 3, listTeamsTestRepository.teams[0].Id)
	assert.Equal(t, 1, listTeamsTestRepository.teams[1].Id)
	assert.Equal(t, 2, listTeamsTestRepository.teams[2].Id)
}

func TestTeamService_FindTeamByNameStartsWith(t *testing.T) {
	service := TeamService{Repository: listTeamsTestRepository}

	teams := service.FindTeamsByNameStartsWith("東")

	assert.Equal(t, 2, len(teams))
	assert.Equal(t, "東海大", teams[0].Name)
	assert.Equal(t, "東洋大", teams[1].Name)
}
