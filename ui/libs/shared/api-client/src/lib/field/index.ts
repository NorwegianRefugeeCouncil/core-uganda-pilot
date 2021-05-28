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
    const p = NewPath(name, ...otherNames);
    p.root()._parent = this;
    return p;
  }

  index(idx: number): Path {
    return new Path({ index: idx + '', parent: this });
  }

  key(key: string): Path {
    return new Path({ index: key, parent: this });
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

export const NewPath = (name: string, ...otherNames: string[]): Path => {
  let path = new Path({ name });
  for (let otherName of otherNames) {
    path = new Path({ name: otherName, parent: path });
  }
  return path;
};

export const PathFrom = (path: string): Path => {
  let currentPath: Path;
  let indexer = false;
  let j = 0;
  for (let i = 0; i < path.length; i++) {
    const char = path[i];
    let done = false;

    if (indexer && char === ']') {
      indexer = false;
      const index = path.substring(j, i);
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
    }

    if (!indexer && !done && (char === '.' || char === '[' || i === path.length - 1)) {
      let name = path.substring(j, i + 1);
      if (name.endsWith('[') || name.endsWith('.')) {
        name = name.substring(0, name.length - 1);
      }
      if (!name) {
        j++;
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
      j = i + 1;
    }


  }

  return currentPath;
};
