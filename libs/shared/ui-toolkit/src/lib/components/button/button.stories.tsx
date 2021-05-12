import { Button, CloseButton } from './button.component';
import { Card, CardBody, CardTitle } from '../card/card.component';
import { Container } from '../container/container.component';

export default {
  title: 'Button',
  decorators: [
    (Story: any) => (
      <Container centerContent>
        <Story />
      </Container>
    ),
  ],
};

export const basic = () => (
  <>
    <Card>
      <CardBody>
        <CardTitle>Buttons</CardTitle>
        <Button kind="primary">Primary</Button>
        <Button kind="secondary">Secondary</Button>
        <Button kind="danger">Danger</Button>
        <Button kind="success">Success</Button>
        <Button kind="warning">Warning</Button>
        <Button kind="info">Info</Button>
        <Button kind="dark">Dark</Button>
        <Button kind="light">Light</Button>
        <Button kind="link">Link</Button>
      </CardBody>
    </Card>
    <Card>
      <CardBody>
        <CardTitle>Close Buttons</CardTitle>
        <CloseButton size="sm" />
        <CloseButton />
        <CloseButton size="lg" />
      </CardBody>
    </Card>
  </>
);

export const close = () => (
  <>
    <CloseButton size="sm" />
    <CloseButton />
    <CloseButton size="lg" />
    <div style={{ background: 'black' }}>
      <CloseButton white size="sm" />
      <CloseButton white />
      <CloseButton white size="lg" />
    </div>
  </>
);
