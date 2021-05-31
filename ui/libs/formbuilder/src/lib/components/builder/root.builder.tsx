import { FormElement } from '@core/api-client';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { addValue } from '../../reducers';
import { FormElementContainer } from './element.container';

export type RootBuilderContainerProps = {
  root: FormElement,
  path: string
}


/**
 * Renders the builder for the root element of a form version
 * @param props
 * @param context
 * @constructor
 */
export const RootBuilderContainer: React.FC<RootBuilderContainerProps> = (props, context) => {

  const { root, path } = props;

  const dispatch = useDispatch();

  // Adds a new field to the root element
  const doAddField = () => {
    const action = addValue({ path: path + '.children', value: { type: 'shortText' } });
    dispatch(action);
  };

  return <div>

    {/* map each child of the root element */}
    {root?.children?.map((c, idx) => {
      return <FormElementContainer
        key={idx}
        element={c}
        path={path + '.children[' + idx + ']'}
      />;
    })}

    {/*button to add a new field to the root element*/}
    <button className={'btn btn-primary shadow-sm w-100'}
            onClick={() => doAddField()}>
      <i className={'bi bi-plus'} /> Add
    </button>

  </div>;
};
