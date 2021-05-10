import { storiesOf } from '@storybook/react';
import React from 'react';
import { FieldBuilder, Tab } from './fieldbuilder.component'
import {
  Card, CardBody
} from '@nrc.no/ui-toolkit'
import { GenericFieldConfig } from '../fieldconfig/fieldconfig.component';

storiesOf('Field Builder', module).add('default', () => (
  <Card>
    <CardBody>
      <FieldBuilder 
        tab={Tab.info}
        fieldConfig={{} as GenericFieldConfig}
        fieldType={undefined}
        name=""
        description=""
      />
    </CardBody>
  </Card>
));
