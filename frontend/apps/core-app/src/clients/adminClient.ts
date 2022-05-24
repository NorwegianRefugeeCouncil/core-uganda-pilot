import Constants from 'expo-constants';
import { AdminClient } from 'core-api-client';

export const adminClient = new AdminClient(
  Constants.manifest?.extra?.server_uri,
);
