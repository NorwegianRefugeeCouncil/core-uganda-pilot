import * as Random from 'expo-random';

export const getEncryptionKey = () => {
    const bytes = Random.getRandomBytes(40);
    return new TextDecoder().decode(bytes);
};
