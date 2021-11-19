import client, {ClientDefinition} from "@core/api-client";
import {useMemo} from "react";
import host from "../constants/host";
import axios from "axios";

export const axiosInstance = axios.create()

export default function useApiClient(): ClientDefinition {
    return useMemo(() => {
        return new client(host, axiosInstance)
    }, [1])
}
