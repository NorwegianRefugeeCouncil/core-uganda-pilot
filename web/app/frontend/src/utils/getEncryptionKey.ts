import base64 from 'react-native-base64'
import 'react-native-get-random-values'

export const getEncryptionKey = () => {
    const array = new Uint32Array(10);
    return base64.encode(crypto.getRandomValues(array).toString())
};
