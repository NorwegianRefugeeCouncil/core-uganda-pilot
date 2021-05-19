import { Button } from '../button/button.component';
import { Card } from '../card/card.component';
import { Form } from './form.component';

export default {
  title: 'Form',
  component: Form,
  decorators: [
    (Story: any) => (
      <Card>
        <Card.Body>
          <Story />
        </Card.Body>
      </Card>
    ),
  ],
};

export const Basic = () => (
  <>
    <Card.Title>Basic form example</Card.Title>
    <Form>
      <Form.Group controlId="exampleEmail" className="mb-3">
        <Form.Label>Email address</Form.Label>
        <Form.Control type="email" />
        <Form.Text>We'll never share your email with anyone else.</Form.Text>
      </Form.Group>
      <Form.Group controlId="examplePassword" className="mb-3">
        <Form.Label>Password</Form.Label>
        <Form.Control type="password" />
        <Form.Text>
          Your password must be 8-20 characters long, contain letters and
          numbers, and must not contain spaces, special characters, or emoji.{' '}
        </Form.Text>
      </Form.Group>
      <Form.Check
        id="exampleCheckbox"
        label="Keep me signed in"
        className="mb-3"
      />
      <Button type="submit">Submit</Button>
    </Form>
  </>
);

export const Form_control = () => <></>;
