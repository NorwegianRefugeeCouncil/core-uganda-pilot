type Tree<T> = T & { childNodes: Tree<T>[] };

export const createDataTree = <T>(
  dataset: T[],
  idProperty: keyof T,
  parentIdProperty: keyof T,
): Tree<T>[] => {
  const hashTable = Object.create(null);
  dataset.forEach((aData) => {
    hashTable[aData[idProperty]] = { ...aData, childNodes: [] };
  });
  const dataTree: Tree<T>[] = [];
  dataset.forEach((aData) => {
    if (aData[parentIdProperty])
      hashTable[aData[parentIdProperty]].childNodes.push(
        hashTable[aData[idProperty]],
      );
    else dataTree.push(hashTable[aData[idProperty]]);
  });
  return dataTree;
};

export const getLeafNodes = <T>(
  nodes: Tree<T>[],
  result: Tree<T>[] = [],
): Tree<T>[] => {
  for (let i = 0, { length } = nodes; i < length; i++) {
    if (!nodes[i].childNodes || nodes[i].childNodes.length === 0) {
      result.push(nodes[i]);
    } else {
      result = getLeafNodes(nodes[i].childNodes, result);
    }
  }
  return result;
};
