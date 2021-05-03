import {
  createContext,
  FunctionComponent,
  HTMLAttributes,
  useCallback,
  useMemo,
  useRef,
  useState
} from 'react';
import classNames from 'classnames';
import BaseModal, { ModalProps as BaseModalProps } from 'react-overlays/Modal';
import useCallbackRef from '@restart/hooks/useCallbackRef';
import transitionEnd from 'dom-helpers/transitionEnd';
import { useEventCallback } from '@restart/hooks';
import ModalManager from 'react-overlays/ModalManager';
import canUseDOM from 'dom-helpers/canUseDOM';
import { ownerDocument, removeEventListener } from 'dom-helpers';
import getScrollbarSize from 'dom-helpers/scrollbarSize';
import useWillUnmount from '@restart/hooks/useWillUnmount';

interface ModalContextProps {
  onHide: () => void
}

const ModalContext = createContext<ModalContextProps>({
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  onHide: () => {
  }
});

type ModalProps = {
  size?: 'lg' | 'sm' | 'xl'
  backdrop?: 'static' | true | false
  fade?: boolean
  handleHide?: () => void
  scrollable?: boolean
  show?: boolean
  verticallyCentered?: boolean
  fullscreen?: boolean
  fullscreenBelow?: 'sm' | 'md' | 'lg' | 'xl' | 'xxl'
} & BaseModalProps

export const Modal: FunctionComponent<ModalProps>
  = ({

       style,
       size,
       fullscreen,
       fullscreenBelow,

       fade,
       scrollable,
       verticallyCentered,
       className,
       children,

       /* BaseModal props*/
       show,
       animation,
       backdrop,
       keyboard,
       onEscapeKeyDown,
       onShow,
       onHide,
       container,
       autoFocus,
       enforceFocus,
       restoreFocus,
       restoreFocusOptions,
       onEntered,
       onExit,
       onExiting,
       onEnter,
       onEntering,
       onExited,
       backdropClassName,

       ...props
     }, ref) => {

  // We use a react context to wrap the Modal.

  const [modalStyle, setStyle] = useState({});
  const [animateStaticModal, setAnimateStaticModal] = useState(false);
  const waitingForMouseUpRef = useRef(false);
  const ignoreBackdropClickRef = useRef(false);
  const removeStaticModalAnimationRef = useRef<(() => void) | null>(null);

  const [modal, setModalRef] = useCallbackRef<{ dialog: any }>();
  const handleHide = useEventCallback(onHide);

  const modalContext = useMemo(() => ({ onHide: handleHide }), [handleHide]);

  function getModalManager() {
    return new ModalManager();
  }

  function updateDialogStyle(node) {

    if (!canUseDOM) {
      return;
    }

    const containerIsOverflowing = getModalManager().isContainerOverflowing(modal as any);
    const modalIsOverflowing = node.scrollHeight > ownerDocument(node).documentElement.clientHeight;

    setStyle({
      paddingRight:
        containerIsOverflowing && !modalIsOverflowing
          ? getScrollbarSize()
          : undefined,
      paddingLeft:
        !containerIsOverflowing && modalIsOverflowing
          ? getScrollbarSize()
          : undefined
    });
  }

  const handleWindowResize = useEventCallback(() => {
    if (modal) {
      updateDialogStyle(modal.dialog);
    }
  });

  useWillUnmount(() => {
    removeEventListener(window as any, 'resize', handleWindowResize);
    if (removeStaticModalAnimationRef.current) {
      removeStaticModalAnimationRef.current();
    }
  });

  const handleDialogMouseDown = () => {
    waitingForMouseUpRef.current = true;
  };

  const handleMouseUp = (e) => {
    if (waitingForMouseUpRef.current && modal && e.target === modal.dialog) {
      ignoreBackdropClickRef.current = true;
    }
    waitingForMouseUpRef.current = false;
  };

  const handleStaticModalAnimation = () => {
    setAnimateStaticModal(true);
    removeStaticModalAnimationRef.current = transitionEnd(
      modal?.dialog,
      () => {
        setAnimateStaticModal(false);
      }
    );
  };

  const handleStaticBackdropClick = (e) => {
    if (e.target !== e.currentTarget) {
      return;
    }
    handleStaticModalAnimation();
  };

  const handleClick = (e) => {
    if (backdrop === 'static') {
      handleStaticBackdropClick(e);
      return;
    }
    if (ignoreBackdropClickRef.current || e.target !== e.currentTarget) {
      ignoreBackdropClickRef.current = false;
      return;
    }

    if (onHide) {
      onHide();
    }

  };

  const handleEscapeKeyDown = (e) => {
    if (!keyboard && backdrop === 'static') {
      e.preventDefault();
      handleStaticModalAnimation();
    } else if (keyboard && onEscapeKeyDown) {
      onEscapeKeyDown(e);
    }
  };

  const handleEnter = (node) => {
    if (node) {
      node.style.display = 'block';
      updateDialogStyle(node);
    }
    if (onEnter) {
      onEnter(node);
    }
  };

  const handleExit = (node) => {
    if (removeStaticModalAnimationRef.current) {
      removeStaticModalAnimationRef.current();
    }
    if (onExit) {
      onExit(node);
    }
  };

  const handleEntering = (node) => {
    if (onEntering) {
      onEntering(node);
    }
    // eslint-disable-next-line no-restricted-globals
    addEventListener('resize', handleWindowResize);
  };

  const handleExited = (node) => {
    if (node) {
      node.style.display = '';
    }
    if (onExited) {
      onExited(node);
    }
    removeEventListener(window as any, 'resize', handleWindowResize);
  };

  className = classNames(
    className,
    'modal',
    { fade: fade }
  );

  if (scrollable) {
    className = classNames(className, 'modal-dialog-scrollable');
  }
  if (verticallyCentered) {
    className = classNames(className, 'modal-dialog-centered');
  }

  const renderBackdrop = useCallback((backdropProps) => {

    return <div
      {...backdropProps}
      className={classNames('modal-backdrop', backdropClassName, !animation && 'show')}
    />;

  }, [animation, backdropClassName]);

  const baseModalStyle = { ...style, ...modalStyle };

  if (!animation) {
    baseModalStyle.display = 'block';
  }

  const renderDialog = (dialogProps) => {

    return <div
      role='dialog'
      {...dialogProps}
      style={baseModalStyle}
      className={classNames(className, 'modal')}
      onClick={backdrop ? handleClick : undefined}
      onMouseUp={handleMouseUp}
    >
      <ModalDialog
        size={size}
        className={classNames(
          animateStaticModal && 'modal-static',
          fullscreen && 'modal-fullscreen',
          fullscreenBelow && 'modal-fullscreen-' + fullscreenBelow + '-down'
        )}
        onMouseDown={handleDialogMouseDown}
      >
        {children}
      </ModalDialog>
    </div>;
  };

  return (
    <ModalContext.Provider value={modalContext}>
      <BaseModal
        show={show}
        ref={setModalRef}
        backdrop={backdrop}
        container={container}
        keyboard
        autoFocus={autoFocus}
        enforceFocus={enforceFocus}
        restoreFocus={restoreFocus}
        restoreFocusOptions={restoreFocusOptions}
        onEscapeKeyDown={handleEscapeKeyDown}
        onShow={onShow}
        onHide={onHide}
        onEnter={handleEnter}
        onEntering={handleEntering}
        onEntered={onEntered}
        onExit={handleExit}
        onExiting={onExiting}
        onExited={handleExited}
        containerClassName={'modal-open'}
        renderBackdrop={renderBackdrop}
        renderDialog={renderDialog}
      >
        <div{...props} className={className}>
          {children}
        </div>
      </BaseModal>
    </ModalContext.Provider>);
};

Modal.defaultProps = {
  show: false,
  backdrop: true,
  keyboard: true,
  autoFocus: true,
  enforceFocus: true,
  restoreFocus: true,
  animation: false
};

type ModalDialogProps = HTMLAttributes<HTMLDivElement> & {
  size?: 'lg' | 'sm' | 'xl'
}

export const ModalDialog: FunctionComponent<ModalDialogProps> = ({ size, ...props }) => {
  return <div {...props} className={classNames(props.className, 'modal-dialog', size && 'modal-' + size)}>
    {props.children}
  </div>;
};

export const ModalContent: FunctionComponent<HTMLAttributes<HTMLDivElement>> = props => {
  return <div {...props} className={classNames(props.className, 'modal-content')}>
    {props.children}
  </div>;
};

export const ModalHeader: FunctionComponent<HTMLAttributes<HTMLDivElement>> = props => {
  return <div {...props} className={classNames(props.className, 'modal-header')}>
    {props.children}
  </div>;
};

export const ModalTitle: FunctionComponent<HTMLAttributes<HTMLDivElement>> = props => {
  return <h5 {...props} className={classNames(props.className, 'modal-title')}>
    {props.children}
  </h5>;
};

export const ModalBody: FunctionComponent<HTMLAttributes<HTMLDivElement>> = props => {
  return <h5 {...props} className={classNames(props.className, 'modal-body')}>
    {props.children}
  </h5>;
};
export const ModalFooter: FunctionComponent<HTMLAttributes<HTMLDivElement>> = props => {
  return <h5 {...props} className={classNames(props.className, 'modal-footer')}>
    {props.children}
  </h5>;
};
