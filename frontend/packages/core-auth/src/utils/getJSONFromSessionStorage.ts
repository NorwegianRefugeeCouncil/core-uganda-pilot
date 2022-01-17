export const getJSONFromSessionStorage = (key: string): any | undefined => {
  const stored = sessionStorage.getItem(key);
  if (stored) {
    try {
      return JSON.parse(stored);
    } catch (e) {
      return undefined;
    }
  }
};
