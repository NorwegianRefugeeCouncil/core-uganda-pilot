import { TabLink, Tabs, Tab } from './tabs.component';
import { storiesOf } from '@storybook/react';
import { Container } from '../container/container.component';

export default {
  title: 'Tabs',
  decorators: [
    (Story: any) => (
      <Container centerContent>
        <Story />
      </Container>
    ),
  ],
};

export const basic = () => (
  <Tabs>
    <Tab>
      <TabLink>Boogle</TabLink>
    </Tab>
    <Tab>
      <TabLink>Doogle</TabLink>
    </Tab>
    <Tab>
      <TabLink>Poogle</TabLink>
    </Tab>
    <Tab>
      <TabLink>Foogle</TabLink>
    </Tab>
  </Tabs>
);
