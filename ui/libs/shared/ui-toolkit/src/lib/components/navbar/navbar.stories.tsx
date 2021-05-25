import { Nav } from '../nav/nav.component';
import { Navbar } from './navbar.component';

export default {
  title: 'Navbar',
  component: Navbar
};

export const Basic = () => (
  <Navbar>
    <Navbar.Brand>Navbar</Navbar.Brand>
    <Navbar.Nav>
      <Nav.Link active>Home</Nav.Link>
      <Nav.Link>Beneficiaries</Nav.Link>
      <Nav.Link>Profile</Nav.Link>
    </Navbar.Nav>
    <Navbar.Text>Headline text</Navbar.Text>
  </Navbar>
);
