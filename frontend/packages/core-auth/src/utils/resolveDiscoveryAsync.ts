import axios from 'axios';

import { DiscoveryDocument, IssuerOrDiscovery } from '../types/types';

const discoveryPath = '.well-known/openid-configuration';
export default async function resolveDiscoveryAsync(issuerOrDiscovery: IssuerOrDiscovery): Promise<DiscoveryDocument | null> {
  let issuer: string;
  if (typeof issuerOrDiscovery === 'string') {
    issuer = issuerOrDiscovery;
  } else {
    issuer = issuerOrDiscovery.issuer;
  }
  const metadataEndpoint = issuer.endsWith('/') ? `${issuer}${discoveryPath}` : `${issuer}/${discoveryPath}`;
  return axios
    .get<DiscoveryDocument>(metadataEndpoint)
    .then((value) => value.data)
    .catch((err) => {
      return null;
    });
}
