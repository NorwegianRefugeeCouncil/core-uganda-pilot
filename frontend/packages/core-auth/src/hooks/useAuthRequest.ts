import { AuthRequestConfig, AuthRequestPromptOptions, AuthSessionResult, DiscoveryDocument } from '../types/types';
import { AuthRequest } from '../types/authrequest';
import Browser from '../types/browser';

import useLoadedAuthRequest from './useLoadedAuthRequest';
import useAuthRequestResult from './useAuthRequestResult';

export default function useAuthRequest(
  config: AuthRequestConfig,
  discovery: DiscoveryDocument | null,
  browser: Browser,
): [AuthRequest | null, AuthSessionResult | null, (options?: AuthRequestPromptOptions) => Promise<AuthSessionResult>] {
  const request = useLoadedAuthRequest(config, discovery, AuthRequest, browser);
  const [result, promptAsync] = useAuthRequestResult(request, discovery);
  return [request, result, promptAsync];
}
