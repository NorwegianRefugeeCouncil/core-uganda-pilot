import {NotAllowedCheck} from "./models";

export type FunctionWithParametersType<P extends unknown[], R = void> = (
    ...args: P
) => R;

export const REGISTERED_ACTION_TYPES: { [actionType: string]: number } = {};

export function resetRegisteredActionTypes() {
    for (const key of Object.keys(REGISTERED_ACTION_TYPES)) {
        delete REGISTERED_ACTION_TYPES[key];
    }
}

export interface Action {
    type: string;
}

export type Creator<P extends any[] = any[],
    R extends object = object> = FunctionWithParametersType<P, R>;

export type ActionCreator<T extends string = string,
    C extends Creator = Creator> = C & TypedAction<T>;

export interface ActionCreatorProps<T> {
    _as: 'props';
    _p: T;
}

export type ActionType<A> = A extends ActionCreator<infer T, infer C>
    ? ReturnType<C> & { type: T }
    : never;

// declare to make it property-renaming safe
export declare interface TypedAction<T extends string> extends Action {
    readonly type: T;
}

function defineType<T extends string>(
    type: T,
    creator: Creator
): ActionCreator<T> {
    return Object.defineProperty(creator, 'type', {
        value: type,
        writable: false,
    }) as ActionCreator<T>;
}

export function props<P extends object>(): ActionCreatorProps<P> {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion, @typescript-eslint/naming-convention
    return {_as: 'props', _p: undefined!};
}

export function createAction<T extends string>(
    type: T
): ActionCreator<T, () => TypedAction<T>>;

export function createAction<T extends string, P extends object>(
    type: T,
    config: ActionCreatorProps<P> & NotAllowedCheck<P>
): ActionCreator<T, (props: P & NotAllowedCheck<P>) => P & TypedAction<T>>;

export function createAction<T extends string,
    P extends any[],
    R extends object>(
    type: T,
    creator: Creator<P, R & NotAllowedCheck<R>>
): FunctionWithParametersType<P, R & TypedAction<T>> & TypedAction<T>;

export function createAction<T extends string, C extends Creator>(
    type: T,
    config?: { _as: 'props' } | C
): ActionCreator<T> {
    REGISTERED_ACTION_TYPES[type] = (REGISTERED_ACTION_TYPES[type] || 0) + 1;

    if (typeof config === 'function') {
        return defineType(type, (...args: any[]) => ({
            ...config(...args),
            type,
        }));
    }
    const as = config ? config._as : 'empty';
    switch (as) {
        case 'empty':
            return defineType(type, () => ({type}));
        case 'props':
            return defineType(type, (props: object) => ({
                ...props,
                type,
            }));
        default:
            throw new Error('Unexpected config.');
    }
}

