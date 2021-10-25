import {filter, Observable, OperatorFunction} from "rxjs";
import {Action, ActionCreator, Creator} from "./actions";

export interface EffectConfig {
    dispatch?: boolean;
    useEffectsErrorHandler?: boolean;
}

export const DEFAULT_EFFECT_CONFIG: Readonly<Required<EffectConfig>> = {
    dispatch: true,
    useEffectsErrorHandler: true,
};

export const CREATE_EFFECT_METADATA_KEY = '__@ngrx/effects_create__';

export interface CreateEffectMetadata {
    [CREATE_EFFECT_METADATA_KEY]: EffectConfig;
}

type DispatchType<T> = T extends { dispatch: infer U } ? U : true;
type ObservableType<T, OriginalType> = T extends false ? OriginalType : Action;
type EffectResult<OT> = Observable<OT> | ((...args: any[]) => Observable<OT>);
type ConditionallyDisallowActionCreator<DT, Result> = DT extends false
    ? unknown // If DT (DispatchType is false, then we don't enforce any return types)
    : Result extends EffectResult<infer OT>
        ? OT extends ActionCreator
            ? 'ActionCreator cannot be dispatched. Did you forget to call the action creator function?'
            : unknown
        : unknown;

export function createEffect<C extends EffectConfig,
    DT extends DispatchType<C>,
    OT extends ObservableType<DT, OT>,
    R extends EffectResult<OT>>(
    source: () => R & ConditionallyDisallowActionCreator<DT, R>,
    config?: Partial<C>
): R & CreateEffectMetadata {
    const effect = source();
    const value: EffectConfig = {
        ...DEFAULT_EFFECT_CONFIG,
        ...config, // Overrides any defaults if values are provided
    };
    Object.defineProperty(effect, CREATE_EFFECT_METADATA_KEY, {
        value,
    });
    return effect as typeof effect & CreateEffectMetadata;
}

// Module-private helper type
type ActionExtractor<T extends string | AC,
    AC extends ActionCreator<string, Creator>,
    E> = T extends string ? E : ReturnType<Extract<T, AC>>;

export function ofType<AC extends ActionCreator<string, Creator>[],
    U extends Action = Action,
    V = ReturnType<AC[number]>>(...allowedTypes: AC): OperatorFunction<U, V>;

export function ofType<E extends Extract<U, { type: T1 }>,
    AC extends ActionCreator<string, Creator>,
    T1 extends string | AC,
    U extends Action = Action,
    V = T1 extends string ? E : ReturnType<Extract<T1, AC>>>(t1: T1): OperatorFunction<U, V>;
export function ofType<E extends Extract<U, { type: T1 | T2 }>,
    AC extends ActionCreator<string, Creator>,
    T1 extends string | AC,
    T2 extends string | AC,
    U extends Action = Action,
    V = ActionExtractor<T1 | T2, AC, E>>(t1: T1, t2: T2): OperatorFunction<U, V>;
export function ofType<E extends Extract<U, { type: T1 | T2 | T3 }>,
    AC extends ActionCreator<string, Creator>,
    T1 extends string | AC,
    T2 extends string | AC,
    T3 extends string | AC,
    U extends Action = Action,
    V = ActionExtractor<T1 | T2 | T3, AC, E>>(t1: T1, t2: T2, t3: T3): OperatorFunction<U, V>;
export function ofType<E extends Extract<U, { type: T1 | T2 | T3 | T4 }>,
    AC extends ActionCreator<string, Creator>,
    T1 extends string | AC,
    T2 extends string | AC,
    T3 extends string | AC,
    T4 extends string | AC,
    U extends Action = Action,
    V = ActionExtractor<T1 | T2 | T3 | T4, AC, E>>(t1: T1, t2: T2, t3: T3, t4: T4): OperatorFunction<U, V>;
export function ofType<E extends Extract<U, { type: T1 | T2 | T3 | T4 | T5 }>,
    AC extends ActionCreator<string, Creator>,
    T1 extends string | AC,
    T2 extends string | AC,
    T3 extends string | AC,
    T4 extends string | AC,
    T5 extends string | AC,
    U extends Action = Action,
    V = ActionExtractor<T1 | T2 | T3 | T4 | T5, AC, E>>(t1: T1, t2: T2, t3: T3, t4: T4, t5: T5): OperatorFunction<U, V>;
/**
 * Fallback for more than 5 arguments.
 * There is no inference, so the return type is the same as the input -
 * Observable<Action>.
 *
 * We provide a type parameter, even though TS will not infer it from the
 * arguments, to preserve backwards compatibility with old versions of ngrx.
 */
export function ofType<V extends Action>(
    ...allowedTypes: Array<string | ActionCreator<string, Creator>>
): OperatorFunction<Action, V>;
/**
 * `ofType` filters an Observable of `Actions` into an Observable of the actions
 * whose type strings are passed to it.
 *
 * For example, if `actions` has type `Actions<AdditionAction|SubstractionAction>`, and
 * the type of the `Addition` action is `add`, then
 * `actions.pipe(ofType('add'))` returns an `Observable<AdditionAction>`.
 *
 * Properly typing this function is hard and requires some advanced TS tricks
 * below.
 *
 * Type narrowing automatically works, as long as your `actions` object
 * starts with a `Actions<SomeUnionOfActions>` instead of generic `Actions`.
 *
 * For backwards compatibility, when one passes a single type argument
 * `ofType<T>('something')` the result is an `Observable<T>`. Note, that `T`
 * completely overrides any possible inference from 'something'.
 *
 * Unfortunately, for unknown 'actions: Actions' these types will produce
 * 'Observable<never>'. In such cases one has to manually set the generic type
 * like `actions.ofType<AdditionAction>('add')`.
 *
 * @usageNotes
 *
 * Filter the Actions stream on the "customers page loaded" action
 *
 * ```ts
 * import { ofType } from '@ngrx/effects';
 * import * fromCustomers from '../customers';
 *
 * this.actions$.pipe(
 *  ofType(fromCustomers.pageLoaded)
 * )
 * ```
 */
export function ofType(
    ...allowedTypes: Array<string | ActionCreator<string, Creator>>
): OperatorFunction<Action, Action> {
    return filter((action: Action) =>
        allowedTypes.some((typeOrActionCreator) => {
            if (typeof typeOrActionCreator === 'string') {
                // Comparing the string to type
                return typeOrActionCreator === action.type;
            }

            // We are filtering by ActionCreator
            return typeOrActionCreator.type === action.type;
        })
    );
}

