import { FC, useEffect, useState } from 'react';
import { Link, useMatch } from 'react-router-dom';

import { SectionTitle } from '../sectiontitle/SectionTitle';
import { useApiClient } from '../../hooks/hooks';
import { OAuth2Client } from '../../types/types';

export const Clients: FC = () => {
  const apiClient = useApiClient();
  const match = useMatch('/clients/');
  const basePath = match?.pattern.path || '';

  const [clients, setClients] = useState<OAuth2Client[]>([]);
  useEffect(() => {
    if (!apiClient) {
      return;
    }
    apiClient.listOAuth2Clients({}).then((resp) => {
      if (resp.response) {
        setClients(resp.response.items);
      }
    });
  }, [apiClient]);

  return (
    <div className="container mt-3">
      <div className="row">
        <div className="col">
          <div className="card card-darkula">
            <div className="card-body">
              <SectionTitle title="OAuth2 Clients">
                <Link to={`${basePath}add`} className="btn btn-sm btn-success">
                  Add OAuth2 Client
                </Link>
              </SectionTitle>
              <div className="list-group list-group-darkula">
                {clients.length === 0 && (
                  <div className="list-group-item">No Clients</div>
                )}
                {clients.map((c) => (
                  <Link
                    key={c.id}
                    to={`${basePath}${c.id}`}
                    className="list-group-item"
                  >
                    {c.clientName}
                  </Link>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
