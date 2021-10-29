package schema

func getFKActionDDL(action SQLForeignKeyAction) string {
	switch action {
	case ActionCascade:
		return "CASCADE"
	case ActionNoAction:
		return "NO ACTION"
	case ActionRestrict:
		return "RESTRICT"
	case ActionSetDefault:
		return "SET DEFAULT"
	case ActionSetNull:
		return "SET NULL"
	default:
		return ""
	}
}
