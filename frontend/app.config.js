import 'dotenv/config';

export default {
    extra: {
        server_default_hostname: 'localhost:9000',
        server_hostname: process.env.SERVER_HOSTNAME,
    },
};
