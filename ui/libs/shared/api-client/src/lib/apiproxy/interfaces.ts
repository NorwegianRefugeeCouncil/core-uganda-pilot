import { Observable } from 'rxjs';

export interface Store extends Iterable<string> {
  getItem(key: string): string | null

  setItem(key: string, value: string)

  removeItem(key: string)

  events(): Observable<StorageEvent>

  clear()
}
