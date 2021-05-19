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

export const Form_control = () => (
  <>
    <Card.Title>Input types</Card.Title>
    <Form>
      <Form.Group controlId="exampleEmail" className="mb-3">
        <Form.Label>Example email field</Form.Label>
        <Form.Control type="email" />
      </Form.Group>
      <Form.Group controlId="examplePassword" className="mb-3">
        <Form.Label>Example password fied</Form.Label>
        <Form.Control type="password" />
      </Form.Group>
      <Form.Group controlId="exampleText" className="mb-3">
        <Form.Label>Example text fied</Form.Label>
        <Form.Control type="text" />
      </Form.Group>
      <Form.Group controlId="exampleTextarea" className="mb-3">
        <Form.Label>Example textarea fied</Form.Label>
        <Form.Control type="textarea" />
      </Form.Group>
      <Form.Group controlId="exampleFile" className="mb-3">
        <Form.Label>Example file fied</Form.Label>
        <Form.Control type="file" />
      </Form.Group>
      <Form.Group controlId="exampleColor" className="mb-3">
        <Form.Label>Example color fied</Form.Label>
        <Form.Control type="color" />
      </Form.Group>
    </Form>
    <Card.Title>Sizing</Card.Title>
    <Form.Group controlId="exampleSizes">
      <Form.Control
        type="text"
        displaySize="lg"
        placeholder="Large"
        className="mb-3"
      />
      <Form.Control type="text" placeholder="Default" className="mb-3" />
      <Form.Control
        type="text"
        displaySize="sm"
        placeholder="Small"
        className="mb-3"
      />
    </Form.Group>
    <Card.Title>Disabled</Card.Title>
    <Form.Group controlId="exampleSizes">
      <Form.Control
        type="text"
        placeholder="Disabled input"
        className="mb-3"
        disabled
      />
    </Form.Group>
    <Card.Title>Readonly plain text</Card.Title>
    <Form.Group controlId="exampleEmail" className="mb-3 row">
      <Form.Label className="col-sm-2 col-form-label">Email</Form.Label>
      <div className="col-sm-10">
        <Form.Control
          type="email"
          readOnly
          plaintext
          value="spongebob@bikinibottom.com"
        />
      </div>
    </Form.Group>
    <Form.Group controlId="examplePassword" className="mb-3 row">
      <Form.Label className="col-sm-2 col-form-label">Password</Form.Label>
      <div className="col-sm-10">
        <Form.Control type="password" />
      </div>
    </Form.Group>
    <Card.Title>File input</Card.Title>
    <Form.Group controlId="exampleFile" className="mb-3">
      <Form.Label>Default file input</Form.Label>
      <Form.Control type="file" />
    </Form.Group>
    <Form.Group controlId="exampleFile" className="mb-3">
      <Form.Label>Multiple files input</Form.Label>
      <Form.Control type="file" multiple />
    </Form.Group>
    <Form.Group controlId="exampleFile" className="mb-3">
      <Form.Label>Disabled file input</Form.Label>
      <Form.Control type="file" disabled />
    </Form.Group>
    <Form.Group controlId="exampleFile" className="mb-3">
      <Form.Label>Samall file input</Form.Label>
      <Form.Control type="file" displaySize="sm" />
    </Form.Group>
    <Form.Group controlId="exampleFile" className="mb-3">
      <Form.Label>Large file input</Form.Label>
      <Form.Control type="file" displaySize="lg" />
    </Form.Group>
  </>
);

export const Select = () => (
  <>
    <Card.Title>Default</Card.Title>
    <Form.Group controlId="selectExample">
      <Form.Select>
        <option selected>Open this select menu</option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </Form.Select>
    </Form.Group>
    <Card.Title>Sizing</Card.Title>
    <Form.Group controlId="selectExample">
      <Form.Select displaySize="lg" className="mb-3">
        <option selected>Open this select menu</option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </Form.Select>
      <Form.Select displaySize="sm" className="mb-3">
        <option selected>Open this select menu</option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </Form.Select>
    </Form.Group>
    <Card.Title>Multiple</Card.Title>
    <Form.Group controlId="selectExample">
      <Form.Select multiple className="mb-3">
        <option selected>Open this select menu</option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </Form.Select>
      <p>
        also supports the <code>size</code> attribute:
      </p>
      <Form.Select size={3} className="mb-3">
        <option selected>Open this select menu</option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </Form.Select>
    </Form.Group>
    <Card.Title>Disabled</Card.Title>
    <Form.Group controlId="selectExample">
      <Form.Select disabled className="mb-3">
        <option selected>Open this select menu</option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </Form.Select>
    </Form.Group>
  </>
);

export const Checks_and_radios = () => (
  <>
    <Card.Title>Checks</Card.Title>
    <Form.Check id="exampleCheckbox" label="Default checkbox" />
    <Form.Check id="exampleCheckbox2" label="Checked checkbox" defaultChecked />
    <Form.Check id="exampleCheckbox3" label="Disabled checkbox" disabled />
    <Form.Check
      id="exampleCheckbox4"
      label="Disabled checked checkbox"
      disabled
      defaultChecked
    />
    <Card.Title>Radios</Card.Title>
    <Form.Check
      type="radio"
      name="radioDefault"
      id="exampleRadio"
      label="Default radio"
    />
    <Form.Check
      type="radio"
      name="radioDefault"
      id="exampleRadio2"
      label="Checked radio"
      defaultChecked
    />
    <Form.Check
      type="radio"
      id="exampleRadio3"
      label="Disabled radio"
      disabled
    />
    <Form.Check
      type="radio"
      id="exampleRadio4"
      label="Disabled checked radio"
      disabled
      defaultChecked
    />
    <Card.Title>Switchesa</Card.Title>
    <Form.Check type="switch" id="exampleSwitch" label="Default switch" />
    <Form.Check
      type="switch"
      id="exampleSwitch2"
      label="Checked switch"
      defaultChecked
    />
    <Form.Check
      type="switch"
      id="exampleSwitch3"
      label="Disabled switch"
      disabled
    />
    <Form.Check
      type="switch"
      id="exampleSwitch4"
      label="Disabled checked switch"
      disabled
      defaultChecked
    />

    <Card.Title>Inline</Card.Title>
    <Form.Check inline id="exampleCheckbox" label="1" />
    <Form.Check inline id="exampleCheckbox2" label="2" />
    <Form.Check inline id="exampleCheckbox3" label="3 (disabled)" disabled />
    <br />
    <Form.Check
      inline
      type="radio"
      name="defaultRadio"
      id="exampleCheckbox"
      label="1"
    />
    <Form.Check
      inline
      type="radio"
      name="defaultRadio"
      id="exampleCheckbox2"
      label="2"
    />
    <Form.Check
      inline
      type="radio"
      name="defaultRadio"
      id="exampleCheckbox3"
      label="3 (disabled)"
      disabled
    />
    <Card.Title>No labels</Card.Title>
    <Form.Check id="exampleCheckbox" />
    <Form.Check type="radio" name="defaultRadio" id="exampleCheckbox" />
  </>
);
