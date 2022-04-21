type Route = {
  name: string;
  title: string;
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
    title: 'Recipient List',
  },
  recipientsList: {
    name: 'RecipientList',
    title: 'Recipient List',
  },
  recipientsRegistration: {
    name: 'RecipientRegistration',
    title: 'Recipient Registration',
  },
  recipientsProfile: {
    name: 'RecipientProfile',
    title: 'Recipient Profile',
  },
};
