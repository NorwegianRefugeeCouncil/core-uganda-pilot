import 'dotenv/config'

export default ({ config }) => {
    switch (process.env.NODE_ENV) {
        case 'development':
        case 'dev':
            return {
                ...config,
                extra: {
                    server_default_hostname: 'localhost:9000',
                    server_hostname: process.env.SERVER_HOSTNAME,
                },
            }
        case 'production':
        case 'prod':
            return config
        case 'staging':
            return config
        default:
            return config
    }
}

