import Constants from 'expo-constants';
import { Platform } from 'react-native';

const host =
    Platform.OS === 'web'
        ? Constants.manifest?.extra?.server_default_hostname
        : Constants.manifest?.extra?.server_hostname;

export default host;
