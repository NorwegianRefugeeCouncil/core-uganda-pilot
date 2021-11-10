import 'expo-dev-client';

import 'react-native-gesture-handler';
import { registerRootComponent } from 'expo';
import { en, registerTranslation } from 'react-native-paper-dates';

import App from './App';

// registerRootComponent calls AppRegistry.registerComponent('main', () => App);
// It also ensures that whether you load the app in Expo Go or in a native build,
// the environment is set up appropriately
registerTranslation('en', en);
registerRootComponent(App);
