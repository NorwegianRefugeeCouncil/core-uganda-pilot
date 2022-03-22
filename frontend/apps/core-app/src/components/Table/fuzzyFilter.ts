export function fuzzyTextFilterFn(rows, id, filterValue) {
  // TODO use fusejs
  // return matchSorter(rows, filterValue, { keys: [(row) => row.values[id]] });
}
fuzzyTextFilterFn.autoRemove = (val) => !val;
