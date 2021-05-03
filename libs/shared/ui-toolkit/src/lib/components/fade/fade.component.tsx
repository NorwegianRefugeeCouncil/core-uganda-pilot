import Transition, {
  ENTERED,
  ENTERING, TransitionProps
} from 'react-transition-group/Transition';
import React, { cloneElement, useCallback } from 'react';
import { triggerBrowserReflow, transitionEndListener } from '../../utils/utils';
import classNames from 'classnames';

const fadeStyles = {
  [ENTERING]: 'show',
  [ENTERED]: 'show'
};

const Fade = React.forwardRef<Transition<any>, TransitionProps<undefined>>(
  ({ className, children, ...props }, ref) => {

    const handleEnter = useCallback(
      args => {
        triggerBrowserReflow(args);
        if (props.onEnter) {
          props.onEnter(args, false);
        }
      }, [props]
    );

    return <Transition
      ref={ref}
      onEnter={handleEnter}
      addEndListener={transitionEndListener}
      {...props}
    >
      {(status, innerProps) =>
        cloneElement((children as any), {
          ...innerProps,
          className: classNames(
            'fade',
            className,
            (children as any).props.className,
            fadeStyles[status]
          )
        })
      }
    </Transition>;
  }
);

function wrapper<RefElement extends undefined | HTMLElement = undefined>() {
  return <Fade></Fade>;
}

function a() {
  <Fade></Fade>;
}

const defaultProps: Partial<TransitionProps> = {
  in: false,
  timeout: 300,
  mountOnEnter: false,
  unmountOnExit: false
};

Fade.defaultProps = defaultProps;

export default Fade;
