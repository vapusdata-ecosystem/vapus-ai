package mysql

import (
	"context"
)

func (m *MysqlStore) RunDDL(Ctx context.Context, query *string) error {
	// query the mysql database
	_, err := m.DB.Raw(*query).Rows()
	if err != nil {
		return err
	}
	return nil
}
