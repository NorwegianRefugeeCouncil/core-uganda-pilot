import Card from '../card/card.component';
import Dropdown from './dropdown.component';

export default {
  title: 'Dropdown',
  component: Dropdown,
  decorators: [
    (Story: any) => (
      <Card className="mb-4">
        <Card.Body>
          <Story />
        </Card.Body>
      </Card>
    ),
  ],
};

export const Demo = () => (
  <>
    <div className="mb-4">
      <h3>Single button</h3>
      <Dropdown label="Dropdown button" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </div>
    <div className="mb-4">
      <h3>Colors</h3>
      <Dropdown label="Primary" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown theme="secondary" label="Secondary" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown theme="success" label="Success" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown theme="info" label="Info" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown theme="warning" label="Warning" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown theme="danger" label="Danger" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </div>
    <div className="mb-4">
      <h3>Split button</h3>
      <Dropdown split label="Split button" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </div>
    <div className="mb-4">
      <h3>Dark variant</h3>
      <Dropdown theme="dark" label="Dropdown button" className="m-2 d-inline">
        <Dropdown.Menu dark>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </div>
  </>
);

export const Menu_Items = () => (
  <>
    <div className="mb-4">
      <h3>Active item</h3>
      <Dropdown label="Active item" className="m-2">
        <Dropdown.Menu>
          <Dropdown.Item>Regular link</Dropdown.Item>
          <Dropdown.Item active>Active link</Dropdown.Item>
          <Dropdown.Item>Regular link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </div>
    <div className="mb-4">
      <h3>Disabled item</h3>
      <Dropdown label="Disabled item" className="m-2">
        <Dropdown.Menu>
          <Dropdown.Item>Regular link</Dropdown.Item>
          <Dropdown.Item disabled>Disabled link</Dropdown.Item>
          <Dropdown.Item>Regular link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </div>
    <div className="mb-4">
      <h3>Alignment</h3>
      <Dropdown label="Left-aligned" className="m-2 d-inline">
        <Dropdown.Menu alignStart>
          <Dropdown.Item>Regular link</Dropdown.Item>
          <Dropdown.Item disabled>Disabled link</Dropdown.Item>
          <Dropdown.Item>Regular link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown label="Right-aligned" className="m-2 d-inline">
        <Dropdown.Menu alignEnd>
          <Dropdown.Item>Regular link</Dropdown.Item>
          <Dropdown.Item disabled>Disabled link</Dropdown.Item>
          <Dropdown.Item>Regular link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </div>
    <div className="mb-4">
      <h3>Directions</h3>
      <Dropdown dropDir="up" label="Dropup" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown dropDir="right" label="Dropright" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown dropDir="left" label="Dropleft" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <p className="mb-2" />
      <Dropdown split dropDir="up" label="Up split" className="m-2 d-inline">
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown
        split
        dropDir="right"
        label="Right split"
        className="m-2 d-inline"
      >
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
      <Dropdown
        split
        dropDir="left"
        label="Left split"
        className="m-2 d-inline"
      >
        <Dropdown.Menu>
          <Dropdown.Item>Action</Dropdown.Item>
          <Dropdown.Item>Another action</Dropdown.Item>
          <Dropdown.Item>Something else</Dropdown.Item>
          <Dropdown.Divider />
          <Dropdown.Item>Separated link</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </div>
  </>
);
