import { AxiosResponse } from 'axios';

import { Response } from '../types';

const errorResponse = <TRequest, TBody>(request: TRequest, r: AxiosResponse<TBody>): Response<TRequest, TBody> => {
  return {
    request,
    response: undefined,
    status: r.request,
    statusCode: r.status,
    error: r.data as any,
    success: false,
  };
};

const successResponse = <TRequest, TBody>(request: TRequest, r: AxiosResponse<TBody>): Response<TRequest, TBody> => {
  return {
    request,
    response: r.data as TBody,
    status: r.statusText,
    statusCode: r.status,
    error: undefined,
    success: true,
  };
};

export const clientResponse = <TRequest, TBody>(
  r: AxiosResponse<TBody>,
  request: TRequest,
  expectedStatusCode: number,
): Response<TRequest, TBody> => {
  return r.status !== expectedStatusCode
    ? errorResponse<TRequest, TBody>(request, r)
    : successResponse<TRequest, TBody>(request, r);
};
