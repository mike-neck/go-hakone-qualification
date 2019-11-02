package usecase

import (
	"fmt"
	"github.com/mike-neck/go-hakone-qualification/hakone"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Top10RecordsRepoTestImpl struct {
	Records []hakone.Record
}

func newRecord(order int, name, team string, grade, hour, sec int) hakone.Record {
	return hakone.Record{
		Order:      order,
		Runner:     hakone.Runner(name),
		Grade:      hakone.Grade(fmt.Sprintf("(%d)", grade)),
		Team:       hakone.TeamName(team),
		FinishTime: hakone.Time(hour*60*60 + sec),
	}
}

func (tr *Top10RecordsRepoTestImpl) FindTop10FinishTimeRecordsByTeamName(name hakone.TeamName) []hakone.Record {
	result := make([]hakone.Record, 10)
	records := tr.Records
	max := len(records)
	index := 0
	for i := 0; i < max && index < 10; i++ {
		rec := records[i]
		if rec.Team == name {
			result[index] = rec
			index++
		}
	}
	return result
}

func TestRecordService_FindTop10RecordsByNames(t *testing.T) {
	records := makeRecords()
	service := RecordService{
		TeamRepository:  listTeamsTestRepository,
		Top10Repository: &Top10RecordsRepoTestImpl{Records: records},
	}

	result := service.FindTop10RecordsByNames([]hakone.TeamName{"東洋大", "日本体育大"})

	recs := result.Records
	assert.Equal(t, 2, len(recs))
	if len(recs) != 2 {
		return
	}
	team1 := recs[0]
	assert.Equal(t, 10, len(team1.Records))
	assert.Equal(t, "東洋大", team1.Team.Name)
	assert.Equal(t, 1, team1.Records[0].RankAmongAll)
	assert.Equal(t, 16, team1.Records[1].RankAmongAll)
	assert.Equal(t, 136, team1.Records[9].RankAmongAll)
	team2 := recs[1]
	assert.Equal(t, 10, len(team2.Records))
	assert.Equal(t, "日本体育大", team2.Team.Name)
	assert.Equal(t, 3, team2.Records[0].RankAmongAll)
	assert.Equal(t, 21, team2.Records[1].RankAmongAll)
	assert.Equal(t, 165, team2.Records[9].RankAmongAll)
}

func makeRecords() []hakone.Record {
	records := make([]hakone.Record, 33)
	for i := 0; i < 11; i++ {
		// order 1,16,31,46,61,76,91,106,121,136,151
		records[0+i*3] =
			newRecord(1+i*15, fmt.Sprintf("東洋大ランナー-%d", i), "東洋大", (i%4)+1, 1, 61+i*15)
		// order 2,14,26,38,50,62,74,86,98,110,122
		records[1+i*3] =
			newRecord(2+i*12, fmt.Sprintf("東海大ランナー-%d", i), "東海大", (i%4)+1, 1, 62+i*12)
		// order 3,21,39,57,75,93,111,129,147,165,183
		records[2+i*3] =
			newRecord(3+i*18, fmt.Sprintf("日体大ランナー-%d", i), "日本体育大", (i%4)+1, 1, 63+i*18)
	}
	return records
}

func TestRecordService_FindTop10RecordsByNames_WithUnknownTeam(t *testing.T) {
	records := makeRecords()
	service := RecordService{
		TeamRepository:  listTeamsTestRepository,
		Top10Repository: &Top10RecordsRepoTestImpl{Records: records},
	}

	result := service.FindTop10RecordsByNames([]hakone.TeamName{"筑波大", "鹿屋体育大"})

	recs := result.Records
	assert.Equal(t, 0, len(recs))
}
