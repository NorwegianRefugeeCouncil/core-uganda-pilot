import { AxiosResponse } from "axios";
import { Response } from "../types/types";
export declare function clientResponse<TRequest, TBody>(r: AxiosResponse<TBody>, request: TRequest, expectedStatusCode: number): Response<TRequest, TBody>;
