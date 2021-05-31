import * as React from 'react';
import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Button, Card, CardBody, CardTitle, ListGroup, ListGroupItem } from '@core/ui-toolkit';
import {
  clearCommitted,
  clearUncommitted,
  commit,
  entitiesSet,
  Entity,
  entityAdded,
  HasObjectMeta,
  push,
  rollback,
  selectCommittedEntities,
  selectCommittedOperations,
  selectEntities,
  selectNextQueuedOperation,
  selectUncommittedEntities,
  selectUncommittedOperations,
  TypeMeta,
  submit, isInCondition, acknowledgeAcceptedOperations, selectOperationBatches
} from '@core/api-client';

type TestEntity = TypeMeta & HasObjectMeta & {
  property: string
}

const anEntity = (name: string): TestEntity => {
  return {
    apiVersion: 'test.nrc.no/v1',
    kind: 'TestEntity',
    metadata: {
      name
    },
    property: ''
  };
};

export const StoreView: React.FC = props => {

  const dispatch = useDispatch();
  const stored = useSelector(selectEntities);
  const committed = useSelector(selectCommittedEntities);
  const uncommitted = useSelector(selectUncommittedEntities);
  const committedOperations = useSelector(selectCommittedOperations);
  const uncommittedOperations = useSelector(selectUncommittedOperations);
  const batches = useSelector(selectOperationBatches);
  const first = useSelector(selectNextQueuedOperation);
  const [online, setOnline] = useState(false);
  const [counter, setCounter] = useState(1);

  useEffect(() => {
    if (!online) {
      console.log('offline');
    } else {
      console.log('online');
      if (first) {
        dispatch(submit());
      }
    }
  }, [online, first]);


  return <div className={'container mt-5'}>
    <div className={'row mt-4 bg-light p-3 shadow'}>
      <div className={'col'}>
        <div className={'d-flex flex-row'}>
          <Button className={'text-white shadow'} onClick={() => {
            dispatch(entityAdded(anEntity('entity ' + counter)));
            setCounter(counter + 1);
          }}>Add entity from server</Button>

          <Button className={'text-white shadow ms-2'} onClick={() => {
            dispatch(push({
              operation: {
                apiVersion: 'offline.nrc.no/v1',
                kind: 'Operation',
                metadata: {},
                spec: {
                  entity: anEntity('entity ' + counter),
                  operationType: 'create'
                }
              }
            }));
            setCounter(counter + 1);
          }}>Add entity on client</Button>

          <Button className={'text-white shadow ms-2'} onClick={() => {
            dispatch(commit());
            setCounter(counter + 1);
          }}>Commit</Button>

          <Button className={'text-white shadow ms-2'} onClick={() => {
            dispatch(rollback());
            setCounter(counter + 1);
          }}>Rollback</Button>

          <Button className={'text-white shadow ms-2'} onClick={() => {
            dispatch(clearCommitted());
            dispatch(entitiesSet(committed));
            setCounter(counter + 1);
          }}>Mark committed as synced</Button>

          <Button className={'text-white shadow ms-2'} onClick={() => {
            dispatch(clearCommitted());
            dispatch(entitiesSet([]));
            dispatch(clearUncommitted());
            dispatch(clearCommitted());
            setCounter(0);
          }}>Clear</Button>

          <Button className={'text-white shadow ms-2'} onClick={() => {
            dispatch(acknowledgeAcceptedOperations());
            setCounter(0);
          }}>Acknowledge</Button>

          <div className='form-check form-switch ms-2'>
            <input className='form-check-input' type='checkbox' id='flexSwitchCheckDefault'
                   checked={online}
                   onChange={() => setOnline(!online)}
            />
            <label className='form-check-label' htmlFor='flexSwitchCheckDefault'>Online</label>
          </div>
        </div>

      </div>
    </div>

    <div className={'row mt-4'}>
      {renderEntities('Synced', stored)}
      {renderEntities('Committed', committed)}
      {renderEntities('Uncommitted', uncommitted)}
      <div className={'col'}>
        <div className={'card shadow'}>
          <div className={'card-body px-0'}>
            <CardTitle className={'px-3 pb-2'}>Operations</CardTitle>
            <ListGroup flush={true} className={'border-top'}>

              {
                batches.map(b => {

                  return <Card className={"m-2"}>
                    <ListGroup flush={true}>

                      {b.map(e => {

                        const isAccepted = isInCondition(e, 'OperationAccepted', 'True');
                        const isRejected = isInCondition(e, 'OperationAccepted', 'False');
                        const isPending = isInCondition(e, 'OperationPending', 'True');

                        let color = 'black';
                        if (isAccepted) {
                          color = 'green';
                        } else if (isRejected) {
                          color = 'red';
                        } else if (isPending) {
                          color = 'blue';
                        }

                        return <ListGroupItem>
                          {isAccepted ? <i className={'bi bi-check'} /> : null}
                          {isRejected ? <i className={'bi bi-x'} /> : null}
                          {isPending ? <i className={'bi bi-three-dots'} /> : null}
                          <span style={{ color }}>
              {e.spec.operationType}
                            -
                            {e.spec.entity.metadata.name}</span>
                        </ListGroupItem>;

                      })}

                    </ListGroup>
                  </Card>;

                })
              }

            </ListGroup>
          </div>
        </div>
      </div>
    </div>

  </div>;
};

const renderEntities = (title: string, entities: Entity[]) => {
  return <div className={'col'}>
    <div className={'card shadow'}>
      <div className={'card-body px-0'}>
        <CardTitle className={'px-3 pb-2'}>{title}</CardTitle>
        <ListGroup flush={true} className={'border-top'}>

          {
            entities.length === 0
              ? <ListGroupItem>No items</ListGroupItem>
              : entities.map(e => <ListGroupItem>
                {e.metadata.name}
              </ListGroupItem>)
          }

        </ListGroup>
      </div>
    </div>
  </div>;
};


/**
 *
 * +-------------------
 * |Beneficiary
 * |
 * |      +-----------
 * |      | Household
 * |      |
 *
 *
 * 1. Create benef. B1  -> add ben. in state
 * 2. Create hous. H    -> add hous. in state
 *   2.1 Select head of household -> select B1
 *   2.2 Link B1 -> H has to be established
 *
 *
 *
 *  key: headOfHousehold
 *  type: Entity
 *  kind: Beneficiary
 *  apiVersion: core.nrc.no/v1
 *  required: false
 *
 *
 *  key: household
 *  type: Entity
 *  kind: Household
 *  apiVersion: core.nrc.no/v1
 *
 *
 */
