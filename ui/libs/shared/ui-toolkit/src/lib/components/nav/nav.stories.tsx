import { Card } from '../card/card.component';
import { Nav } from './nav.component';

export default {
  title: 'Nav',
  component: Nav
};

export const Basic = () => (
  <>
    <h3>Basic Nav</h3>
    <Card className="mb-4">
      <Card.Body>
        <Nav>
          <Nav.Link active>Active</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link disabled>Disabled</Nav.Link>
        </Nav>
      </Card.Body>
    </Card>
    <Card className="mb-4">
      <Card.Body>
        <Nav className="justify-content-center">
          <Nav.Link active>Active</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link disabled>Disabled</Nav.Link>
        </Nav>
      </Card.Body>
    </Card>
    <Card className="mb-4">
      <Card.Body>
        <Nav className="justify-content-end">
          <Nav.Link active>Active</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link disabled>Disabled</Nav.Link>
        </Nav>
      </Card.Body>
    </Card>
    <h3>Vertical</h3>
    <Card className="mb-4">
      <Card.Body>
        <Nav className="flex-column">
          <Nav.Link active>Active</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link disabled>Disabled</Nav.Link>
        </Nav>
      </Card.Body>
    </Card>
    <h3>Tabs</h3>
    <Card className="mb-4">
      <Card.Body>
        <Nav variant="tabs">
          <Nav.Link active>Active</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link disabled>Disabled</Nav.Link>
        </Nav>
      </Card.Body>
    </Card>
    <h3>Pills</h3>
    <Card className="mb-4">
      <Card.Body>
        <Nav variant="pills">
          <Nav.Link active>Active</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link disabled>Disabled</Nav.Link>
        </Nav>
      </Card.Body>
    </Card>
    <h3>Fill</h3>
    <Card className="mb-4">
      <Card.Body>
        <Nav variant="pills" fill>
          <Nav.Link active>Active</Nav.Link>
          <Nav.Link>Much longer nav link</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link disabled>Disabled</Nav.Link>
        </Nav>
      </Card.Body>
    </Card>
    <h3>Justify</h3>
    <Card className="mb-4">
      <Card.Body>
        <Nav variant="pills" justified>
          <Nav.Link active>Active</Nav.Link>
          <Nav.Link>Much longer nav link</Nav.Link>
          <Nav.Link>Link</Nav.Link>
          <Nav.Link disabled>Disabled</Nav.Link>
        </Nav>
      </Card.Body>
    </Card>
  </>
);
