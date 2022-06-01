import { FC, useCallback, useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import classNames from 'classnames';

import { useApiClient, useFormValidation } from '../../../hooks/hooks';
import { IdentityProvider, Organization } from '../../../types/types';

type Props = {
  id?: string;
  organization: Organization;
};

type FormData = {
  name: string;
  issuer: string;
  clientId: string;
  clientSecret: string;
  organizationId: string;
  emailDomain: string;
  scopes: string;
  claimMappings: { Version: string; Mappings: any };
};

export const IdentityProviderEditor: FC<Props> = (props) => {
  const { id, organization } = props;

  const isNew = useMemo(() => !id, [id]);

  const apiClient = useApiClient();

  const form = useForm<FormData>({ mode: 'onChange' });

  const [version, setVersion] = useState<string>('0');

  const {
    register,
    handleSubmit,
    setValue,
    formState: { isSubmitting },
  } = form;

  const { fieldErrors, fieldClasses } = useFormValidation(isNew, form);

  const setData = (data: IdentityProvider) => {
    setValue('name', data.name);
    setValue('clientId', data.clientId);
    setValue('organizationId', data.organizationId);
    setValue('issuer', data.domain);
    setValue('emailDomain', data.emailDomain);
    setValue('clientSecret', '');
    setValue('scopes', data.scopes);
    setValue(
      'claimMappings.Mappings',
      JSON.stringify(data.claimMappings.Mappings),
    );
    setVersion(data.claimMappings.Version);
  };

  useEffect(() => {
    if (id) {
      apiClient.getIdentityProvider({ id }).then((resp) => {
        if (resp.response) {
          console.log('RESP', resp.response);
          setData(resp.response);
        }
      });
    }
  }, [apiClient, id, setValue]);

  const onSubmit = useCallback(
    async (args: FormData) => {
      const newVersion = `${parseInt(version, 10) + 1}`;

      const object = {
        id,
        name: args.name,
        clientId: args.clientId,
        clientSecret: args.clientSecret,
        domain: args.issuer,
        organizationId: organization.id,
        emailDomain: args.emailDomain,
        scopes: args.scopes,
        claimMappings: {
          Version: newVersion,
          Mappings: JSON.parse(args.claimMappings.Mappings),
        },
      };
      let resp;
      if (id) {
        resp = await apiClient.updateIdentityProvider({
          object,
        });
      } else {
        resp = await apiClient.createIdentityProvider({
          object,
        });
      }
      return resp.response && setData(resp.response);
    },
    [apiClient, id, organization.id, version],
  );

  return (
    <div className={classNames('card bg-dark border-secondary')}>
      <div className="card-body">
        <form
          className="needs-validation"
          noValidate
          onSubmit={handleSubmit(onSubmit)}
        >
          <div className={classNames('form-group mb-2')}>
            <label className="form-label text-light">Name</label>
            <input
              {...register('name', {
                required: true,
                pattern: /^[a-zA-Z0-9\-_ ]+$/,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('name'),
              )}
            />
            {fieldErrors('name')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Issuer</label>
            <input
              {...register('issuer', {
                required: true,
                pattern: /^https?:\/\/[a-zA-Z0-9.\-_]+(:[0-9]+)?$/,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('issuer'),
              )}
            />
            {fieldErrors('issuer')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Email Domain</label>
            <input
              {...register('emailDomain', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('emailDomain'),
              )}
            />
            {fieldErrors('emailDomain')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Client ID</label>
            <input
              {...register('clientId', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('clientId'),
              )}
            />
            {fieldErrors('clientId')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Client Secret</label>
            <input
              type="password"
              {...register('clientSecret', {
                required: isNew,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('clientSecret'),
              )}
              placeholder={isNew ? '' : '********'}
            />
            {fieldErrors('clientSecret')}
          </div>

          <div className="form-group mb-2">
            <label className="form-label text-light">Scopes</label>
            <input
              {...register('scopes', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('scopes'),
              )}
            />
            {fieldErrors('scopes')}
          </div>

          <div className="form-group mb-2">
            <label className="form-label text-light">
              Claim Mapping (Version: {version})
            </label>
            <textarea
              {...register('claimMappings.Mappings', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('claimMappings.Mappings'),
              )}
            />
          </div>
          {fieldErrors('claimMappings')}

          <button disabled={isSubmitting} className="btn btn-success mt-2">
            {props.id ? 'Update Identity Provider' : 'Create Identity Provider'}
          </button>
        </form>
      </div>
    </div>
  );
};
