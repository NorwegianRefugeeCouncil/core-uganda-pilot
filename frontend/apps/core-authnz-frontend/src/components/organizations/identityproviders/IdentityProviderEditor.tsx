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
  claimMappings: {
    version: string;
    mappings: {
      subject: string;
      displayName: string;
      fullName: string;
      email: string;
      emailVerified: string;
    };
  };
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
      'claimMappings.mappings.subject',
      JSON.stringify(data.claimMappings.mappings.subject),
    );
    setValue(
      'claimMappings.mappings.displayName',
      JSON.stringify(data.claimMappings.mappings.displayName),
    );
    setValue(
      'claimMappings.mappings.fullName',
      JSON.stringify(data.claimMappings.mappings.fullName),
    );
    setValue(
      'claimMappings.mappings.email',
      JSON.stringify(data.claimMappings.mappings.email),
    );
    setValue(
      'claimMappings.mappings.emailVerified',
      JSON.stringify(data.claimMappings.mappings.emailVerified),
    );
    setVersion(data.claimMappings.version);
  };

  useEffect(() => {
    if (id) {
      apiClient.getIdentityProvider({ id }).then((resp) => {
        if (resp.response) {
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
          version: newVersion,
          mappings: {
            subject: JSON.parse(args.claimMappings.mappings.subject),
            displayName: JSON.parse(args.claimMappings.mappings.displayName),
            fullName: JSON.parse(args.claimMappings.mappings.fullName),
            email: JSON.parse(args.claimMappings.mappings.email),
            emailVerified: JSON.parse(
              args.claimMappings.mappings.emailVerified,
            ),
          },
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

          <h6 className="text-light">
            Current Claim Mapping Version: {version}
          </h6>

          <div className="form-group mb-2">
            <label className="form-label text-light">
              Subject Claim Mapping
            </label>
            <input
              {...register('claimMappings.mappings.subject', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('claimMappings.mappings.subject'),
              )}
            />
            <label className="form-label text-light">
              DisplayName Claim Mapping
            </label>
            <input
              {...register('claimMappings.mappings.displayName', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('claimMappings.mappings.displayName'),
              )}
            />
            <label className="form-label text-light">
              FullName Claim Mapping
            </label>
            <input
              {...register('claimMappings.mappings.fullName', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('claimMappings.mappings.fullName'),
              )}
            />
            <label className="form-label text-light">Email Claim Mapping</label>
            <input
              {...register('claimMappings.mappings.email', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('claimMappings.mappings.email'),
              )}
            />
            <label className="form-label text-light">
              EmailVerified Claim Mapping
            </label>
            <input
              {...register('claimMappings.mappings.emailVerified', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('claimMappings.mappings.emailVerified'),
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
