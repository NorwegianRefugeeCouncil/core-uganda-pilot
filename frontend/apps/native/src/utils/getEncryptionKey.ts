import 'react-native-get-random-values';

export const getEncryptionKey = () => {
  const array = new Uint32Array(10);
  return crypto.getRandomValues(array).toString();
};
