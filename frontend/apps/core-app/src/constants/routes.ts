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
    name: 'recipientsRoot',
    title: 'Recipient List',
  },
  recipientsList: {
    name: 'recipientsList',
    title: 'Recipient List',
  },
  recipientsRegistration: {
    name: 'recipientsRegistration',
    title: 'Recipient Registration',
  },
  recipientsProfile: {
    name: 'recipientsProfile',
    title: 'Recipient Profile',
  },
};
