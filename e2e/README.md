# Cypress E2E Test Suide

## Dev-Setup

1. Create `.env` file in the root of the project `./E2E/.env`
2. Copy this content to the `.env` file with your credentials
```bash
# These variables will be read automatically
# CYPRESS_ part will be removed,
# access them with Cypress.env("USERNAME"), Cypress.env("PASSWORD")
# Changes in here requires Cypress restart (Stop and `npm run open`)
CYPRESS_USERNAME=user@fuu.bar
CYPRESS_PASSWORD=fuubar
```
3. Install and Run
```bash
npm install
npm run open
```