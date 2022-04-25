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
    title: 'Beneficiaries',
  },
  recipientsList: {
    name: 'recipientsList',
    title: 'Beneficiaries',
  },
  recipientsRegistration: {
    name: 'recipientsRegistration',
    title: 'Beneficiary Registration',
  },
  recipientsProfile: {
    name: 'recipientsProfile',
    title: 'Beneficiary Profile',
  },
};
