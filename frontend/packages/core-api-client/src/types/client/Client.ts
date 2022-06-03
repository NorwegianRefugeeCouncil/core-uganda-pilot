import { DatabaseClientDefinition } from './Database';
import { FolderClientDefinition } from './Folder';
import { FormClientDefinition } from './Form';
import { RecordClientDefinition } from './Record';
import { RecipientClientDefinition } from './Recipient';
import { IdentityClientDefinition } from './Identity';

export interface ClientDefinition {
  Database: DatabaseClientDefinition;
  Folder: FolderClientDefinition;
  Form: FormClientDefinition;
  Record: RecordClientDefinition;
  Recipient: RecipientClientDefinition;
  Identity: IdentityClientDefinition;
}
