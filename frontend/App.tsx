import React from 'react';
import { Platform } from 'react-native';
import { CMSClient } from 'core-js-api-client';
import { Provider as PaperProvider } from 'react-native-paper';
import theme from './src/constants/theme';
import Layout from './src/components/Layout';

export default function App() {
  return (
    <PaperProvider theme={theme}>
      <Layout />
    </PaperProvider>
  );
}

const host = (Platform.OS === 'web' ? 'localhost' : '192.168.178.40') + ':9000';
const cmsClient = new CMSClient(host);
cmsClient.Cases().Get('dba43642-8093-4685-a197-f8848d4cbaaa').subscribe();