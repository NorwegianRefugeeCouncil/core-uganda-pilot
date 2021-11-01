// import {useApiClient} from "core-js-api-client";
import {client} from "core-js-api-client";
import {useMemo} from "react";
// import {useAuth} from "../../../../admin/src/app/hooks";
// import host from "../constants/host";



export default function useApiClient(): client {
    // const authCtx = useAuth()
    return useMemo(() => {
        return new client()
    }, [1])
}

// const apiClient = useApiClient();

// const iamClient = new useApiClient(
//     'http',
//     host,
//     {
//         'X-Authenticated-User-Subject': ['066a0268-fdc6-495a-9e4b-d60cfae2d81a']
//     });
//
// export default iamClient;

// export default apiClient;
