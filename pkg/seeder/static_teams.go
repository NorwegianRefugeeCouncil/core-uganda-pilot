package seeder

import (
	"github.com/nrc-no/core/pkg/iam"
)

var (
	teams                   []iam.Team
	KampalaRegistrationTeam = team("5ccc5a95-cc90-4740-a037-3a57557745eb", "Kampala Registration Team")
	KampalaProtectionTeam   = team("ac9b8d7d-d04d-4850-9a7f-3f93324c0d1e", "Kampala Protection Team")
	KampalaICLATeam         = team("a43f84d5-3f8a-48c4-a896-5fb0fcd3e42b", "Kampala ICLA Team")
	ColombiaTeam            = team("a6bc6436-fcea-4738-bde8-593e6480e1ad", "Colombia Team")

	// Nationalities
	_ = nationality("3ea0eea3-fef8-4bef-983a-b52a7814efbd", KampalaRegistrationTeam, ugandaCountry)
	_ = nationality("b58e4d26-fe8e-4442-8449-7ec4ca3d9066", KampalaProtectionTeam, ugandaCountry)
	_ = nationality("23e3eb5e-592e-42e2-8bbf-ee097d93034c", KampalaICLATeam, ugandaCountry)
	_ = nationality("7ba6d2ee-1af9-447c-8000-7719467b3414", ColombiaTeam, colombiaCountry)
)
