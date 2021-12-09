import React from 'react';

import { DiscoveryDocument, IssuerOrDiscovery } from '../types/types';
import resolveDiscoveryAsync from '../utils/resolveDiscoveryAsync';

export default function useDiscovery(issuerOrDiscovery: IssuerOrDiscovery): DiscoveryDocument | null {
  console.log('ORIGINAL DISCOVERY');
  const [discoveryDocument, setDiscoveryDocument] = React.useState<DiscoveryDocument | null>(null);
  React.useEffect(() => {
    let isAllowed = true;
    resolveDiscoveryAsync(issuerOrDiscovery).then((discovery) => {
      if (isAllowed) {
        setDiscoveryDocument(discovery);
      }
    });
    return () => {
      isAllowed = false;
    };
  }, [issuerOrDiscovery]);
  return discoveryDocument;
}
