import {client} from "core-js-api-client";
import {useMemo} from "react";
import host from "../constants/host";

export default function useApiClient(): client {
    return useMemo(() => {
        return new client(host)
    }, [1])
}
