import 'dotenv/config';


const runtimeDefault = {
    plugins: [
        "@config-plugins/android-jsc-intl"
    ],
    extra: {
        server_hostname: 'localhost:9000',
    }
}

export default ({config}) => {
    switch (process.env.NODE_ENV) {
        case 'development':
        case 'dev':
            return {
                ...config,
                ...runtimeDefault,
                extra: {
                    server_hostname: process.env.SERVER_HOSTNAME,
                },
            }
        case 'production':
        case 'prod':
            return {...config, ...runtimeDefault}
        case 'staging':
            return {...config, ...runtimeDefault}
        default:
            return {...config, ...runtimeDefault}
    }
}

