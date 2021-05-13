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
        <Button className="m-2" colorTheme="primary">
          Primary
        </Button>
        <Button className="m-2" colorTheme="secondary">
          Secondary
        </Button>
        <Button className="m-2" colorTheme="danger">
          Danger
        </Button>
        <Button className="m-2" colorTheme="success">
          Success
        </Button>
        <Button className="m-2" colorTheme="warning">
          Warning
        </Button>
        <Button className="m-2" colorTheme="info">
          Info
        </Button>
        <Button className="m-2" colorTheme="light">
          Light
        </Button>
        <Button className="m-2" colorTheme="dark">
          Dark
        </Button>
        <Button className="m-2" colorTheme="link">
          Link
        </Button>
      </Card.Body>
    </Card>
    <Card className="mb-2">
      <Card.Body>
        <Card.Title>Outline Style</Card.Title>
        <Button outline className="m-2" colorTheme="primary">
          Primary
        </Button>
        <Button outline className="m-2" colorTheme="secondary">
          Secondary
        </Button>
        <Button outline className="m-2" colorTheme="danger">
          Danger
        </Button>
        <Button outline className="m-2" colorTheme="success">
          Success
        </Button>
        <Button outline className="m-2" colorTheme="warning">
          Warning
        </Button>
        <Button outline className="m-2" colorTheme="info">
          Info
        </Button>
        <Button outline className="m-2" colorTheme="light">
          Light
        </Button>
        <Button outline className="m-2" colorTheme="dark">
          Dark
        </Button>
        <Button outline className="m-2" colorTheme="link">
          Link
        </Button>
      </Card.Body>
    </Card>
    <Card className="mb-2">
      <Card.Body>
        <Card.Title>Sizes</Card.Title>
        <Button className="m-2" colorTheme="primary" size="lg">
          Large button
        </Button>
        <Button className="m-2" colorTheme="secondary" size="lg">
          Large button
        </Button>
      </Card.Body>
      <Card.Body>
        <Button className="m-2" colorTheme="primary" size="sm">
          Small button
        </Button>
        <Button className="m-2" colorTheme="secondary" size="sm">
          Small button
        </Button>
      </Card.Body>
    </Card>
    <Card>
      <Card.Body>
        <Card.Title>Disabled state</Card.Title>
        <Button className="m-2" colorTheme="primary" disabled>
          Disabled button
        </Button>
        <Button className="m-2" colorTheme="secondary" disabled>
          Disabled button
        </Button>
      </Card.Body>
    </Card>
  </>
);
