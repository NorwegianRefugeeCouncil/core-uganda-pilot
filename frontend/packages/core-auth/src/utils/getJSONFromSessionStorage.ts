export const getJSONFromSessionStorage = (key: string): any | undefined => {
  const stored = sessionStorage.getItem(key);
  if (stored) {
    return JSON.parse(stored);
  }
};
