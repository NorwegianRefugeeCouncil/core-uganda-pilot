import { ReactComponent as Logo } from './logo.svg';
import nrc from './nrc.svg';
// import { Route, Link } from 'react-router-dom';

import { Accordion, Button, Card, Container, Nav } from '@nrc.no/ui-toolkit';

export function App() {
  return (
    <Container className="p-4" centerContent>
      <Card>
        <Card.Img src={nrc} style={{ width: '700px' }} />
        <Card.Header className="text-center">
          Welcome to the NRC Core showcase!
        </Card.Header>
        <Card.Body>
          <Card.Title>UI Toolkit</Card.Title>
          <Card.Text>
            Our custom component library built with React+TS and Bootstrap 5{' '}
            <a href="assets/shared-ui-toolkit/index.html?path=/story/button--basic">
              Check it out
            </a>
          </Card.Text>
          <Card.Title>FormRenderer</Card.Title>
          <Card.Text>
            A tool that transforms forms schemas provided by the server into
            dynamic forms using the UI toolkit{' '}
            <a href="assets/formrenderer/index.html?path=/story/formrenderer--demo">
              Check it out
            </a>
          </Card.Text>
        </Card.Body>
      </Card>
      <div>
        <Nav role="navigation">
          <Nav.Item>
            <Nav.Link href="assets/shared-ui-toolkit/index.html?path=/story/button--basic">
              UI Toolkit
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link href="assets/formrenderer/index.html?path=/story/formrenderer--demo">
              FormRenderer
            </Nav.Link>
          </Nav.Item>
        </Nav>
      </div>
    </Container>
  );
}

export default App;
