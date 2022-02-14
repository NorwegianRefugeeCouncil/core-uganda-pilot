import { DatabaseClientDefinition } from './Database';
import { FolderClientDefinition } from './Folder';
import { FormClientDefinition } from './Form';
import { RecordClientDefinition } from './Record';

export interface ClientDefinition {
  Database: DatabaseClientDefinition;
  Folder: FolderClientDefinition;
  Form: FormClientDefinition;
  Record: RecordClientDefinition;
}
