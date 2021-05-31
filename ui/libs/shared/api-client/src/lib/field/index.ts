export class Path {
  private _name: string;
  private _index: string;
  private _parent: Path;

  constructor(args: { name?: string, index?: string, parent?: Path }) {
    this._name = args.name;
    this._index = args.index;
    this._parent = args.parent;
  }

  root(): Path {
    if (!this._parent) {
      return this;
    }
    return this._parent.root();
  }

  child(name: string, ...otherNames: string[]): Path {
    const p = newPath(name, ...otherNames);
    p.root()._parent = this;
    return p;
  }

  index(idx: number): Path {
    return new Path({ index: idx + '', parent: this });
  }

  key(key: string): Path {
    return new Path({ index: key, parent: this });
  }

  set(obj: any, value: any): void {
    this.ensurePath(obj);

    let parent: any;
    if (!this._parent) {
      parent = obj;
    } else {
      parent = this._parent.get(obj);
    }

    let key: string;
    if (this._index) {
      key = this._index;
    } else {
      key = this._name;
    }
    parent[key] = value;
  }

  setIndexed(obj: any, key: string, value: any): void {
    this.ensurePath(obj);

    let parent = this.get(obj);
    if (!parent) {
      parent = [];
      this.set(obj, parent);
    }
    if (!Array.isArray(parent)) {
      throw 'cannot add value as object at path ' + this.string() + ' is not an array';
    }

    const arr = parent as any[];

    const keyValue = value[key];
    let found = false;
    for (let i = 0; i < arr.length; i++) {
      const arrElement = arr[i];
      const elementKeyValue = arrElement[key];
      if (elementKeyValue === keyValue) {
        arr[i] = value;
        found = true;
      }
    }

    if (!found) {
      arr.push(value);
    }
  }

  removeIndexed(obj: any, key: string, keyValue: string): void {
    this.ensurePath(obj);

    let parent = this.get(obj);
    if (!parent) {
      return;
    }
    if (!Array.isArray(parent)) {
      throw 'cannot add value as object at path ' + this.string() + ' is not an array';
    }

    const arr = parent as any[];

    for (let i = 0; i < arr.length; i++) {
      const arrElement = arr[i];
      const elementKeyValue = arrElement[key];
      if (elementKeyValue === keyValue) {
        arr.splice(i, 1);
        return;
      }
    }
  }

  ensurePath(obj: any) {
    let paths: Path[] = [];
    let currentPath = this as Path;
    while (currentPath) {
      paths.push(currentPath);
      currentPath = currentPath._parent;
    }
    paths = paths.reverse();

    let currentObj = obj;

    for (let i = 0; i < paths.length - 1; i++) {
      const path = paths[i];
      const nextPath = paths[i + 1];

      if (path.isArrayIndexer()) {
        const idx = path.getArrayIndex();
        if (!currentObj.hasOwnProperty(idx)) {
          if (nextPath && nextPath.isArrayIndexer()) {
            const nextObj = [];
            currentObj[idx] = nextObj;
            currentObj = nextObj;
          } else {
            const nextObj = {};
            currentObj[idx] = nextObj;
            currentObj = nextObj;
          }
        } else {
          currentObj = currentObj[idx];
        }
      }

      if (path.isObjectIndexer()) {
        const key = path.getObjectIndexer() as string;
        if (!currentObj.hasOwnProperty(key)) {
          if (nextPath && nextPath.isArrayIndexer()) {
            const nextObj = [];
            currentObj[key] = nextObj;
            currentObj = nextObj;
          } else {
            const nextObj = {};
            currentObj[key] = nextObj;
            currentObj = nextObj;
          }
        } else {
          currentObj = currentObj[key];
        }
      }
    }
  }

  private isArrayIndexer(): boolean {
    if (this._index === undefined) {
      return false;
    }
    const idx = parseInt(this._index);
    return !isNaN(idx);
  }

  private isObjectIndexer(): boolean {
    return !this.isArrayIndexer();
  }

  private getArrayIndex(): number {
    return parseInt(this._index);
  }

  private getObjectIndexer(): string {
    if (this._index) {
      return this._index;
    }
    return this._name;
  }

  add(obj: any, elem: any): Path {
    let value = this.get(obj);
    if (!value) {
      value = [];
      this.set(obj, value);
    }
    if (!Array.isArray(value)) {
      throw 'cannot add value as object at path ' + this.string() + ' is not an array';
    }
    const newLength = value.push(elem);
    return this.index(newLength - 1);
  }

  insert(obj: any, index: number, elem: any): Path {
    let value = this.get(obj) as any[];
    if (!value) {
      value = [];
      this.set(obj, value);
    }
    if (!Array.isArray(value)) {
      throw 'cannot add value as object at path ' + this.string() + ' is not an array';
    }

    value.splice(index, 0, elem);
    return this.index(index);
  }

  remove(obj: any): void {
    let value: any;
    if (!this._parent) {
      value = obj;
    } else {
      const parentPath = this._parent;
      value = parentPath.get(obj);
    }

    if (Array.isArray(value)) {
      const idx = parseInt(this._index);
      if (isNaN(idx)) {
        throw 'cannot remove value at path ' + this.string() + ' as index is not a number';
      }
      if (value.length - 1 < idx) {
        throw 'cannot remove value at path ' + this.string() + ' as index is out of bounds';
      }
      value.splice(idx, 1);
    } else if (typeof value === 'object') {
      const objValue = value as object;
      let key = this._name;
      if (this._index) {
        key = this._index;
      }
      if (!objValue.hasOwnProperty(key)) {
        return;
      }
      delete objValue[key];
    } else {
      throw 'cannot remove value as object at path ' + this.string() + ' is not an array or an object';
    }

  }

  get(obj: any): any {
    let paths: Path[] = [];
    let parent = this as Path;
    while (parent) {
      paths.push(parent);
      parent = parent._parent;
    }
    paths = paths.reverse();
    let value = obj;
    for (let i = 0; i < paths.length; i++) {
      const currentPath = paths[i];
      let key: string;
      if (currentPath._index) {
        key = currentPath._index;
      } else {
        key = currentPath._name;
      }
      if (!value) {
        return undefined;
      }
      if (!value.hasOwnProperty(key)) {
        return undefined;
      }
      value = value[key];
    }
    return value;
  }

  string(): string {
    const elems: Path[] = [];
    let current = this as Path;
    while (current) {
      elems.push(current);
      current = current._parent;
    }
    let path = '';
    for (let i = elems.length - 1; i >= 0; i--) {
      const current = elems[i];
      if (current._parent && current._name) {
        path = path + '.';
      }
      if (current._name) {
        path = path + current._name;
      } else {
        path = path + '[' + current._index + ']';
      }
    }
    return path;
  }
}

export const newPath = (name: string, ...otherNames: string[]): Path => {
  let path = new Path({ name });
  for (let otherName of otherNames) {
    path = new Path({ name: otherName, parent: path });
  }
  return path;
};

export const pathFrom = (path: string): Path => {
  let currentPath: Path;
  let indexer = false;
  let j = 0;

  for (let i = 0; i < path.length; i++) {

    const char = path[i];
    const nextChar = path[i + 1];

    if (char === ']') {
      j = i + 1;
      continue;
    }

    if (char === '.') {
      j = i + 1;
      continue;
    }

    if (char === '[') {
      indexer = true;
      j = i + 1;
      continue;
    }

    if (indexer && nextChar === ']') {
      indexer = false;
      let index = path.substring(j, i + 1);
      if (!currentPath) {
        currentPath = new Path({ index: index });
      } else {
        const idx = parseInt(index);
        if (isNaN(idx)) {
          currentPath = currentPath.key(index);
        } else {
          currentPath = currentPath.index(idx);
        }
      }
      j = i + 1;
      continue;
    }

    if (!indexer && (nextChar === '.' || nextChar === '[' || i === path.length - 1)) {
      let name = path.substring(j, i + 1);
      if (!name) {
        continue;
      }
      if (char === '[') {
        indexer = true;
      }
      if (!currentPath) {
        currentPath = new Path({ name: name });
      } else {
        currentPath = currentPath.child(name);
      }
      j = i;
      continue;
    }

  }

  return currentPath;
};

export type ErrorType =
  'FieldValueNotFound'
  | 'FieldValueRequired'
  | 'FieldValueDuplicate'
  | 'FieldValueInvalid'
  | 'FieldValueNotSupported'
  | 'FieldValueForbidden'
  | 'FieldValueTooLong'
  | 'FieldValueTooMany'
  | 'InternalError'

export class Error {
  public constructor(
    public type: ErrorType,
    public field: string,
    public badValue: any,
    public detail: string
  ) {
  }
}

export type ErrorList = Error[]

export const NotFound = (field: Path, value: any): Error => {
  return new Error('FieldValueNotFound', field.string(), value, '');
};

export const Required = (field: Path, detail: string): Error => {
  return new Error('FieldValueRequired', field.string(), '', detail);
};

export const Duplicate = (field: Path, value: any): Error => {
  return new Error('FieldValueDuplicate', field.string(), value, '');
};

export const Invalid = (field: Path, value: any, detail: string) => {
  return new Error('FieldValueInvalid', field.string(), value, detail);
};

