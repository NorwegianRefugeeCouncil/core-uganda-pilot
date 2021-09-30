package seeder

var (
	// Individuals - UG Beneficiaries
	JohnDoe     = ugandaIndividual(individual("c529d679-3bb6-4a20-8f06-c096f4d9adc1", "John Sinclair Doe", "John Doe", "1983-04-23", "john.doe", "Refugee", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "No", "Yes", "Moderate", "No", "", "No", "", "Kenya", "Kiswahili, English", "English", "123 Main Street, Kampala", "0123456789", "", "Email", "No"), "1983-04-23", "0", "0", "0", "0", "0", "0")
	MaryPoppins = ugandaIndividual(individual("bbf539fd-ebaa-4438-ae4f-8aca8b327f42", "Mary Poppins", "Mary Poppins", "1983-04-23", "mary.poppins", "Internally Displaced Person", "Female", "Yes", "https://link-to-consent.proof", "No", "No", "No", "No", "", "No", "", "No", "", "Uganda", "Rukiga, English", "Rukiga", "901 First Avenue, Kampala", "0123456789", "", "Telegram", "Yes"), "1983-04-23", "0", "0", "0", "0", "0", "0")
	BoDiddley   = ugandaIndividual(individual("26335292-c839-48b6-8ad5-81271ee51e7b", "Ellas McDaniel", "Bo Diddley", "1983-04-23", "bo.diddley", "Host Community", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "Yes", "No", "", "No", "", "No", "", "Somalia", "Somali, Arabic, English", "English", "101 Main Street, Kampala", "0123456789", "", "Whatsapp", "No"), "1983-04-23", "0", "0", "0", "0", "0", "0")

	// Beneficiaries
	_ = beneficiary(JohnDoe)
	_ = beneficiary(MaryPoppins)
	_ = beneficiary(BoDiddley)
)
