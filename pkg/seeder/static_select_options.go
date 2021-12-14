package seeder

import "github.com/nrc-no/core/pkg/api/types"

// General options
var (
	yesNoChoice = []*types.SelectOption{
		{Name: "Yes"},
		{Name: "No"},
	}
	wgShortSet = []*types.SelectOption{
		{Name: "Moderate Impairment"},
		{Name: "Severe Impairment"},
	}
	globalDisplacementStatuses = []*types.SelectOption{
		{Name: "Refugee"},
		{Name: "Internally Displaced Person (DP)"},
		{Name: "Host Community"},
		{Name: "Other"},
	}
	globalGenders = []*types.SelectOption{
		{Name: "Male"},
		{Name: "Female"},
		{Name: "Non-Binary"},
		{Name: "Other"},
	}
) // Uganda Specific Options
var (
	ugCC = []*types.SelectOption{
		{Name: "ICLA"},
		{Name: "Education"},
		{Name: "LFS"},
		{Name: "Protection"},
		{Name: "WASH"},
		{Name: "Shelter"},
	}
	ugNationality = []*types.SelectOption{
		{Name: "Uganda"},
		{Name: "South Sudan"},
		{Name: "Eritrea"},
		{Name: "Somalia"},
		{Name: "Democratic Republic of Congo"},
		{Name: "Rwanda"},
		{Name: "Burundi"},
		{Name: "Sudan"},
	}
	ugRelationshipToHH = []*types.SelectOption{
		// TBD
	}
	ugTypeOfIdentification = []*types.SelectOption{
		{Name: "Family Attestation (OPM)"},
		{Name: "Refugee ID (OPM)"},
		{Name: "Asylum Certificate (OPM)"},
		{Name: "Resident ID"},
		{Name: "National ID"},
		{Name: "Passport"},
		{Name: "Conventional Travel Document"},
	}
	ugLocationOfIdentification = []*types.SelectOption{
		{Name: "Kabusu Access Center"},
		{Name: "Nsambya Access Center"},
		{Name: "Kisenyi ICLA Center"},
		{Name: "Lukuli ICLA Center"},
		{Name: "Kawempe ICLA Center"},
		{Name: "Ndejje ICLA Center"},
		{Name: "Mengo Field Office"},
		{Name: "Community (Specify location)"},
		{Name: "Home Visit"},
		{Name: "Phone"},
		{Name: "Other (Specify)"},
	}
	ugSourceOfIdentification = []*types.SelectOption{
		{Name: "Walk-in Center"},
		{Name: "FFRM Referral"},
		{Name: "Internal Referral (Other - Specify)"},
		{Name: "ICLA Outreach Team"},
		{Name: "External Referral (Community Leader/Contact"},
		{Name: "External Referral (INGO/LNGO)"},
		{Name: "External Referral (UN Agency)"},
		{Name: "External Referral (Government)"},
		{Name: "External Referral (Other - Specify)"},
		{Name: "Self (Telephone)"},
		{Name: "Self (Email)"},
		{Name: "Internal Referral (Other NRC Sector – Specify)"},
		{Name: "CBP Outreach Team "},
		{Name: "Other NRC Outreach Team (Specify)"},
		{Name: "Other - Specify"},
	}
	ugAdmin2 = []*types.SelectOption{
		{Name: "Abim"},
		{Name: "Adjumani"},
		{Name: "Alebtong"},
		{Name: "Amolatar"},
		{Name: "Amudat"},
		{Name: "Amuria"},
		{Name: "Amuru"},
		{Name: "Apac"},
		{Name: "Budaka"},
		{Name: "Bugiri"},
		{Name: "Buikwe"},
		{Name: "Bukomansimbi"},
		{Name: "Bukwo"},
		{Name: "Bulambuli"},
		{Name: "Buliisa"},
		{Name: "Bundibugyo"},
		{Name: "Bushenyi"},
		{Name: "Buyende"},
		{Name: "Dokolo"},
		{Name: "Butambala"},
		{Name: "Hoima"},
		{Name: "Iganga"},
		{Name: "Kaabong"},
		{Name: "Kabale"},
		{Name: "Kabarole"},
		{Name: "Kalangala"},
		{Name: "Kaliro"},
		{Name: "Kalungu"},
		{Name: "Kamuli"},
		{Name: "Kanungu"},
		{Name: "Kapchorwa"},
		{Name: "Katakwi"},
		{Name: "Kayunga"},
		{Name: "Sheema"},
		{Name: "Kitgum"},
		{Name: "Koboko"},
		{Name: "Kole"},
		{Name: "Kotido"},
		{Name: "Kisoro"},
		{Name: "Kween"},
		{Name: "Lamwo"},
		{Name: "Lira"},
		{Name: "Luuka"},
		{Name: "Lyantonde"},
		{Name: "Manafwa"},
		{Name: "Masaka"},
		{Name: "Masindi"},
		{Name: "Mayuge"},
		{Name: "Mbale"},
		{Name: "Mbarara"},
		{Name: "Moroto"},
		{Name: "Moyo"},
		{Name: "Nakapiripirit"},
		{Name: "Nakaseke"},
		{Name: "Nakasongola"},
		{Name: "Namutumba"},
		{Name: "Napak"},
		{Name: "Nebbi"},
		{Name: "Ngora"},
		{Name: "Buhweju"},
		{Name: "Ntoroko"},
		{Name: "Maracha"},
		{Name: "Otuke"},
		{Name: "Oyam"},
		{Name: "Pader"},
		{Name: "Rubirizi"},
		{Name: "Sironko"},
		{Name: "Soroti"},
		{Name: "Wakiso"},
		{Name: "Yumbe"},
		{Name: "Zombo"},
		{Name: "Isingiro"},
		{Name: "Mitooma"},
		{Name: "Kyegegwa"},
		{Name: "Ntungamo"},
		{Name: "Rukungiri"},
		{Name: "Kamwenge"},
		{Name: "Ibanda"},
		{Name: "Kasese"},
		{Name: "Kiruhura"},
		{Name: "Kyenjojo"},
		{Name: "Mubende"},
		{Name: "Gomba"},
		{Name: "Kiboga"},
		{Name: "Mpigi"},
		{Name: "Kyankwanzi"},
		{Name: "Kakumiro"},
		{Name: "Nwoya"},
		{Name: "Kiryandongo"},
		{Name: "Serere"},
		{Name: "Omoro"},
		{Name: "Arua"},
		{Name: "Lwengo"},
		{Name: "Sembabule"},
		{Name: "Rakai"},
		{Name: "Mityana"},
		{Name: "Luwero"},
		{Name: "Mukono"},
		{Name: "Kampala"},
		{Name: "Buvuma"},
		{Name: "Jinja"},
		{Name: "Namayingo"},
		{Name: "Busia"},
		{Name: "Bududa"},
		{Name: "Tororo"},
		{Name: "Butaleja"},
		{Name: "Bukedea"},
		{Name: "Kumi"},
		{Name: "Pallisa"},
		{Name: "Kibuku"},
		{Name: "Kaberamaido"},
		{Name: "Agago"},
		{Name: "Kagadi"},
		{Name: "Kibaale"},
		{Name: "Gulu"},
		{Name: "Rubanda"},
	}

	ugMeansOfContact = []*types.SelectOption{
		{Name: "Phone Call"},
		{Name: "Text Message"},
		{Name: "WhatsApp"},
		{Name: "Signal"},
		{Name: "Telegram"},
		{Name: "Email"},
		{Name: "Home Visit"},
		{Name: "Other - Specify"},
	}
	ugServicesRequested = []*types.SelectOption{
		{Name: "Health care (including medication)"},
		{Name: "Legal assistance"},
		{Name: "Education"},
		{Name: "Mental Health"},
		{Name: "Transportation"},
		{Name: "Food"},
		{Name: "Non-food items (including hygiene items)"},
		{Name: "Disability"},
		{Name: "MPC"},
		{Name: "Shelter/Housing"},
		{Name: "Shelter construction/repair"},
		{Name: "Youth Livelihoods (e.g. vocational training)"},
		{Name: "Small/Medium Business Grants"},
		{Name: "Other livelihood activities"},
	}
	ugPerceivedPriorityResponseLevel = []*types.SelectOption{
		{Name: "High - Immediate contact within 24 hrs"},
		{Name: "Medium - Contact within 72hrs"},
		{Name: "Low – Contact within five days "},
	}
	ugSpecificNeed = []*types.SelectOption{
		{Name: "Pregnant woman"},
		{Name: "Elderly taking care of minors alone"},
		{Name: "Single parent"},
		{Name: "Chronic illness"},
		{Name: "Legal protection needs"},
	}
	ugProtectionRisk = []*types.SelectOption{
		{Name: "Violence"},
		{Name: "Coercion"},
		{Name: "Discrimination"},
		{Name: "Deprivation"},
		{Name: "Other"},
	}
	ugProtectionConcern = []*types.SelectOption{
		{Name: "Physical violence"},
		{Name: "Neglect"},
		{Name: "Family separation"},
		{Name: "Arrest"},
		{Name: "Denial of resources"},
		{Name: "Psychosocial violence"},
	}
	ugProtectionVulnerability = []*types.SelectOption{
		{Name: "Child at risk"},
		{Name: "Elder at risk"},
		{Name: "Single parent"},
		{Name: "Separated child"},
		{Name: "Disability"},
		{Name: "Woman at risk"},
		{Name: "Legal & physical protection"},
		{Name: "Legal & physical protection"},
		{Name: "Medical condition"},
		{Name: "Pregnant/Lactating woman"},
	}
	ugProtectionSpecificNeeds = []*types.SelectOption{
		{Name: "Disability"},
		{Name: "Pregnant woman"},
		{Name: "Elder taking care of minors alone"},
		{Name: "Elder living alone"},
		{Name: "Single parent"},
		{Name: "Chronic illness"},
		{Name: "Legal protection needs"},
		{Name: "Child"},
		{Name: "Other"},
	}
	ugProtectionHomeSituation = []*types.SelectOption{
		{Name: "Lives alone"},
		{Name: "Lives with family"},
		{Name: "Hosted by relative"},
		{Name: "Other"},
	}
	ugProtectionDisability = []*types.SelectOption{
		{Name: "None"},
		{Name: "Moderate physical impairment"},
		{Name: "Severe physical impairment"},
		{Name: "Moderate sensory impairment"},
		{Name: "Severe sensory impairment"},
		{Name: "Moderate mental disability"},
		{Name: "Severe mental disability"},
	}
	ugProtectionHHAbilityToMeetNeeds = []*types.SelectOption{
		{Name: "We can meet all our needs without worry"},
		{Name: "We can meet our needs"},
		{Name: "We can barely meet our needs"},
		{Name: "We are unable to meet our needs"},
		{Name: "We are totally unable to meet our needs"},
	}
	ugProtectionFoodNeedsObstacles = []*types.SelectOption{
		{Name: "Insufficient funds"},
		{Name: "Distance issues"},
		{Name: "Insecurity"},
		{Name: "Social discrimination"},
		{Name: "Insufficient quantity of goods"},
		{Name: "Inadequate quality of goods"},
		{Name: "Insufficient service providers"},
		{Name: "Other"},
	}
	ugProtectionAccommodationNeedsObstacles = []*types.SelectOption{
		{Name: "Insufficient funds"},
		{Name: "Distance issues"},
		{Name: "Insecurity"},
		{Name: "Social discrimination"},
		{Name: "Insufficient quantity of goods"},
		{Name: "Inadequate quality of goods/services"},
		{Name: "Insufficient capabilities/competencies"},
		{Name: "Other"},
	}
	ugProtectionWASHNeedsObstacles = []*types.SelectOption{
		{Name: "Insufficient funds"},
		{Name: "Distance issues"},
		{Name: "Insecurity"},
		{Name: "Social discrimination"},
		{Name: "Insufficient quantity of goods/services"},
		{Name: "Inadequate quality of goods"},
		{Name: "Other"},
	}
	ugProtectionPriorityOpts = []*types.SelectOption{
		{Name: "High (follow-up requested in 24 hours)"},
		{Name: "Medium (follow-up in 3 days)"},
		{Name: "Low (follow-up in 7 days)"}}
	ugICLALegalIssues = []*types.SelectOption{
		{Name: "RSD"},
		{Name: "ELP"},
		{Name: "HLP"},
		{Name: "IDP registration"},
		{Name: "Other"},
	}
	ugPriority = []*types.SelectOption{
		{Name: "High"},
		{Name: "Medium"},
		{Name: "Low"},
	}
	ugICLADisplacementStatus = []*types.SelectOption{
		{Name: "Unregistered asylum seeker"},
		{Name: "Registered asylum seeker"},
		{Name: "Refugee"},
	}
	ugICLARSDDocuments = []*types.SelectOption{
		{Name: "Family Attesttion"},
		{Name: "Refugee ID"},
		{Name: "Asylum certificate"},
		{Name: "Rejection decision"},
		{Name: "Other"},
	}
	ugICLASpecificHLPConcern = []*types.SelectOption{
		{Name: "Housing"},
		{Name: "Land"},
		{Name: "Property"},
	}
	ugICLAHomeOwnership = []*types.SelectOption{
		{Name: "Own house"},
		{Name: "Rent"},
		{Name: "Other"},
	}
	ugICLAEvictionDocuments = []*types.SelectOption{
		{Name: "Eviction Notice"},
		{Name: "Other"},
	}
	ugICLANatureLandTenure = []*types.SelectOption{
		{Name: "Joint ownership"},
		{Name: "Co-ownership"},
		{Name: "Individual ownership"},
		{Name: "Other"},
	}
	ugICLANatureTenure = []*types.SelectOption{
		{Name: "Mailo"},
		{Name: "Lease"},
		{Name: "Freehold"},
		{Name: "Sustomary"},
	}
	ugICLATypeOfDocumentation = []*types.SelectOption{
		{Name: "Legal"},
		{Name: "Civil"},
	}
	ugICLATypeOfAgreement = []*types.SelectOption{
		{Name: "Oral"},
		{Name: "Written"},
	}
	ugICLATypeOfChallenge = []*types.SelectOption{
		{Name: "Employment"},
		{Name: "Business"},
	}
	ugICLATypeOfActions = []*types.SelectOption{
		{Name: "(1) Discussion with supervisor/team leader"},
		{Name: "(2) Conducting legal analysis, including the study of judicial practice"},
		{Name: "(3) Preparing letters, inquiries to various authorities"},
		{Name: "(4) Drafting of other legal documents (such leases or contracts)"},
		{Name: "(5) Lodging of a court application"},
		{Name: "(6) Attending of court session/hearing"},
		{Name: "(7) Review of the decision/appeal"},
		{Name: "(8) Negotiation"},
		{Name: "(9) Follow up with relevant administrative authority or other entities"},
		{Name: "(10) Accompaniment"},
		{Name: "(11) Other"},
	}
	ugICLATypesOfServices = []*types.SelectOption{
		{Name: "Legal counselling"},
		{Name: "Referral"},
		{Name: "Relocation"},
		{Name: "Livelihood"},
		{Name: "Business support"},
	}
	ugICLACaseClosureReason = []*types.SelectOption{
		{Name: "All action plan objectives achieved"},
		{Name: "Other priorities by beneficiary"},
		{Name: "Beneficiary unreachable"},
		{Name: "Other"},
	}
	ugAgreedFollowupMeans = []*types.SelectOption{
		{Name: "Schedule in-person meeting with beneficiary"},
		{Name: "Schedule phone call"},
		{Name: "Other"},
	}
	ugEILFSIncomeActivity = []*types.SelectOption{
		{Name: "Employment"},
		{Name: "Business"},
		{Name: "Skill training/Education"},
		{Name: "Farming"},
		{Name: "Bodaboda"},
		{Name: "Causal laborer"},
		{Name: "Other"},
	}
	ugEILFSStartupCapitalSource = []*types.SelectOption{
		{Name: "Donation (NGO, well-wishers, family & friends)"},
		{Name: "Member savings"},
		{Name: "Remittances"},
		{Name: "Loan/Credit"},
		{Name: "Other"},
	}
	ugEILFSProfitUse = []*types.SelectOption{
		{Name: "Reinvested"},
		{Name: "Shared to members"},
		{Name: "Loanable fund"},
		{Name: "Other"},
	}
	ugEILFSTrainingKind = []*types.SelectOption{
		{Name: "Record keeping/Tracking"},
		{Name: "Savings and loans management"},
		{Name: "Financial literacy"},
		{Name: "Business management"},
		{Name: "Agronomic practices"},
		{Name: "Farming as a business"},
		{Name: "Post-harvest handling"},
		{Name: "Skills training"},
		{Name: "Basic numeracy training"},
		{Name: "Lifeskills and entrepreneurship"},
		{Name: "Other"},
	}
	ugEILFSSkillsTraining = []*types.SelectOption{
		{Name: "Tailoring and garment cutting"},
		{Name: "Sewing"},
		{Name: "Arts & crafts"},
		{Name: "Agriculture"},
		{Name: "Briquettes making"},
		{Name: "Reusable packing production"},
		{Name: "Energy stoves"},
		{Name: "Hair dressing"},
		{Name: "Bakery and confectioneries"},
		{Name: "Motor mechanics and driving"},
		{Name: "Other"},
	}
	ugEILFSMarketChallenges = []*types.SelectOption{
		{Name: "Business registration"},
		{Name: "Taxes"},
		{Name: "Distance to markets"},
		{Name: "Product quality"},
		{Name: "Product certification"},
	}
	ugEILFSGroupBusinessArea = []*types.SelectOption{
		{Name: "Savings"},
		{Name: "Credit & loans"},
		{Name: "Arts & crafts"},
		{Name: "Skills training"},
		{Name: "Farming/Agriculture"},
		{Name: "Tailoring"},
		{Name: "Soap making"},
		{Name: "Bakery"},
		{Name: "Salon"},
		{Name: "Other"},
	}
	ugEILFSFinancialServices = []*types.SelectOption{
		{Name: "Credit"},
		{Name: "Savings"},
		{Name: "Insurance"},
		{Name: "Training"},
	}
	ugEILFSFinancialServiceProviders = []*types.SelectOption{
		{Name: "Banks"},
		{Name: "Micro-finance/micro-deposit institutions"},
		{Name: "SACCOs"},
		{Name: "VSLAs"},
		{Name: "Mobile money"},
		{Name: "Other"},
	}
	ugEducationVulnerability = []*types.SelectOption{
		{Name: "(1) Person with disability"},
		{Name: "(2) Living in child-headed household"},
		{Name: "(3) Child mother/father"},
		{Name: "(4) Orphan who's rights are not fulfilled"},
		{Name: "(5) Affected with chronic illness"},
		{Name: "(6) Experiencing forms of abuse"},
		{Name: "(7) Living on the street / Abandoned / Neglected"},
		{Name: "(8) In contact with the law"},
		{Name: "(9) Living in a (poverty-stricken) impoverished household"},
		{Name: "(10) Living with / elderly caregiver / guardian"},
		{Name: "(11) Affected by armed conflict"},
		{Name: "(12) Living in worst forms of labour"},
	}
	ugAction = []*types.SelectOption{
		{Name: "Open"},
		{Name: "Closed"},
		{Name: "Referred"},
		{Name: "Follow-up"},
	}
	ugReferralType = []*types.SelectOption{
		{Name: "External"},
		{Name: "Internal: Shelter/NFI"},
		{Name: "Internal: Livelihood/Food security"},
		{Name: "Internal: Education"},
		{Name: "Internal: WASH"},
		{Name: "Internal: Camp Management/UDOC"},
	}
	ugReferralReason = []*types.SelectOption{
		{Name: ""}, // TBD
	}
	ugReferralMeans = []*types.SelectOption{
		{Name: "(1) Phone"},
		{Name: "(2) E-mail"},
		{Name: "(3) Personal meeting"},
		{Name: "(4) Other"},
	}
	ugReferralPriority = []*types.SelectOption{
		{Name: "High (follow-up requested in 24 hours)"},
		{Name: "Medium (follow-up in 3 days)"},
		{Name: "Low (follow-up in 7 days)"},
	}
)
