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
import { IdentificationDocumentComponent } from '../components/identificationdocument/identificationdocument.component';
import { HouseholdComponent } from '../components/household/household.component';
import { ServiceCatalogComponent } from '../components/servicecatalog/service-catalog.component';

export function App() {
  return <>
    <TopNav />
    <div style={{ height: 'calc(100vh - 3.5rem)', maxHeight: 'calc(100vh - 3.5rem)', overflowY: 'auto' }}>
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
        <Route path='/identificationdocuments/:id'>
          <IdentificationDocumentComponent />
        </Route>
        <Route path='/households/:id'>
          <HouseholdComponent />
        </Route>
        <Route path='/settings'>
          <SettingsComponent />
        </Route>
        <Route path='/servicecatalog'>
          <ServiceCatalogComponent />
        </Route>
      </Switch>
    </div>
  </>;
}

export default App;
