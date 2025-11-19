//go:build !wasm
// +build !wasm

package formPlatform

import (
	. "github.com/cdvelop/tinystring"
)

const (
	dbFieldTypeInt    dbFieldType = "INT"
	dbFieldTypeString dbFieldType = "VARCHAR(255)"
)

type dbFieldType string

func (t entity) CreateTableSQL() string {
	var sb = Convert()
	sb.Write(Fmt("CREATE TABLE IF NOT EXISTS %s (\n", t.TableName))

	for i, column := range t.Fields {
		sb.Write(Fmt("    %s %s", column.Name, column.DbType))

		if column.Unique {
			sb.Write(" UNIQUE")
		}

		if column.PrimaryKey {
			sb.Write(" PRIMARY KEY")
			if column.DbType == dbFieldTypeInt {
				// sb.Write(" AUTO_INCREMENT")
			}
		}

		if column.NotNull {
			sb.Write(" NOT NULL")
		}

		if column.ForeignKey != nil {
			// CONSTRAINT fk_departments FOREIGN KEY (id_department) REFERENCES departments(id_department)
			sb.Write(Fmt(",\n CONSTRAINT fk_%v FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE CASCADE",
				column.ForeignKey.TableName, column.Name, column.ForeignKey.TableName, column.Name))
		}

		if i < len(t.Fields)-1 {
			sb.Write(",")
		}
		sb.Write("\n")
	}

	sb.Write(");")
	return sb.String()
}
