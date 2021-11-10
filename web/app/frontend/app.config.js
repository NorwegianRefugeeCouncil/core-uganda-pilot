import 'dotenv/config'

export default ({ config }) => {
    switch (process.env.NODE_ENV) {
        case 'development':
        case 'dev':
            return {
                ...config,
                plugins: [
                    '@config-plugins/android-jsc-intl'
                ],
                extra: {
                    server_default_hostname: 'https://core.dev:8443',
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

