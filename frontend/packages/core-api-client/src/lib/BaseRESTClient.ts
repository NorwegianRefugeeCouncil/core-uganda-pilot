import axios, {
  AxiosError,
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  Method,
} from 'axios';

import { RequestOptions, Response } from './types';
import { clientResponse } from './utils/responses';

export class BaseRESTClient {
  protected readonly axiosInstance: AxiosInstance;

  private token = '';

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
    let value: AxiosResponse<TBody> | AxiosError<TBody>;
    try {
      value = await this.axiosInstance.request<TBody>({
        responseType: 'json',
        method,
        url,
        data,
        headers,
        withCredentials: true,
      });
    } catch (err) {
      value = err as AxiosError<TBody>;
    }
    return clientResponse<TRequest, TBody>(value, request, expectStatusCode);
  }
}
