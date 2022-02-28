import React from 'react';
import { storiesOf } from '@storybook/react-native';
import { Text, VStack } from 'native-base';
import { Accordion } from 'core-design-system';
import { Icon } from 'core-design-system/src';

import CenterView from '../CenterView';

storiesOf('Accordion', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Overview, Variants', () => {
    const sections = [
      {
        header: 'First Header',
        content: <Icon name="save" color="primary.500" />,
      },
      {
        header: 'Second Header',
        content: <Text variant="caption">Lorem ipsum...</Text>,
      },
    ];
    return (
      <VStack space="2">
        {sections.map((s, i) => {
          return (
            <Accordion header={s.header} key={i}>
              {s.content}
            </Accordion>
          );
        })}
      </VStack>
    );
  });
