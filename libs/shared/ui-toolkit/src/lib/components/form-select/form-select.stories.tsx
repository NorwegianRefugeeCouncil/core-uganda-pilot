import { FormSelect } from './form-select.component';
import { FormLabel } from '../form-label/label.component';

export default {
  title: 'Form Select',
  component: FormSelect,
};

export const Basic = () => (
  <div className={'container'}>
    <p>Simple select</p>
    <div className={'row'}>
      <div className={'col-12 mb-2'}>
        <FormLabel>Select:</FormLabel>
        <FormSelect>
          <option>One</option>
          <option>Two</option>
          <option>Three</option>
        </FormSelect>
      </div>
    </div>
    <div className={'row'}>
      <p>Select multiple</p>
      <div className={'col-12 mb-2'}>
        <FormLabel>Select:</FormLabel>
        <FormSelect multiple={true}>
          <option>One</option>
          <option>Two</option>
          <option>Three</option>
        </FormSelect>
      </div>
    </div>
    <div className={'row'}>
      <p>Select disabled</p>
      <div className={'col-12 mb-2'}>
        <FormLabel>Select:</FormLabel>
        <FormSelect disabled={true}>
          <option>One</option>
          <option>Two</option>
          <option>Three</option>
        </FormSelect>
      </div>
    </div>
  </div>
);
