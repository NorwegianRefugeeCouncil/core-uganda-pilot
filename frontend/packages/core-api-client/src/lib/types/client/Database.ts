import { Database, DatabaseList } from '..';

import { DataOperation, PartialObjectWrapper, Response } from './utils';

export type DatabaseCreateRequest = PartialObjectWrapper<Database>;
export type DatabaseCreateResponse = Response<DatabaseCreateRequest, Database>;

export type DatabaseListRequest = {} | undefined;
export type DatabaseListResponse = Response<DatabaseListRequest, DatabaseList>;

export interface DatabaseClientDefinition {
  create: DataOperation<DatabaseCreateRequest, DatabaseCreateResponse>;
  list: DataOperation<DatabaseListRequest, DatabaseListResponse>;
}
