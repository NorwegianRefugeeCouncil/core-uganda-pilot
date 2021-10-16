import {Action, ActionCreator, ActionType} from "./actions";
import {ActionReducer} from "./models";

export interface ReducerTypes<State,
    Creators extends readonly ActionCreator[]> {
    reducer: OnReducer<State, Creators>;
    types: ExtractActionTypes<Creators>;
}

// Specialized Reducer that is aware of the Action type it needs to handle
export interface OnReducer<State, Creators extends readonly ActionCreator[]> {
    (state: State, action: ActionType<Creators[number]>): State;
}

type ExtractActionTypes<Creators extends readonly ActionCreator[]> = {
    [Key in keyof Creators]: Creators[Key] extends ActionCreator<infer T>
        ? T
        : never;
};

export function on<State, Creators extends readonly ActionCreator[]>(
    ...args: [
        ...creators: Creators,
        reducer: OnReducer<State extends infer S ? S : never, Creators>
    ]
): ReducerTypes<State, Creators> {
    // This could be refactored when TS releases the version with this fix:
    // https://github.com/microsoft/TypeScript/pull/41544
    const reducer = args.pop() as OnReducer<any, Creators>;
    const types = (((args as unknown) as Creators).map(
        (creator) => creator.type
    ) as unknown) as ExtractActionTypes<Creators>;
    return {reducer, types};
}

export function createReducer<S, A extends Action = Action>(
    initialState: S,
    ...ons: ReducerTypes<S, readonly ActionCreator[]>[]
): ActionReducer<S, A> {
    const map = new Map<string, OnReducer<S, ActionCreator[]>>();
    for (const on of ons) {
        for (const type of on.types) {
            const existingReducer = map.get(type);
            if (existingReducer) {
                const newReducer: typeof existingReducer = (state, action) =>
                    on.reducer(existingReducer(state, action), action);
                map.set(type, newReducer);
            } else {
                map.set(type, on.reducer);
            }
        }
    }

    return function (state: S = initialState, action: A): S {
        const reducer = map.get(action.type);
        return reducer ? reducer(state, action) : state;
    };
}
