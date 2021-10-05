import { StatusBar } from 'expo-status-bar';
import React from 'react';
import { Platform, StyleSheet, Text, View } from 'react-native';
import { CMSClient } from 'core-js-api-client';

export default function App() {
  return (
    <View style={styles.container}>
      <Text>Open up App.tsx to start working on your app!</Text>
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

const host = (Platform.OS === 'web' ? 'localhost' : '192.168.0.67') + ':9000';
const cmsClient = new CMSClient(host);
cmsClient.Cases().Get('dba43642-8093-4685-a197-f8848d4cbaaa').subscribe();
