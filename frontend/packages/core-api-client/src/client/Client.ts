import { ClientDefinition } from '../types';

import { BaseRESTClient } from './BaseRESTClient';
import { DatabaseClient } from './Database';
import { FolderClient } from './Folder';
import { FormClient } from './Form';
import { RecordClient } from './Record';
import { RecipientClient } from './Recipient';

export class Client extends BaseRESTClient implements ClientDefinition {
  static corev1 = 'apis/core.nrc.no/v1';

  public Database: DatabaseClient;

  public Folder: FolderClient;

  public Form: FormClient;

  public Record: RecordClient;

  public Recipient: RecipientClient;

  constructor(address: string) {
    super(`${address}/${Client.corev1}`);
    this.Database = new DatabaseClient(this);
    this.Folder = new FolderClient(this);
    this.Form = new FormClient(this);
    this.Record = new RecordClient(this, this.Form);
    this.Recipient = new RecipientClient(this.Record, this.Form);
  }
}
