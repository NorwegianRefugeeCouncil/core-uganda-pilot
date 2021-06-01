import './app.css';
import { TopNav } from '../components/topnav/topnav.component';
import { Route, Switch } from 'react-router-dom';
import FormDefinitions from '../components/formdefinitions';
import FormDefinitionContainer from '../components/formdefinition';
import React from 'react';
import { StoreView } from '../components/store/store';
import { HomeComponent } from '../components/home/home.component';
import { BeneficiariesComponent } from '../components/beneficiaries/beneficiaries.component';
import { SettingsComponent } from '../components/settings/settings.component';

export function App() {
  return <>
    <TopNav />
    <div style={{ marginBottom: '6rem' }}>
      <Switch>
        <Route exact path='/'>
          <HomeComponent />
        </Route>
        <Route exact path='/formdefinitions'>
          <FormDefinitions />
        </Route>
        <Route exact path='/formdefinitions/:id'>
          <FormDefinitionContainer />
        </Route>
        <Route exact path='/store'>
          <StoreView />
        </Route>
        <Route path='/beneficiaries'>
          <BeneficiariesComponent />
        </Route>
        <Route path='/settings'>
          <SettingsComponent />
        </Route>
      </Switch>
    </div>
  </>;
}

export default App;
