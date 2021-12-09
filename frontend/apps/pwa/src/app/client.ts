import {Client} from "core-api-client";
import axios from "axios";

export const axiosInstance = axios.create()
export default new Client(`${process.env.REACT_APP_SERVER_URL}`, axiosInstance)
