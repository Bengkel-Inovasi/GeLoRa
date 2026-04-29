package adaptersoutboundrepositoryalert

import "github.com/Masterminds/squirrel"

func (s *sqliteImpl) queryCreateAlert(userId int64, message string) (query string, args []any, err error) {
	return s.sqrQuestion.
		Insert("gelora_alerts").
		Columns("user_id", "message").
		Values(userId, message).
		ToSql()
}

func (s *sqliteImpl) queryReadAlerts(unacknowledgedOnly bool) (query string, args []any, err error) {
	q := s.sqrQuestion.
		Select("id", "user_id", "message", "acknowledged_at", "created_at").
		From("gelora_alerts").
		OrderBy("created_at DESC")

	if unacknowledgedOnly {
		q = q.Where(squirrel.Eq{"acknowledged_at": nil})
	}

	return q.ToSql()
}

func (s *sqliteImpl) queryUpdateAlertAcknowledge(id int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Update("gelora_alerts").
		Set("acknowledged_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
}
