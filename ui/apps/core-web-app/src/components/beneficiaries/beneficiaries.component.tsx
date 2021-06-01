import * as React from 'react';
import { ListGroup, ListGroupItem } from '@core/ui-toolkit';
import { Route, Switch, useRouteMatch, Link, useParams } from 'react-router-dom';
import * as uuid from 'uuid';
import { BeneficiaryProfileComponent } from './profile/beneficiary-profile.component';
import { Params } from '../formdefinition';

export const BeneficiariesComponent: React.FC = props => {

  const names = [
    'Mildred Garza',
    'Ismael Mccarthy',
    'Rudolph Barnett',
    'Allen Stephens',
    'Mindy Kelley',
    'Marguerite Bowers',
    'Alfonso Sharp',
    'Marty Holmes',
    'Matthew Singleton',
    'Stephanie Roberts',
    'Robyn Sharp',
    'Emily Tate',
    'Alberta Wilkerson',
    'Ginger Harmon',
    'Martin Peterson',
    'Derek Lawrence',
    'Carroll Mann',
    'Craig Bryan',
    'Laverne Sanders',
    'Grace Carter'
  ];

  let { path, url } = useRouteMatch();

  const { beneficiaryId } = useParams<{beneficiaryId: string}>();

  return (

    <Switch>
      <Route exact path={path}>
        <div className={'bg-white'}>
          <div style={{ marginBottom: '6rem' }}>
            <div className={'input-group w-100 sticky-top border-bottom'}>
              <span className='input-group-text border-0 rounded-0'>@</span>
              <input type='text' className={'form-input flex-grow-1 ps-1 border-0 rounded-0'} placeholder={'Search'} />
            </div>
            <ListGroup>
              {names.map((n, idx) => {
                let className = 'rounded-0';
                if (idx === 0) {
                  className += ' border-top-0';
                }
                return <ListGroupItem className={className}>
                  <Link to={`${url}/${uuid.v4()}`}>{n}</Link>
                </ListGroupItem>;
              })}
            </ListGroup>
          </div>
        </div>
      </Route>
      <Route path={`${path}/:beneficiaryId`}>
        <BeneficiaryProfileComponent beneficiaryId={beneficiaryId} />
      </Route>
    </Switch>


  );

};
