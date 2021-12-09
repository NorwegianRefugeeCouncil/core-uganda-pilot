export default function getCurrentTimeInSeconds(): number {
  return Math.floor(Date.now() / 1000);
}
