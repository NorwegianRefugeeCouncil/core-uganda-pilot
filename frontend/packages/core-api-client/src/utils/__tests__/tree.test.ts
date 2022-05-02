import { createDataTree, getLeafNodes } from '../tree';

describe('createDataTree', () => {
  const makeTree = (parentIds: (string | null)[]) =>
    createDataTree(
      parentIds.map((id, index) => ({
        id: (index + 1).toString(),
        parentId: id,
      })),
      'id',
      'parentId',
    );

  it('should create a tree with a single root', () => {
    const tree = makeTree([null, '1', '1', '2', '2']);
    expect(tree).toEqual([
      {
        id: '1',
        parentId: null,
        childNodes: [
          {
            id: '2',
            parentId: '1',
            childNodes: [
              {
                id: '4',
                parentId: '2',
                childNodes: [],
              },
              {
                id: '5',
                parentId: '2',
                childNodes: [],
              },
            ],
          },
          {
            id: '3',
            parentId: '1',
            childNodes: [],
          },
        ],
      },
    ]);
  });

  it('should create a tree with multiple roots', () => {
    const tree = makeTree([null, null, '1', '1', '2', '2', '5']);
    expect(tree).toEqual([
      {
        id: '1',
        parentId: null,
        childNodes: [
          {
            id: '3',
            parentId: '1',
            childNodes: [],
          },
          {
            id: '4',
            parentId: '1',
            childNodes: [],
          },
        ],
      },
      {
        id: '2',
        parentId: null,
        childNodes: [
          {
            id: '5',
            parentId: '2',
            childNodes: [
              {
                id: '7',
                parentId: '5',
                childNodes: [],
              },
            ],
          },
          {
            id: '6',
            parentId: '2',
            childNodes: [],
          },
        ],
      },
    ]);
  });
});

describe('getLeafNodes', () => {
  it('should return a list of leaf nodes', () => {
    const dataset = [
      { id: '1', parentId: null },
      { id: '2', parentId: '1' },
      { id: '3', parentId: '1' },
      { id: '4', parentId: '2' },
      { id: '5', parentId: '2' },
    ];
    const tree = createDataTree(dataset, 'id', 'parentId');
    expect(getLeafNodes(tree)).toEqual([
      {
        id: '4',
        parentId: '2',
        childNodes: [],
      },
      {
        id: '5',
        parentId: '2',
        childNodes: [],
      },
      {
        id: '3',
        parentId: '1',
        childNodes: [],
      },
    ]);
  });
});
