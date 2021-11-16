package sqlstore

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"

	"github.com/mattermost/mattermost-server/v6/shared/mlog"

	sq "github.com/Masterminds/squirrel"
)

var (
	errUnsupportedOperation = errors.New("unsupported operation")
)

func (s *SQLStore) upsertTeamSignupToken(db sq.BaseRunner, team model.Team) error {
	now := utils.GetMillis()

	query := s.getQueryBuilder(db).
		Insert(s.tablePrefix+"teams").
		Columns(
			"id",
			"signup_token",
			"modified_by",
			"update_at",
		).
		Values(
			team.ID,
			team.SignupToken,
			team.ModifiedBy,
			now,
		)
	if s.dbType == mysqlDBType {
		query = query.Suffix("ON DUPLICATE KEY UPDATE signup_token = ?, modified_by = ?, update_at = ?",
			team.SignupToken, team.ModifiedBy, now)
	} else {
		query = query.Suffix(
			`ON CONFLICT (id)
			 DO UPDATE SET signup_token = EXCLUDED.signup_token, modified_by = EXCLUDED.modified_by, update_at = EXCLUDED.update_at`,
		)
	}

	_, err := query.Exec()
	return err
}

func (s *SQLStore) upsertTeamSettings(db sq.BaseRunner, team model.Team) error {
	now := utils.GetMillis()
	signupToken := utils.NewID(utils.IDTypeToken)

	settingsJSON, err := json.Marshal(team.Settings)
	if err != nil {
		return err
	}

	query := s.getQueryBuilder(db).
		Insert(s.tablePrefix+"teams").
		Columns(
			"id",
			"signup_token",
			"settings",
			"modified_by",
			"update_at",
		).
		Values(
			team.ID,
			signupToken,
			settingsJSON,
			team.ModifiedBy,
			now,
		)
	if s.dbType == mysqlDBType {
		query = query.Suffix("ON DUPLICATE KEY UPDATE settings = ?, modified_by = ?, update_at = ?", settingsJSON, team.ModifiedBy, now)
	} else {
		query = query.Suffix(
			`ON CONFLICT (id)
			 DO UPDATE SET settings = EXCLUDED.settings, modified_by = EXCLUDED.modified_by, update_at = EXCLUDED.update_at`,
		)
	}

	_, err = query.Exec()
	return err
}

func (s *SQLStore) getTeam(db sq.BaseRunner, id string) (*model.Team, error) {
	var settingsJSON string

	query := s.getQueryBuilder(db).
		Select(
			"id",
			"signup_token",
			"COALESCE(settings, '{}')",
			"modified_by",
			"update_at",
		).
		From(s.tablePrefix + "teams").
		Where(sq.Eq{"id": id})
	row := query.QueryRow()
	team := model.Team{}

	err := row.Scan(
		&team.ID,
		&team.SignupToken,
		&settingsJSON,
		&team.ModifiedBy,
		&team.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(settingsJSON), &team.Settings)
	if err != nil {
		s.logger.Error(`ERROR GetTeam settings json.Unmarshal`, mlog.Err(err))
		return nil, err
	}

	return &team, nil
}

func (s *SQLStore) hasTeamAccess(db sq.BaseRunner, userID string, teamID string) (bool, error) {
	return true, nil
}

func (s *SQLStore) getTeamCount(db sq.BaseRunner) (int64, error) {
	query := s.getQueryBuilder(db).
		Select(
			"COUNT(*) AS count",
		).
		From(s.tablePrefix + "teams")

	rows, err := query.Query()
	if err != nil {
		s.logger.Error("ERROR GetTeamCount", mlog.Err(err))
		return 0, err
	}
	defer s.CloseRows(rows)

	var count int64

	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		s.logger.Error("Failed to fetch team count", mlog.Err(err))
		return 0, err
	}
	return count, nil
}

// ToDo: maybe remove
func (s *SQLStore) getUserTeams(_ sq.BaseRunner, _ string) ([]model.UserTeam, error) {
	return nil, fmt.Errorf("GetUserTeams %w", errUnsupportedOperation)
}
