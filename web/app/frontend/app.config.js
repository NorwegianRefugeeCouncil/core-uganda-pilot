import 'dotenv/config';

const config = setConfig();

function setConfig() {
	switch (process.env.NODE_ENV) {
		case 'development':
		case 'dev':
			return {
				extra: {
					server_default_hostname: 'localhost:9000',
					server_hostname: process.env.SERVER_HOSTNAME,
				},
			}
		case 'production':
		case 'prod':
			return {}
		case 'staging':
			return {}
		default:
			return {}
	}
}

export default config;
