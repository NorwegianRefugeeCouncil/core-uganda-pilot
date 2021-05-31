import { TypeMeta } from './types';
import { Path } from '../../../field';
import { ErrorList, Required } from '../../../field';

export const validateTypeMeta = (typeMeta: TypeMeta, fldPath: Path): ErrorList => {
  const errorList: ErrorList = [];
  if (typeMeta.kind === '') {
    errorList.push(Required(fldPath.child('kind'), 'kind is required'));
  }
  if (!typeMeta.apiVersion) {
    errorList.push(Required(fldPath.child('apiVersion'), 'apiVersion is required'));
  }
  return errorList;
};
