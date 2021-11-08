import "@formatjs/intl-locale/polyfill";

import "@formatjs/intl-getcanonicallocales/polyfill";

import "@formatjs/intl-datetimeformat/polyfill";
import "@formatjs/intl-datetimeformat/add-golden-tz.js";
import "@formatjs/intl-datetimeformat/locale-data/en.js"; // USE YOUR OWN LANGUAGE OR MULTIPLE IMPORTS YOU WANT TO SUPPORT

import "@formatjs/intl-relativetimeformat/polyfill";
import "@formatjs/intl-relativetimeformat/locale-data/en.js"; // USE YOUR OWN LANGUAGE OR MULTIPLE IMPORTS YOU WANT TO SUPPORT

import "@formatjs/intl-numberformat/polyfill";
import "@formatjs/intl-numberformat/locale-data/en.js"; // USE YOUR OWN LANGUAGE OR MULTIPLE IMPORTS YOU WANT TO SUPPORT

import "@formatjs/intl-listformat/polyfill";
import "@formatjs/intl-listformat/locale-data/en.js"; // USE YOUR OWN LANGUAGE OR MULTIPLE IMPORTS YOU WANT TO SUPPORT

import "@formatjs/intl-displaynames/polyfill";
import "@formatjs/intl-displaynames/locale-data/en.js"; // USE YOUR OWN LANGUAGE OR MULTIPLE IMPORTS YOU WANT TO SUPPORT

import "@formatjs/intl-pluralrules/polyfill";
import "@formatjs/intl-pluralrules/locale-data/en.js"; // USE YOUR OWN LANGUAGE OR MULTIPLE IMPORTS YOU WANT TO SUPPORT


import {timezone} from "expo-localization";

if ('__setDefaultTimeZone' in Intl.DateTimeFormat) {
    Intl.DateTimeFormat.__setDefaultTimeZone(
        timezone
    );
}

import 'expo-dev-client';

import 'react-native-gesture-handler';
import {registerRootComponent} from 'expo';

import App from './App';

// registerRootComponent calls AppRegistry.registerComponent('main', () => App);
// It also ensures that whether you load the hooks in Expo Go or in a native build,
// the environment is set up appropriately
registerRootComponent(App);
