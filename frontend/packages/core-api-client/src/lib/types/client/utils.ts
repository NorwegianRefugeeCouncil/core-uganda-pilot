export type RequestOptions = {
  headers: { [key: string]: string };
  silentRedirect?: boolean;
};

export type Response<TRequest, TResponse> = {
  request: TRequest;
  response: TResponse | undefined;
  status: string;
  statusCode: number;
  success: boolean;
  error: any;
};

export type PartialObjectWrapper<T> = { object: Partial<T> };
export type DataOperation<TRequest, TResponse> = (
  request: TRequest,
) => Promise<TResponse>;
