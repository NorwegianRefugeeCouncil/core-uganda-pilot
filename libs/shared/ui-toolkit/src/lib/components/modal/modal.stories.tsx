import * as React from 'react'
import { storiesOf } from '@storybook/react';
import { Modal, ModalBody, ModalContent, ModalDialog, ModalFooter, ModalHeader, ModalTitle } from './modal';
import { Button } from '../button/button.component';
import { ModalProps } from 'react-overlays/Modal';

const content: (setShow: (show: boolean) => void) => any = (setShow) => {
  return (<ModalContent>
    <ModalHeader>
      <ModalTitle>Edit Properties</ModalTitle>
    </ModalHeader>
    <ModalBody>
      <p>Here is some content</p>
    </ModalBody>
    <ModalFooter>
      <Button onClick={() => setShow(false)}>Close</Button>
    </ModalFooter>
  </ModalContent>);
};

type StoryProps = {
  storyName: string
} & ModalProps

const ModalStory = (props: React.PropsWithChildren<StoryProps>) => {
  const { storyName, children, ...otherProps } = props;
  const [show, setShow] = React.useState(false);
  return (

    <div className={"mb-2 border p-3"}>
      <Button className={"mb-2"} kind={'primary'} onClick={() => setShow(true)}>{props.storyName}</Button>
      <Modal
        {...otherProps}
        show={show}
        onHide={() => setShow(false)}>
        {content(setShow)}
      </Modal>
      {children}
    </div>

  );
};

storiesOf('Modal', module)
  .add('default', () => {
    return (
      <>
        <ModalStory storyName={'Default'}>
          <p>Just a basic, default modal</p>
        </ModalStory>
        <ModalStory storyName={'Static Backdrop'} backdrop={'static'}>
          <p>A modal with a <code>static</code> backdrop. That means that a click on the backdrop will not close the modal</p>
        </ModalStory>
        <ModalStory storyName={'No Keyboard'} keyboard={false} backdrop={'static'}>
          <p>A modal with a <code>keyboard=false</code> and <code>backdrop=static</code> backdrop. That means that pressing the <code>esc</code> key will not close the backdrop</p>
        </ModalStory>
        <ModalStory storyName={'Fullscreen'} fullscreen={true}>
          <p>A fullscreen modal</p>
        </ModalStory>
        <ModalStory storyName={'Large'} size={'lg'}>
          <p>A large modal</p>
        </ModalStory>
        <ModalStory storyName={'XLarge'} size={'xl'}>
          <p>A X-large modal</p>
        </ModalStory>
        <ModalStory storyName={'Small'} size={'sm'}>
          <p>A small modal</p>
        </ModalStory>
      </>
    );
  });
