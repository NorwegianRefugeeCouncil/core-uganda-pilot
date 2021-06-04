import * as React from 'react';
import { Link, Route, Switch, useRouteMatch } from 'react-router-dom';
import { UgandaAssessmentComponent } from './uganda-assessment/uganda-assessment.component';

export const ServiceCatalogComponent: React.FC = props => {

  let { path, url } = useRouteMatch();

  return <Switch>
    <Route path={`${url}/uganda-assessment`}>
      <UgandaAssessmentComponent />
    </Route>
    <Route path={`${path}`}>
      <div className='list-group list-group-flush'>
        <div className='list-group-item'>
          <Link to={`${path}/uganda-assessment`}>
            Uganda Assessment
          </Link>
        </div>
        <div className='list-group-item'>
          ICLA Legal Council
        </div>
        <div className='list-group-item'>
          Shelter Something
        </div>
        <div className='list-group-item'>
          Protection
        </div>
      </div>
    </Route>
  </Switch>;

};
