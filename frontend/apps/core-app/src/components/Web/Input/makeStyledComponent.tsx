/*
  This is a function copied from native base, but is not exported.
  It is required to apply the input style to a html input tag.
  It is not intended to be used outside of this component.
*/

import * as React from 'react';
import { useStyledSystemPropsResolver } from 'native-base';

export const makeStyledComponent = (Comp: any) => {
  // eslint-disable-next-line react/display-name
  return React.forwardRef((props: any, ref: any) => {
    const [style, restProps] = useStyledSystemPropsResolver(props);
    return (
      <Comp {...restProps} style={style} ref={ref}>
        {props.children}
      </Comp>
    );
  });
};
