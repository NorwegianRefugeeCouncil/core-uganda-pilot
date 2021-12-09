import * as Network from 'expo-network';

export const getNetworkState = async () => {
  return await Network.getNetworkStateAsync();
};
