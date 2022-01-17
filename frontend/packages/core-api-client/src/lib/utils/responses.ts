import {AxiosError, AxiosResponse} from 'axios';

import { Response } from '../types';

const errorResponse = <TRequest, TBody>(request: TRequest, r: AxiosError<TBody>): Response<TRequest, TBody> => {
  const rany = r as any;
  console.log('ERROR RESPONSE', rany)
  return {
    request,
    response: undefined,
    status: rany.response.statusText || '500 Internal Server Error',
    statusCode: rany.response.status || 500,
    error: rany.response.data,
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
  r: AxiosResponse<TBody>|AxiosError<TBody>,
  request: TRequest,
  expectedStatusCode: number,
): Response<TRequest, TBody> => {
  // console.log('STATUS?', r.status !== expectedStatusCode, r.status, expectedStatusCode)
  // const isError = r instanceof AxiosError<TBody>;
  const rany = r as any;
  return rany.isAxiosError || rany.status !== expectedStatusCode
    ? errorResponse<TRequest, TBody>(request, rany)
    : successResponse<TRequest, TBody>(request, rany);
};
