import Constants from 'expo-constants';
import { Client } from 'core-api-client';

export const formsClient = new Client(Constants.manifest?.extra?.server_uri);
