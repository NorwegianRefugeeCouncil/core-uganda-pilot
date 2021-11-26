import * as Network from 'expo-network';

export const useNetworkState = async () => await Network.getNetworkStateAsync();
