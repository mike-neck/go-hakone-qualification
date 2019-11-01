package usecase

import (
	"github.com/mike-neck/go-hakone-qualification/hakone"
	"sort"
	"strings"
)

type TeamRepository interface {
	ListAllTeams() []hakone.Team
}

type TeamService struct {
	Repository TeamRepository
}

func (ts *TeamService) ListAllTeams() []hakone.Team {
	return ts.Repository.ListAllTeams()
}

func (ts *TeamService) ListAllTeamsOrderById() []hakone.Team {
	teams := ts.ListAllTeams()
	sort.SliceStable(teams, func(i, j int) bool {
		return teams[i].Id < teams[j].Id
	})
	return teams
}

func (ts *TeamService) FindTeamsByNameStartsWith(prefix string) []hakone.Team {
	teams := make([]hakone.Team, 0)
	allTeams := ts.Repository.ListAllTeams()
	for _, team := range allTeams {
		if strings.HasPrefix(team.Name, prefix) {
			teams = append(teams, team)
		}
	}
	return teams
}
