package usecase

import (
	"github.com/mike-neck/go-hakone-qualification/hakone"
	"sort"
)

type Top10RecordsRepository interface {
	FindTop10FinishTimeRecordsByTeamName(name hakone.TeamName) []hakone.Record
}

type RecordService struct {
	TeamRepository  TeamRepository
	Top10Repository Top10RecordsRepository
}

type PersonalRecord struct {
	RankAmongAll  int
	RankAmongTeam int
	Grade         hakone.Grade
	Time          hakone.Time
}

type TeamRecords struct {
	Team    hakone.Team
	Records []PersonalRecord
}

type Top10RecordsByTeam struct {
	Records []TeamRecords
}

func (rs *RecordService) FindTop10RecordsByNames(names []hakone.TeamName) Top10RecordsByTeam {
	teamSize := len(names)
	if teamSize == 0 {
		return Top10RecordsByTeam{}
	}
	teamRecords := make([]TeamRecords, 0)
	for _, name := range names {
		team, err := rs.TeamRepository.FindTeamByName(string(name))
		if err != nil {
			continue
		}
		personalRecords := make([]PersonalRecord, 10)
		records := rs.Top10Repository.FindTop10FinishTimeRecordsByTeamName(name)
		sort.Slice(records, func(i, j int) bool {
			return records[i].Order < records[j].Order
		})
		for index, record := range records {
			personalRecords[index] = PersonalRecord{
				RankAmongAll:  record.Order,
				RankAmongTeam: index + 1,
				Grade:         record.Grade,
				Time:          record.FinishTime,
			}
		}
		teamRecords = append(teamRecords, TeamRecords{
			Team:    *team,
			Records: personalRecords,
		})
	}
	return Top10RecordsByTeam{Records: teamRecords}
}
