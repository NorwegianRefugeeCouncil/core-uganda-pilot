import axios, { AxiosInstance, AxiosRequestConfig, Method } from 'axios';

import { RequestOptions, Response } from './types';
import { clientResponse } from './utils/responses';

export class BaseRESTClient {
  protected readonly axiosInstance: AxiosInstance;

  private token = '';

  lastInterceptorId: number | undefined;

  constructor(baseURL: string) {
    this.axiosInstance = axios.create({
      baseURL,
    });
    this.axiosInstance.interceptors.request.use((value: AxiosRequestConfig) => {
      if (!this.getToken()) {
        return value;
      }
      return {
        ...value,
        headers: {
          ...value.headers,
          Authorization: `Bearer ${this.getToken()}`,
        },
      };
    });
  }

  // this method became redundant, could be removed
  public setAuth = (token: string): void => {
    this.setToken(token);
  };

  public getToken = (): string => this.token;

  public setToken = (token: string): void => {
    this.token = token;
  };

  protected async do<TRequest, TBody>(
    request: TRequest,
    url: string,
    method: Method,
    data: unknown,
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
