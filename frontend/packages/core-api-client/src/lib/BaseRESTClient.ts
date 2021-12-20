import axios, { AxiosInstance, Method } from 'axios';

import { RequestOptions, Response } from './types';
import { clientResponse } from './utils/responses';

export class BaseRESTClient {
  protected readonly axiosInstance: AxiosInstance;

  constructor(baseURL: string) {
    this.axiosInstance = axios.create({
      baseURL,
    });
  }

  public setAuth = (token: string): void => {
    this.axiosInstance.interceptors.request.use((value: any) => {
      const result = { ...value };
      if (!token) {
        return value;
      }
      return {
        ...result,
        headers: {
          ...result.headers,
          Authorization: `Bearer ${token}`,
        },
      };
    });
  };

  protected async do<TRequest, TBody>(
    request: TRequest,
    url: string,
    method: Method,
    data: any,
    expectStatusCode: number,
    options?: RequestOptions,
  ): Promise<Response<TRequest, TBody>> {
    const headers: { [key: string]: string } = options?.headers ?? {
      Accept: 'application/json',
    };
    try {
      const value = await this.axiosInstance.request<TBody>({
        responseType: 'json',
        method,
        url,
        data,
        headers,
        withCredentials: true,
      });
      return clientResponse<TRequest, TBody>(value, request, expectStatusCode);
    } catch (err) {
      return {
        request,
        response: undefined,
        status: '500 Internal Server Error',
        statusCode: 500,
        error: axios.isAxiosError(err) ? err.message : 'Unknown error',
        success: false,
      };
    }
  }
}
