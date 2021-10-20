import {IAMClient} from "core-js-api-client";
import host from "../constants/host";

const iamClient = new IAMClient(
    'http',
    host,
    {
        'X-Authenticated-User-Subject': ['066a0268-fdc6-495a-9e4b-d60cfae2d81a']
    });

export default iamClient;
