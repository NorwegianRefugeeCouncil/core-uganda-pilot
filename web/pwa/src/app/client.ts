import {client as apiClient} from "core-js-api-client";
import axios from "axios";

export const axiosInstance = axios.create()
export default new apiClient("https://core.dev:8443", axiosInstance)
