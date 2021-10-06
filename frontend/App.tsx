import { StatusBar } from 'expo-status-bar';
import React from 'react';
import { Platform, Image, View } from 'react-native';
import { CMSClient } from 'core-js-api-client';
import { Appbar, Provider as PaperProvider, Text } from 'react-native-paper';
import theme from './src/constants/theme';
import layout from './src/styles/layout';
import common from './src/styles/common';
import nrc from './assets/nrc.svg';
import { isMobile } from 'react-device-detect';

export default function App() {
  return (
    <PaperProvider theme={theme}>
      <View style={layout.container}>
        {isMobile && (
          <View style={layout.navigation} />
        )}
        <View style={layout.body}>
        {!isMobile && (
          <Appbar style={common.top} theme={{ dark: true }}>
            <Image source={nrc} style={common.logo}/>
            <Appbar.Content title="Title" subtitle={'Subtitle'} />
            <Appbar.Action
              icon="archive"
              onPress={() => console.log('Pressed archive')}
            />
            <Appbar.Action icon="mail" onPress={() => console.log('Pressed mail')} />
            <Appbar.Action icon="label" onPress={() => console.log('Pressed label')} />
            <Appbar.Action
              icon="delete"
              onPress={() => console.log('Pressed delete')}
            />
          </Appbar>
        )}
          <Text>Open up App.tsx to start working on your app!</Text>
          <StatusBar style="auto" />
        </View>
      </View>
    </PaperProvider>
  );
}

const host = (Platform.OS === 'web' ? 'localhost' : '192.168.178.40') + ':9000';
const cmsClient = new CMSClient(host);
cmsClient.Cases().Get('dba43642-8093-4685-a197-f8848d4cbaaa').subscribe();