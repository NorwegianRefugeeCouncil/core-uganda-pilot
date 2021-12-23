import axios, { AxiosInstance, AxiosRequestConfig, Method } from 'axios';

import { RequestOptions, Response } from './types';
import { clientResponse } from './utils/responses';

export class BaseRESTClient {
  protected readonly axiosInstance: AxiosInstance;

  lastInterceptorId: number | undefined;

  constructor(baseURL: string) {
    this.axiosInstance = axios.create({
      baseURL,
    });
  }

  public setAuth = (token: string): void => {
    console.log('Setting token to:', token);

    if (this.lastInterceptorId != null)
      this.axiosInstance.interceptors.request.eject(this.lastInterceptorId);

    this.lastInterceptorId = this.axiosInstance.interceptors.request.use(
      (value: AxiosRequestConfig) => {
        if (!token) {
          return value;
        }
        const headers = {
          ...value.headers,
          Authorization: `Bearer ${token}`,
        };
        value.headers = headers;
        console.log(value.headers.Authorization);
        console.log(value.url);
        return value;
      },
    );
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
