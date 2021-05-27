import * as React from 'react';
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { listFormDefinitions, selectAllFormDefinitions } from '../../reducers/formdefinitions';
import { Container, ListGroup, ListGroupItem } from '@core/ui-toolkit';

const FormDefinitions: React.FC = (props, context) => {

  const dispatch = useDispatch();
  const allFormDefinitions = useSelector(selectAllFormDefinitions);

  useEffect(() => {
    dispatch(listFormDefinitions());
  }, []);

  return (<Container size='fluid'>
    <div className='row'>
      <div className='col'>
        <ListGroup className='mt-2' isActionListGroup={true}>

          {allFormDefinitions.map(f => {
            return <ListGroupItem key={f.metadata.name} href={'/formdefinitions/' + f.metadata.name} isAction={true}>
              {f.metadata.name}
            </ListGroupItem>;
          })}

        </ListGroup>
      </div>
    </div>
  </Container>);

};


export default FormDefinitions;
