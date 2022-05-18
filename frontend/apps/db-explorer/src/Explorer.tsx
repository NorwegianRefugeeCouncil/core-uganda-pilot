import React, {useEffect, useState} from 'react';
import {Client, GetTablesResponseItem, Table} from './client';
import CreateTable from './CreateTable';
import {Link} from 'react-router-dom';

type ExplorerProps = {
  client: Client;
}

function Explorer(props: ExplorerProps) {

  const {client} = props;
  const [tables, setTables] = useState<GetTablesResponseItem[]>();

  useEffect(() => {
    client?.getTables({}).then(t => setTables(t.items));
  }, [client]);

  return (
    <div className="Explorer">
      <div className="container">
        <div className="row">
          <div className="col">
            <h1>DB Explorer</h1>
            <table className="table">
              <thead>
              <tr>
                <th>Table</th>
              </tr>
              </thead>
              <tbody>
              {
                tables?.map(t => <tr key={t.name}>
                  <td><Link to={`/tables/${t.name}`}>{t.name}</Link></td>
                </tr>)
              }
              </tbody>
            </table>
            <CreateTable
              onSubmit={(table) => {
                client?.createTable(table).then(() => {
                  client?.getTables({}).then(t => setTables(t.items));
                });
              }}
              onCancel={() => {

              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default Explorer;
