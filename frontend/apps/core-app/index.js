// Copied from expo-yarn-workspaces as eas build has issues finding __generated__/AppEntry.js

import 'expo/build/Expo.fx';
import { activateKeepAwake } from 'expo-keep-awake';
import registerRootComponent from 'expo/build/launch/registerRootComponent';

import App from './App';

// eslint-disable-next-line no-undef
if (__DEV__) {
  activateKeepAwake();
}

registerRootComponent(App);
