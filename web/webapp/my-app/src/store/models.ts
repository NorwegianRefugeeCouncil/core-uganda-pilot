import {Action, ActionCreator, ActionType} from "./actions";

export type TypeId<T> = () => T;

export type InitialState<T> = Partial<T> | TypeId<Partial<T>> | void;

export interface ActionReducer<T, V extends Action = Action> {
    (state: T | undefined, action: V): T;
}

export type ActionReducerMap<T, V extends Action = Action> = {
    [p in keyof T]: ActionReducer<T[p], V>;
};

export interface ActionReducerFactory<T, V extends Action = Action> {
    (
        reducerMap: ActionReducerMap<T, V>,
        initialState?: InitialState<T>
    ): ActionReducer<T, V>;
}

export const arraysAreNotAllowedMsg =
    'arrays are not allowed in action creators';
type ArraysAreNotAllowed = typeof arraysAreNotAllowedMsg;

export const typePropertyIsNotAllowedMsg =
    'type property is not allowed in action creators';
type TypePropertyIsNotAllowed = typeof typePropertyIsNotAllowedMsg;

export const emptyObjectsAreNotAllowedMsg =
    'empty objects are not allowed in action creators';
type EmptyObjectsAreNotAllowed = typeof emptyObjectsAreNotAllowedMsg;

export type FunctionIsNotAllowed<T,
    ErrorMessage extends string> = T extends Function ? ErrorMessage : T;

export type NotAllowedCheck<T extends object> = T extends any[]
    ? ArraysAreNotAllowed
    : T extends { type: any }
        ? TypePropertyIsNotAllowed
        : keyof T extends never
            ? EmptyObjectsAreNotAllowed
            : unknown;
