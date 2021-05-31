import { FormDefinition } from './types';
import { validateFormDefinition } from './validation';
import { pathFrom, Required } from '../../../field';

const aValidFormDefinition = (): FormDefinition => {
  return {
    apiVersion: 'core.nrc.no/v1',
    kind: 'FormDefinition',
    metadata: {
      name: 'snips.bla.com'
    },
    spec: {
      group: 'bla.com',
      names: {
        plural: 'snips',
        singular: 'snip',
        kind: 'Snip'
      },
      versions: [
        {
          name: 'v1',
          served: true,
          storage: true,
          schema: {
            formSchema: {
              root: {
                type: 'section',
                children: [
                  {
                    key: 'prop1',
                    type: 'shortText',
                    label: [
                      {
                        locale: 'en',
                        value: 'Property 1'
                      }
                    ]
                  }
                ]
              }
            }
          }
        }
      ]
    }
  };
};

describe('formDefinitionValidation', () => {

  it('should validate', function() {
    const errs = validateFormDefinition(aValidFormDefinition());
    expect(errs).toHaveLength(0);
  });

  it('should detect a missing key', function() {
    const f = aValidFormDefinition();
    f.spec.versions[0].schema.formSchema.root.children[0].key = '';
    const errs = validateFormDefinition(f);
    const path = pathFrom('spec.versions[0].schema.formSchema.root.children[0].key');
    expect(errs).toContainEqual(Required(path, 'key is required'));
  });

  it('should detect a missing type', function() {
    const f = aValidFormDefinition();
    f.spec.versions[0].schema.formSchema.root.children[0].type = '' as any;
    const errs = validateFormDefinition(f);
    const path = pathFrom('spec.versions[0].schema.formSchema.root.children[0].type');
    expect(errs).toContainEqual(Required(path, 'type is required'));
  });

  it('should detect empty translation value', function() {
    const f = aValidFormDefinition();
    f.spec.versions[0].schema.formSchema.root.children[0].label[0].value = '';
    const errs = validateFormDefinition(f);
    const path = pathFrom('spec.versions[0].schema.formSchema.root.children[0].label[0].value');
    expect(errs).toContainEqual(Required(path, 'value is required'));
  });

  it('should detect empty translation locale', function() {
    const f = aValidFormDefinition();
    f.spec.versions[0].schema.formSchema.root.children[0].label[0].locale = '';
    const errs = validateFormDefinition(f);
    const path = pathFrom('spec.versions[0].schema.formSchema.root.children[0].label[0].locale');
    expect(errs).toContainEqual(Required(path, 'locale is required'));
  });

});

