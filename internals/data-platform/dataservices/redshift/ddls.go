package redshift

import (
	"context"
)

func (m *RedshiftStore) RunDDL(ctx context.Context, query *string) error {
	// query the mysql database
	if m.Orm != nil {
		_, err := m.Orm.Raw(*query).Rows()
		if err != nil {
			return err
		}
		return nil
	}
	m.logger.Debug().Msgf("Running DDL query: %s", *query)
	_, err := m.Conn.Exec(*query)
	return err
}
