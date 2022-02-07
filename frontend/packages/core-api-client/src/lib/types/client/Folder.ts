import { Folder, FolderList } from '..';

import { DataOperation, PartialObjectWrapper, Response } from './utils';

export type FolderListRequest = Record<string, never> | undefined;
export type FolderListResponse = Response<FolderListRequest, FolderList>;

export type FolderCreateRequest = PartialObjectWrapper<Folder>;
export type FolderCreateResponse = Response<FolderCreateRequest, Folder>;

export interface FolderClientDefinition {
  create: DataOperation<FolderCreateRequest, FolderCreateResponse>;
  list: DataOperation<FolderListRequest, FolderListResponse>;
}
