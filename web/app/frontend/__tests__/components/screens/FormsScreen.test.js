import React from 'react';
import FormsScreen from '../../../src/components/screens/FormsScreen';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import { cleanup, render, waitFor } from '@testing-library/react-native';
import testIds from '../../../src/testIds';
import { Provider as PaperProvider } from 'react-native-paper';
import theme, { NavigationTheme } from '../../../src/constants/theme';

const SingleScreenWrapper = ({ screen }) => {
    const Stack = createStackNavigator();
    return (
        <PaperProvider theme={theme}>
            <NavigationContainer them={NavigationTheme}>
                <Stack.Navigator>
                    <Stack.Screen name={'test'} component={screen}/>
                </Stack.Navigator>
            </NavigationContainer>
        </PaperProvider>
    );
};

describe(FormsScreen.name, () => {
    afterEach(cleanup);

    test.only('renders correctly', async () => {
        const tree = render(<SingleScreenWrapper screen={FormsScreen}/>).toJSON();
        await waitFor(() => expect(tree).toMatchSnapshot());
    });

    test('renders a list of forms', async () => {
        const { findAllByTestId } = render(<SingleScreenWrapper screen={FormsScreen}/>);
        const formListItems = await findAllByTestId(testIds.formListItem, { timeout: 3000 });
        expect(formListItems.length).toBeGreaterThan(0);
    });

    test.todo('lets the user navigate to a form by pressing it');
});
