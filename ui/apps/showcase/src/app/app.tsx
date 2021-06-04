import { ReactComponent as Logo } from './logo.svg';
import nrc from './nrc.svg';
// import { Route, Link } from 'react-router-dom';

import { Accordion, Button, Card, Container, Nav, CardImg, CardText, CardBody, CardTitle, CardHeader } from '@core/ui-toolkit';

export function App() {
  return (
    <Container className="p-4" centerContent>
      <Card>
        <CardImg src={nrc} style={{ width: '700px' }} />
        <CardHeader className="text-center">
          Welcome to the NRC Core showcase!
        </CardHeader>
        <CardBody>
          <CardTitle>UI Toolkit</CardTitle>
          <CardText>
            Our custom component library built with React+TS and Bootstrap 5{' '}
            <a href="assets/shared-ui-toolkit/index.html?path=/story/button--basic">
              Check it out!
            </a>
          </CardText>
          <CardTitle>FormRenderer</CardTitle>
          <CardText>
            A tool that transforms forms schemas provided by the server into
            dynamic forms using the UI toolkit{' '}
            <a href="assets/formrenderer/index.html?path=/story/formrenderer--demo">
              Check it out!
            </a>
          </CardText>
        </CardBody>
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
