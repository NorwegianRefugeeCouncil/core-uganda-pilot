import Button from './button.component';
import Card from '../card/card.component';

export default {
  title: 'Button',
  component: Button,
};

export const basic = () => (
  <>
    <Card className="mb-2">
      <Card.Body>
        <Card.Title>Basic Buttons</Card.Title>
        <Button className="m-2" theme="primary">
          Primary
        </Button>
        <Button className="m-2" theme="secondary">
          Secondary
        </Button>
        <Button className="m-2" theme="danger">
          Danger
        </Button>
        <Button className="m-2" theme="success">
          Success
        </Button>
        <Button className="m-2" theme="warning">
          Warning
        </Button>
        <Button className="m-2" theme="info">
          Info
        </Button>
        <Button className="m-2" theme="light">
          Light
        </Button>
        <Button className="m-2" theme="dark">
          Dark
        </Button>
        <Button className="m-2" theme="link">
          Link
        </Button>
      </Card.Body>
    </Card>
    <Card className="mb-2">
      <Card.Body>
        <Card.Title>Outline Style</Card.Title>
        <Button outline className="m-2" theme="primary">
          Primary
        </Button>
        <Button outline className="m-2" theme="secondary">
          Secondary
        </Button>
        <Button outline className="m-2" theme="danger">
          Danger
        </Button>
        <Button outline className="m-2" theme="success">
          Success
        </Button>
        <Button outline className="m-2" theme="warning">
          Warning
        </Button>
        <Button outline className="m-2" theme="info">
          Info
        </Button>
        <Button outline className="m-2" theme="light">
          Light
        </Button>
        <Button outline className="m-2" theme="dark">
          Dark
        </Button>
        <Button outline className="m-2" theme="link">
          Link
        </Button>
      </Card.Body>
    </Card>
    <Card className="mb-2">
      <Card.Body>
        <Card.Title>Sizes</Card.Title>
        <Button className="m-2" theme="primary" size="lg">
          Large button
        </Button>
        <Button className="m-2" theme="secondary" size="lg">
          Large button
        </Button>
      </Card.Body>
      <Card.Body>
        <Button className="m-2" theme="primary" size="sm">
          Small button
        </Button>
        <Button className="m-2" theme="secondary" size="sm">
          Small button
        </Button>
      </Card.Body>
    </Card>
    <Card>
      <Card.Body>
        <Card.Title>Disabled state</Card.Title>
        <Button className="m-2" theme="primary" disabled>
          Disabled button
        </Button>
        <Button className="m-2" theme="secondary" disabled>
          Disabled button
        </Button>
      </Card.Body>
    </Card>
  </>
);
