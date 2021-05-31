export type RESTConfig = {
  baseUri: string;
}

export class RESTClient {

  public constructor(public config: RESTConfig) {
  }

  verb(v: string): Request {
    return NewRequest(this).verb(v);
  }

  get(): Request {
    return this.verb('GET');
  }

  post(): Request {
    return this.verb('POST');
  }

  put(): Request {
    return this.verb('PUT');
  }

  delete(): Request {
    return this.verb('DELETE');
  }

}

export function NewRequest(cli: RESTClient): Request {
  return new Request(cli);
}

export class Request {
  private _verb: string;
  private _group: string;
  private _version: string;
  private _name: string;
  private _params: { [key: string]: string[] };
  private _headers: Headers;
  private _resource: string;
  private _body: any;

  public constructor(private cli: RESTClient) {
  }

  public verb(verb: string): Request {
    this._verb = verb;
    return this;
  }

  public group(group: string): Request {
    this._group = group;
    return this;
  }

  public version(version: string): Request {
    this._version = version;
    return this;
  }

  public resource(resource: string): Request {
    this._resource = resource;
    return this;
  }

  public name(name: string): Request {
    this._name = name;
    return this;
  }

  public param(name: string, value: string): Request {
    if (!this._params) {
      this._params = {};
    }
    if (this._params[name]) {
      this._params[name] = [];
    }
    this._params[name].push(value);
    return this;
  }

  public header(name: string, value: string): Request {
    if (!this._headers) {
      this._headers = new Headers();
    }
    this._headers.set(name, value);
    return this;
  }

  public body(body: any): Request {
    this._body = body;
    return this;
  }

  url(): string {
    let url = '';
    if (this.cli.config.baseUri) {
      url += this.cli.config.baseUri;
    }
    if (this._group) {
      url += '/' + this._group;
    }
    if (this._version) {
      url += '/' + this._version;
    }
    if (this._resource) {
      url += '/' + this._resource;
    }
    if (this._name) {
      url += '/' + this._name;
    }

    if (this._params) {
      const parts: string[] = [];
      for (const paramName in this._params) {
        parts.push(paramName + '=' + this._params[paramName].map(p => encodeURIComponent(p)).join(','));
      }
      url += '?' + parts.join('&');
    }

    return url;
  }

  public do<R>(): Promise<R> {

    const url = this.url();

    if (!this._headers) {
      this._headers = new Headers();
    }
    if (this._verb == 'POST' || this._verb == 'PUT' || (this._verb == 'DELETE' && this._body)) {
      if (!this._headers.get('Content-Type')) {
        this._headers.set('Content-Type', 'application/json');
      }
    }
    if (!this._headers.get('Accept')) {
      this._headers.set('Accept', 'application/json');
    }

    return fetch(url, {
      body: this._body,
      headers: this._headers,
      method: this._verb
    }).then((r) => {
      // Expect a successful status code from the api
      if (r.status < 200 || r.status > 299) {
        throw this.parseError(r.json());
      }
      return r.json() as Promise<R>;
    }).catch(err => {
      throw this.parseError(err);
    });

  }

  parseError(response: any): any {
    let obj: any;
    if (typeof response === 'string') {
      try {
        obj = JSON.parse(response);
        return obj;
      } catch (e) {
        obj = {
          kind: 'Status',
          status: 'Failure',
          message: response
        };
        return obj;
      }
    } else if (typeof response === 'object') {
      if (response?.kind === 'Status' && response?.status === 'Failure') {
        return response;
      }
    }
    return response;
  }
}
