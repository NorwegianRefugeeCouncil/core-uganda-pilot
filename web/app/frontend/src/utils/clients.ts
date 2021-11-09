import {client} from "core-js-api-client";
import {useMemo} from "react";
import host from "../constants/host";
import axios from "axios";

export const axiosInstance = axios.create()

export default function useApiClient(): client {
    return useMemo(() => {
        return new client(host, axiosInstance)
    }, [1])
}
