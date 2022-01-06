export const getSessionStorage = (key: string): any | undefined => {
  const stored = sessionStorage.getItem(key);
  if (stored) {
    return JSON.parse(stored);
  }
};
