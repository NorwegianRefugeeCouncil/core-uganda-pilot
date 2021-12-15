# Core Server

Core is composed of multiple backend servers.

```
login           Server that handles user login and identity federation 
forms-api       Handles form definitions and records
authnz-bouncer  Authorizes requests
authnz-api      Manages organizations, identity providers, OAuth2 clients
```

The servers can all be run simultaneously, using a single command

```
core serve all
```

Or ran separately

```
core serve forms-api
core serve authnz-api
core serve authnz-bouncer
core serve login
```

## Configuration

Core must be configured using a configuration file

```
core serve forms-api --config=my-config-file.yaml,my-other-config.yaml
```

### Configuration options

```
serve:
  
  # forms api server options
  forms_api:
  	
  	# listener host
  	host: 0.0.0.0
  	
  	# listener port
  	port: 8888
  	
  	# cors options
  	cors:
  	  
  	  # enable/disable cors
  	  enabled: false
  	  
  	  # cors allowed origins
  	  allowed_origins: []
  	  
  	  # cors allowed methods
  	  allowed_methods: []
  	  
  	  # cors allowed headers
  	  allowed_headers: []
  	  
  	  # cors exposed headers
  	  exposed_headers: []  
  	  
  	  # cors allow credentials	  
  	  allow_credentials: false
  	  
  	  # cors passthrough options requests
  	  options_passthrough: false
  	  
  	  # cors max age
  	  max_age: 0
  	  
  	  # debug cors requests
  	  debug: false
  	  
  	# tls options
  	tls:
  	
  	  # enables/disables TLS
  	  enabled: false
  	  cert:
  	  
  	    # path to the TLS cert
  	    path: /cert.pem  	 
  	      
  	  key:
  	  
  	  	# path to the TLS key
  	    path: /key.pem  	    
  	 
  	
  # authn api server options
  authnz_api: 
  
  	# listener host
  	host: 0.0.0.0
  	
  	# listener port
  	port: 8888
  	
  	# cors options
  	cors:
  	  
  	  # enable/disable cors
  	  enabled: false
  	  
  	  # cors allowed origins
  	  allowed_origins: []
  	  
  	  # cors allowed methods
  	  allowed_methods: []
  	  
  	  # cors allowed headers
  	  allowed_headers: []
  	  
  	  # cors exposed headers
  	  exposed_headers: []  
  	  
  	  # cors allow credentials	  
  	  allow_credentials: false
  	  
  	  # cors passthrough options requests
  	  options_passthrough: false
  	  
  	  # cors max age
  	  max_age: 0
  	  
  	  # debug cors requests
  	  debug: false
  	  
  	# tls options
  	tls:
  	
  	  # enables/disables TLS
  	  enabled: false
  	  cert:
  	  
  	    # path to the TLS cert
  	    path: /cert.pem  	 
  	      
  	  key:
  	  
  	  	# path to the TLS key
  	    path: /key.pem  	     
  	 
  # login server options
  login:
    	
  	# listener host
  	host: 0.0.0.0
  	
  	# listener port
  	port: 8888
  	
  	# cors options
  	cors:
  	  
  	  # enable/disable cors
  	  enabled: false
  	  
  	  # cors allowed origins
  	  allowed_origins: []
  	  
  	  # cors allowed methods
  	  allowed_methods: []
  	  
  	  # cors allowed headers
  	  allowed_headers: []
  	  
  	  # cors exposed headers
  	  exposed_headers: []  
  	  
  	  # cors allow credentials	  
  	  allow_credentials: false
  	  
  	  # cors passthrough options requests
  	  options_passthrough: false
  	  
  	  # cors max age
  	  max_age: 0
  	  
  	  # debug cors requests
  	  debug: false
  	  
  	# tls options
  	tls:
  	
  	  # enables/disables TLS
  	  enabled: false
  	  cert:
  	  
  	    # path to the TLS cert
  	    path: /cert.pem  	 
  	      
  	  key:
  	  
  	  	# path to the TLS key
  	    path: /key.pem  	 	    
  	   	 
  # authnz_bouncer server options
  authnz_bouncer:
 
  	# listener host
  	host: 0.0.0.0
  	
  	# listener port
  	port: 8888
  	
  	# cors options
  	cors:  	
  	
  	  # enable/disable cors
  	  enabled: false
  	  
  	  # cors allowed origins
  	  allowed_origins: []
  	  
  	  # cors allowed methods
  	  allowed_methods: []
  	  
  	  # cors allowed headers
  	  allowed_headers: []
  	  
  	  # cors exposed headers
  	  exposed_headers: []  
  	  
  	  # cors allow credentials	  
  	  allow_credentials: false
  	  
  	  # cors passthrough options requests
  	  options_passthrough: false
  	  
  	  # cors max age
  	  max_age: 0
  	  
  	  # debug cors requests
  	  debug: false
  	  
  	# tls options
  	tls:
  	
  	  # enables/disables TLS
  	  enabled: false
  	  cert:
  	  
  	    # path to the TLS cert
  	    path: /cert.pem  	 
  	      
  	  key:
  	  
  	  	# path to the TLS key
  	    path: /key.pem  	     
  	 
  	

# Database connection string
# Additional query parameters can be passed to configure the database connection
#
# max_open_conns=10 set the sql.DB max open connections
# max_idle_conns=10 set the sql.DB max idle connections
# conn_max_idle_time=10 set the sql.DB connection max idle time (seconds)
# conn_max_lifetime=10 set the sql.DB connection max lifetime (seconds)
#
dsn: postgres://user:password@localhost:5432/dbname?sslmode=disable

# Logging configuration
log:
  level: info
  
hydra:
  
  # Hydra Admin endpoint configuration
  admin:  
    schemes: [http]
    host: hydra-host.com
    base_path: /my/custom/hydra/path
    
  # Hydra public endpoint configuration
  public: 
    schemes: [http]
    host: hydra-host.com
    base_path: /my/custom/hydra/path
```
