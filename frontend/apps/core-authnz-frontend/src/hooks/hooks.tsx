import { useCallback, useEffect, useMemo, useState } from 'react';
import classNames from 'classnames';
import { Path, UseFormReturn } from 'react-hook-form';

import Client from '../client/client';
import { IdentityProvider, Organization } from '../types/types';
import { client } from '../app/client';

export function useApiClient(): Client {
  return useMemo(() => client, []);
}

export function useOrganization(
  organizationId?: string,
): Organization | undefined {
  const apiClient = useApiClient();
  const [organization, setOrganization] = useState<Organization>();
  useEffect(() => {
    if (!organizationId) {
      return;
    }
    apiClient.getOrganization({ id: organizationId }).then((resp) => {
      if (resp.response) {
        setOrganization(resp.response);
      }
    });
  }, [apiClient, organizationId]);
  return organization;
}

export function useIdentityProviders(
  organizationId?: string,
): IdentityProvider[] {
  const apiClient = useApiClient();
  const [idps, setIdps] = useState<IdentityProvider[]>([]);
  useEffect(() => {
    if (!organizationId) {
      return;
    }
    apiClient.listIdentityProviders({ organizationId }).then((resp) => {
      if (resp.response) {
        setIdps(resp.response ? resp.response.items : []);
      }
    });
  }, [apiClient, organizationId]);
  return idps;
}

export function useFormValidation<
  T extends Record<string, unknown> = { [key: string]: any },
>(isNew: boolean, form: UseFormReturn<T>) {
  const {
    formState: { dirtyFields, errors, isSubmitted, touchedFields },
  } = form;

  const fieldClasses = useCallback(
    (field: Path<T>) => {
      const cls = classNames('form-control form-control-darkula');
      const touched = (touchedFields as any)[field];
      const dirty = (dirtyFields as any)[field];
      const hasError = !!(errors as any)[field];

      if (isSubmitted || (isNew && touched) || (isNew && dirty)) {
        if (hasError) {
          return classNames(cls, 'is-invalid');
        }
      }
      return cls;
    },
    [dirtyFields, errors, isNew, isSubmitted, touchedFields],
  );

  const fieldErrors = useCallback(
    (field: keyof T) => {
      const touched = (touchedFields as any)[field];
      const dirty = (dirtyFields as any)[field];
      const hasError = (errors as any)[field];
      const err = (errors as any)[field];
      return (
        hasError &&
        (isSubmitted || (isNew && touched) || (isNew && dirty)) && (
          <div className="invalid-feedback">
            {err?.type === 'required' ? (
              <span>This field is required</span>
            ) : (
              <></>
            )}
            {err?.type === 'pattern' ? <span>Invalid value</span> : <></>}
          </div>
        )
      );
    },
    [touchedFields, dirtyFields, errors, isNew, isSubmitted],
  );

  return { fieldErrors, fieldClasses };
}
