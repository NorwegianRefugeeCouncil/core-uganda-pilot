package devinit

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func (c *Config) makePostgresInit() error {
	sb := &strings.Builder{}
	sb.WriteString("#!/bin/bash\n")
	sb.WriteString("\n")
	sb.WriteString("set -e\n")
	sb.WriteString("set -u\n")
	sb.WriteString("\n")
	sb.WriteString("function create_user_and_database() {\n")
	sb.WriteString("  local database=$1\n")
	sb.WriteString("  local user=$2\n")
	sb.WriteString("  local password=$3\n")
	sb.WriteString("  echo \">> Creating user '$user' and database '$database'\"\n")
	sb.WriteString("  psql -v ON_ERROR_STOP=1 --username \"$POSTGRES_USER\" <<-EOSQL\n")
	sb.WriteString("    CREATE USER $user WITH PASSWORD '$password';\n")
	sb.WriteString("    CREATE DATABASE $database;\n")
	sb.WriteString("    GRANT ALL PRIVILEGES ON DATABASE $database TO $user;\n")
	sb.WriteString("EOSQL\n")
	sb.WriteString("}\n")
	sb.WriteString("\n")
	for _, user := range c.dbUsers {
		sb.WriteString(fmt.Sprintf("create_user_and_database \"%s\" \"%s\" \"%s\"\n",
			user.database,
			user.username,
			user.password))
	}

	if err := os.WriteFile(path.Join(PostgresDir, "init.sh"), []byte(sb.String()), os.ModePerm); err != nil {
		return err
	}

	return nil
}
