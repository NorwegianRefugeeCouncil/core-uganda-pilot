import { ObjectMeta, TypeMeta } from '../meta/v1';

export interface RuntimeObject {
  apiVersion: string
  kind: string
}

export type RuntimeListItem = RuntimeObject & HasObjectMeta

export interface RuntimeList {
  apiVersion: string
  kind: string
  items: RuntimeListItem[]
}

export interface HasObjectMeta {
  metadata: ObjectMeta
}

export function isRuntimeObject(obj: any): obj is RuntimeObject {
  return (obj as TypeMeta).kind !== undefined
    && (obj as TypeMeta).apiVersion !== undefined;
}

export function isObjectMeta(obj: any): obj is HasObjectMeta {
  return (obj as HasObjectMeta).metadata !== undefined;
}

