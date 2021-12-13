# Core Server

Core is composed of multiple backend servers.

```
login    Server that handles user login and identity federation 
public   Handles form definitions and records
bouncer  Authorizes requests
admin    Manages organizations, identity providers, OAuth2 clients
```

The servers can all be run simultaneously, using a single command

```
core serve all
```

Or ran separately

```
core serve public
core serve admin
core serve bouncer
core serve login
```

## Configuration

Core must be configured using a configuration file

```
core serve public --config=my-config-file.yaml,my-other-config.yaml
```

### Configuration options

```
serve:
  
  # public server options
  public:
  	
  	# listener host
  	host: 0.0.0.0
  	
  	# listener port
  	port: 8888
  	
  	# cors options
  	cors:
  	  
  	  enabled: false
  	  allowed_origins: []
  	  allowed_methods: []
  	  allowed_headers: []
  	  exposed_headers: []
  	  allow_credentials: false
  	  options_passthrough: false
  	  max_age: 0
  	  debug: false
  	  
  	# tls options
  	tls:
  	  enabled: false
  	  cert:
  	    path: /cert.pem
  	  key:
  	    path: /key.pem  	    
  	 
  	
  # admin server options
  admin: 
  
  	# listener host
  	host: 0.0.0.0
  	
  	# listener port
  	port: 8888
  	
  	# cors options
  	cors:
  	  
  	  enabled: false
  	  allowed_origins: []
  	  allowed_methods: []
  	  allowed_headers: []
  	  exposed_headers: []
  	  allow_credentials: false
  	  options_passthrough: false
  	  max_age: 0
  	  debug: false
  	  
  	# tls options
  	tls:
  	  enabled: false
  	  cert:
  	    path: /cert.pem
  	  key:
  	    path: /key.pem  	    
  	 
  # login server options
  login:
    	
  	# listener host
  	host: 0.0.0.0
  	
  	# listener port
  	port: 8888
  	
  	# cors options
  	cors:
  	  
  	  enabled: false
  	  allowed_origins: []
  	  allowed_methods: []
  	  allowed_headers: []
  	  exposed_headers: []
  	  allow_credentials: false
  	  options_passthrough: false
  	  max_age: 0
  	  debug: false
  	  
  	# tls options
  	tls:
  	  enabled: false
  	  cert:
  	    path: /cert.pem
  	  key:
  	    path: /key.pem  	    
  	 
  # auth bouncer options
  auth:
 
  	# listener host
  	host: 0.0.0.0
  	
  	# listener port
  	port: 8888
  	
  	# cors options
  	cors:  	
  	  enabled: false
  	  allowed_origins: []
  	  allowed_methods: []
  	  allowed_headers: []
  	  exposed_headers: []
  	  allow_credentials: false
  	  options_passthrough: false
  	  max_age: 0
  	  debug: false
  	  
  	# tls options
  	tls:
  	  enabled: false
  	  cert:
  	    path: /cert.pem
  	  key:
  	    path: /key.pem  	    
  	 
  	

# Database connection string
dsn: host=localhost port=5432 user=core password=postgres= dbname=core sslmode=disable

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
