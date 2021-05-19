import { CloseButton } from './close-button.component';
import { Card } from '../card/card.component';

export default {
  title: 'Close Button',
  component: CloseButton,
};

export const basic = () => (
  <Card>
    <Card.Body>
      <Card.Title>Basic</Card.Title>
      <CloseButton size="sm" />
      <CloseButton />
      <CloseButton size="lg" />
    </Card.Body>
    <Card.Body>
      <Card.Title>Disabled</Card.Title>
      <CloseButton disabled size="sm" />
      <CloseButton disabled />
      <CloseButton disabled size="lg" />
    </Card.Body>
    <Card.Body className="bg-dark text-light">
      <Card.Title>White variant</Card.Title>
      <CloseButton white size="sm" />
      <CloseButton white />
      <CloseButton white size="lg" />
      <br />
      <CloseButton disabled white size="sm" />
      <CloseButton disabled white />
      <CloseButton disabled white size="lg" />
    </Card.Body>
  </Card>
);
