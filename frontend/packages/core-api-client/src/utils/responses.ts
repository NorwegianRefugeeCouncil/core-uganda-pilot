import { AxiosError, AxiosResponse } from 'axios';

import {
  ErrorResponse,
  Response,
  SuccessResponse,
} from '../types/client/utils';

const errorResponse = <TRequest, TBody>(
  request: TRequest,
  r: AxiosError<TBody>,
): ErrorResponse<TRequest> => {
  const errorResp = r as any;
  return {
    request,
    response: undefined,
    status: errorResp.response?.statusText || '500 Internal Server Error',
    statusCode: errorResp.response?.status || 500,
    error: errorResp.response?.data,
    success: false,
  };
};

const successResponse = <TRequest, TBody>(
  request: TRequest,
  r: AxiosResponse<TBody>,
): SuccessResponse<TRequest, TBody> => {
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
  r: AxiosResponse<TBody> | AxiosError<TBody>,
  request: TRequest,
  expectedStatusCode: number,
): Response<TRequest, TBody> => {
  const resp = r as any;
  return resp.isAxiosError || resp.status !== expectedStatusCode
    ? errorResponse<TRequest, TBody>(request, resp)
    : successResponse<TRequest, TBody>(request, resp);
};
