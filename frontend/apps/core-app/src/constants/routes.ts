type Route = {
  name: string;
  title: string;
  linkLabel?: string;
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
    linkLabel: 'Beneficiaries',
  },
  recipientsRegistration: {
    name: 'recipientsRegistration',
    title: 'Beneficiary Registration',
    linkLabel: 'Register Beneficiary',
  },
  recipientsProfile: {
    name: 'recipientsProfile',
    title: 'Beneficiary Profile',
  },
};
