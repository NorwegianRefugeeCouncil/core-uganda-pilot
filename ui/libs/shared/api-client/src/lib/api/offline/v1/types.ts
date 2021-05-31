import { Status, TypeMeta } from '../../meta';
import { HasObjectMeta } from '../../runtime';

export interface Operation extends TypeMeta, HasObjectMeta {
  spec: OperationSpec
  status?: OperationStatus
}

export type OperationType =
  'create' | 'delete' | 'update'

export interface OperationSpec {
  batchId?: string
  operationType?: OperationType
  entity?: TypeMeta & HasObjectMeta
}

export interface OperationStatus {
  conditions: OperationCondition[]
}

export interface OperationCondition {
  type: OperationConditionType
  status: OperationConditionStatus
  reason?: string
  message?: string
  response?: Status
}

export type OperationConditionStatus =
  'True'
  | 'False'
  | 'Unknown'

export type OperationConditionType =
  'OperationQueued'
  | 'OperationPending'
  | 'OperationAccepted'

export const isInCondition = (operation: Operation, conditionType: OperationConditionType, status: OperationConditionStatus): boolean => {
  if (!operation?.status?.conditions){
    return false
  }
  for (let condition of operation.status.conditions) {
    if (condition.type === conditionType) {
      return condition.status === status;
    }
  }
  return false;
};

export const setCondition = (operation: Operation, conditionType: OperationConditionType, status: OperationConditionStatus) => {
  const condition = getOrCreateCondition(operation, conditionType);
  condition.status = status;
};

export const getOrCreateCondition = (operation: Operation, conditionType: OperationConditionType): OperationCondition => {
  if (!operation.status) {
    operation.status = {
      conditions: []
    };
  }
  if (!operation.status.conditions) {
    operation.status.conditions = [];
  }
  let found: OperationCondition;
  for (let c of operation.status.conditions) {
    if (c.type === conditionType) {
      found = c;
      break;
    }
  }
  if (!found) {
    found = {
      type: conditionType,
      status: 'Unknown'
    };
    operation.status.conditions.push(found);
  }
  return found;
};

export const isQueued = (operation: Operation): boolean => {
  return isInCondition(operation, 'OperationQueued', 'True');
};

export const isPending = (operation: Operation): boolean => {
  return isInCondition(operation, 'OperationPending', 'True');
};

export const isAccepted = (operation: Operation): boolean => {
  return isInCondition(operation, 'OperationAccepted', 'True');
};

export const isRejected = (operation: Operation): boolean => {
  return isInCondition(operation, 'OperationAccepted', 'False');
};
