type Route = {
  name: string;
  headerTitle: string;
  pageTitle: string;
};

type Routes = {
  recipientsRoot: Route;
  recipientsList: Route;
  recipientsRegistration: Route;
  recipientsProfile: Route;
};

export const routes: Routes = {
  recipientsRoot: {
    name: 'Recipients',
    headerTitle: 'Recipient List',
    pageTitle: 'Recipient List',
  },
  recipientsList: {
    name: 'RecipientList',
    pageTitle: 'Recipient List',
    headerTitle: 'Recipient List',
  },
  recipientsRegistration: {
    name: 'RecipientRegistration',
    headerTitle: 'Recipient Registration',
    pageTitle: 'Recipient Registration',
  },
  recipientsProfile: {
    name: 'RecipientProfile',
    headerTitle: 'Recipient Profile',
    pageTitle: 'Recipient Profile',
  },
};
