import './app.css';
import { TopNav } from '../components/topnav/topnav.component';
import { Route, Switch } from 'react-router-dom';
import FormDefinitions from '../components/formdefinitions';
import FormDefinitionContainer from '../components/formdefinition';

export function App() {
  return <>
    <TopNav />
    <Switch>
      <Route exact path='/'>
        Home
      </Route>
      <Route exact path='/formdefinitions'>
        <FormDefinitions />
      </Route>
      <Route exact path='/formdefinitions/:id'>
        <FormDefinitionContainer />
      </Route>
    </Switch>
  </>;
}

export default App;
