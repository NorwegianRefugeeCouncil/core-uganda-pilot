import { StatusBar } from 'expo-status-bar';
import React from 'react';
import { Button, Platform, StyleSheet, View } from 'react-native';
import { CMSClient } from 'core-js-api-client';
import Constants from 'expo-constants';

export default function App() {
  get();
  return (
    <View style={styles.container}>
      <Button onPress={get} title="Initiate an API request" />
      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#ecf',
    alignItems: 'center',
    justifyContent: 'center'
  }
});

const host = (Platform.OS === 'web' ? Constants.manifest?.extra?.server_default_hostname : Constants.manifest?.extra?.server_hostname);
const cmsClient = new CMSClient(host);

function get() {
  cmsClient.Cases().Get('dba43642-8093-4685-a197-f8848d4cbaaa').subscribe();
}
