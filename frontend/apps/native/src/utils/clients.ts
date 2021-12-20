import { useMemo } from 'react';
import Constants from 'expo-constants';
import { Client, ClientDefinition } from 'core-api-client';

export const apiClient = new Client(Constants.manifest?.extra?.server_uri);

// Doesn't need to be a hook, but it's not worth the time refactoring existing uses
export default function useApiClient(): ClientDefinition {
  return useMemo(() => apiClient, []);
}
