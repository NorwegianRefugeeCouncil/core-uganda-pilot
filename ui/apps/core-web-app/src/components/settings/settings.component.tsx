import * as React from 'react';
import { ListGroup, ListGroupItem } from '@core/ui-toolkit';
import { Route, Switch } from 'react-router-dom';
import { FormsComponent } from './forms/forms.component';
import { Link, useRouteMatch } from 'react-router-dom';

export const SettingsComponent: React.FC = props => {


  let { path, url } = useRouteMatch();

  return <Switch>
    <Route exact path={path}>

      <h5 className={'p-3'}>Settings</h5>
      <div className='form-check form-switch m-3'>
        <input className='form-check-input' type='checkbox' id='flexSwitchCheckDefault' />
        <label className='form-check-label' htmlFor='flexSwitchCheckDefault'>Offline</label>
        <div id='emailHelp' className='form-text'>Force the application to operate offline</div>
      </div>
      <div className={'border-top'}>
        <button className={'btn btn-primary text-white fw-bold m-3'}>Clear all data</button>
      </div>


      <div className={'d-flex align-items-center px-3 py-2 border-top'}>
        <i className='bi bi-input-cursor-text fs-3 pe-2' />
        <Link to={`${url}/forms`}>Forms</Link>
      </div>

    </Route>
    <Route path={`${path}/forms`}>
      <FormsComponent />
    </Route>
  </Switch>;

};
