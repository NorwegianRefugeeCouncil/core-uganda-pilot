import { ListGroup, ListGroupItem } from './list-group.component';

export default {
  title: 'List Group',
  component: ListGroup,
};

export const Basic = () => (
  <div className={'container'}>
    <div className={'row'}>
      <div className={'col-12 mb-2'}>
        <p>Standard list group</p>
        <ListGroup>
          <ListGroupItem>Item 1</ListGroupItem>
          <ListGroupItem>Item 2</ListGroupItem>
          <ListGroupItem>Item 3</ListGroupItem>
        </ListGroup>
      </div>
    </div>
    <div className={'row'}>
      <div className={'col-12 mb-2'}>
        <p>Standard list group with active item</p>
        <ListGroup>
          <ListGroupItem active={true}>Item 1</ListGroupItem>
          <ListGroupItem>Item 2</ListGroupItem>
          <ListGroupItem>Item 3</ListGroupItem>
        </ListGroup>
      </div>
    </div>
    <div className={'row'}>
      <div className={'col-12 mb-2'}>
        <p>Standard list group with disabled item</p>
        <ListGroup>
          <ListGroupItem disabled={true}>Item 1</ListGroupItem>
          <ListGroupItem>Item 2</ListGroupItem>
          <ListGroupItem>Item 3</ListGroupItem>
        </ListGroup>
      </div>
    </div>
    <div className={'row'}>
      <div className={'col-12 mb-2'}>
        <p>With action list</p>
        <ListGroup isActionListGroup={true}>
          <ListGroupItem href={'#'} isAction={true}>
            Item 1
          </ListGroupItem>
          <ListGroupItem href={'#'} isAction={true}>
            Item 2
          </ListGroupItem>
          <ListGroupItem href={'#'} isAction={true}>
            Item 3
          </ListGroupItem>
        </ListGroup>
      </div>
    </div>
    <div className={'row'}>
      <div className={'col-12 mb-2'}>
        <p>Flush list group</p>
        <ListGroup flush={true}>
          <ListGroupItem>Item 1</ListGroupItem>
          <ListGroupItem>Item 2</ListGroupItem>
          <ListGroupItem>Item 3</ListGroupItem>
        </ListGroup>
      </div>
    </div>
    <div className={'row'}>
      <div className={'col-12 mb-2'}>
        <p>Numbered list group</p>
        <ListGroup numbered={true}>
          <ListGroupItem>Item 1</ListGroupItem>
          <ListGroupItem>Item 2</ListGroupItem>
          <ListGroupItem>Item 3</ListGroupItem>
        </ListGroup>
      </div>
    </div>
  </div>
);
