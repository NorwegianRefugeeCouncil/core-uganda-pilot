import { RESTClient } from '../rest';
import { ClientSet } from './index';

describe('client', () => {
  it('should work', (done) => {

    const client = new RESTClient({
      baseUri: 'http://localhost:8000/apis'
    });

    const cs = new ClientSet(client);

    cs.core().v1().formDefinitions().get('bla').subscribe(r => {
      console.log(r);
      done();
    }, err => {
      console.error(err);
      done();
    });

  });
});
