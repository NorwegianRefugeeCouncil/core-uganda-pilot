import * as axios from 'axios';

export type StringValue = { 'string': string; }

export type Value =
  StringValue
  | { 'float': string }
  | { int: string }
  | { 'boolean': string }
  | { 'null': boolean };

export function isStringValue(value: Value): value is StringValue {
  return (value as any)['string'] !== undefined;
}

export type Attributes = { [key: string]: Value };

export type Record = {
  id: string;
  revision: string;
  table: string;
  attributes: Attributes;
}

export type RecordList = {
  items: Record[];
}

export type GetRecordRequest = {
  id: string;
  table: string;
  revision: string;
}

export type GetRecordsRequest = {
  table: string;
  recordIds?: string[];
  revisions?: boolean;
}

export type PutRecordRequest = {
  isReplication?: boolean;
  record: Record
}

export type ColumnConstraint = { notNull: {}; }

export type Column = {
  name: string;
  type: string;
  constraints?: ColumnConstraint[];
}

export type TableConstraint = { primaryKey: { columns: string[] }; }

export type Table = {
  name: string;
  columns: Column[];
  constraints?: TableConstraint[];
}

export type Change = {
  sequence: number;
  recordId: string;
  table: string;
  revision: string;
}

export type Changes = {
  items: Change[];
}

export type GetChangesRequest = {
  since: number
}

export type GetTablesResponse = {
  items: GetTablesResponseItem[]
}

export type GetTablesResponseItem = {
  name: string
}

export class Client {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  public async getRecord(request: GetRecordRequest): Promise<Record> {
    try {
      const resp = await axios.default.get<any, axios.AxiosResponse<Record>>(`${this.baseURL}/apis/data.nrc.no/v1/tables/${request.table}/records/${request.id}`);
      return resp.data;
    } catch (e) {
      console.log(e);
      throw e;
    }
  }

  public async getRecords(request: GetRecordsRequest): Promise<RecordList> {
    try {
      const params: any = {};
      if (request?.revisions) {
        params.revisions = request.revisions;
      }
      const resp = await axios.default.get<any, axios.AxiosResponse<RecordList>>(`${this.baseURL}/apis/data.nrc.no/v1/tables/${request.table}/records`, {params});
      return resp.data;
    } catch (e) {
      console.log(e);
      throw e;
    }
  }

  public async putRecord(request: PutRecordRequest): Promise<Record> {
    try {
      const resp = await axios.default.put<any, axios.AxiosResponse<Record>>(`${this.baseURL}/apis/data.nrc.no/v1/tables/${request.record.table}/records/${request.record.id}`, request.record);
      return resp.data;
    } catch (e) {
      console.log(e);
      throw e;
    }
  }

  public async createTable(request: Table): Promise<Table> {
    try {
      const resp = await axios.default.put<any, axios.AxiosResponse<Table>>(`${this.baseURL}/apis/data.nrc.no/v1/tables/${request.name}`, request);
      return resp.data;
    } catch (e) {
      console.log(e);
      throw e;
    }
  }

  public async getChanges(request: GetChangesRequest): Promise<Changes> {
    try {
      const resp = await axios.default.get<any, axios.AxiosResponse<Changes>>(`${this.baseURL}/apis/data.nrc.no/v1/changes`, {params: {since: request.since}});
      return resp.data;
    } catch (e) {
      console.log(e);
      throw e;
    }
  }

  public async getTables(request: {}): Promise<GetTablesResponse> {
    try {
      const resp = await axios.default.get<any, axios.AxiosResponse<GetTablesResponse>>(`${this.baseURL}/apis/data.nrc.no/v1/tables`);
      return resp.data;
    } catch (e) {
      console.log(e);
      throw e;
    }
  }

  public async getTable(request: { tableName: string }): Promise<Table> {
    try {
      const resp = await axios.default.get<any, axios.AxiosResponse<Table>>(`${this.baseURL}/apis/data.nrc.no/v1/tables/${request.tableName}`);
      return resp.data;
    } catch (e) {
      console.log(e);
      throw e;
    }
  }

}
