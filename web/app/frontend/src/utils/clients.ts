import client, {ClientDefinition} from "core-js-api-client";
import {useMemo} from "react";
import axios from "axios";
import Constants from "expo-constants";

export const axiosInstance = axios.create()

export default function useApiClient(): ClientDefinition {
    return useMemo(() => {
        return new client(Constants.manifest?.extra?.server_uri, axiosInstance)
    }, [1])
}
