import { Client } from 'core-api-client';
import Constants from 'expo-constants';

export const formsClient = new Client(Constants.manifest?.extra?.server_uri);
