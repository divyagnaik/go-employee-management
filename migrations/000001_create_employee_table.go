package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createTable = `CREATE TABLE IF NOT EXISTS employee
(
    id             int         not null primary key,
    name           varchar(50) not null,
    position         varchar(50)  not null,
    salary varchar(10) not null
);`

func createTableEmployee() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTable)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
