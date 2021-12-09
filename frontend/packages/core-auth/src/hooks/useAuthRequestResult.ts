import React from 'react';

import { AuthRequest } from '../types/authrequest';
import { AuthSessionResult, DiscoveryDocument, PromptMethod } from '../types/types';

export default function useAuthRequestResult(
  request: AuthRequest | null,
  discovery: DiscoveryDocument | null,
): [AuthSessionResult | null, PromptMethod] {
  const [authSessionResult, setAuthSessionResult] = React.useState<AuthSessionResult | null>(null);

  const promptAsync = React.useCallback(
    async () => {
      if (!discovery || !request) {
        throw new Error('Cannot prompt to authenticate until the request has finished loading.');
      }
      const result = await request?.promptAsync(discovery);
      setAuthSessionResult(result);
      return result;
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [request?.url, discovery?.authorization_endpoint],
  );

  return [authSessionResult, promptAsync];
}
