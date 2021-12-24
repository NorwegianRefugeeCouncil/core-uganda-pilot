export const getSessionStorage = (key: string): any | null => {
  const stored = sessionStorage.getItem(key);
  if (!stored) {
    return null;
  }
  return JSON.parse(stored);
};
