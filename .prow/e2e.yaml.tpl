apiVersion: v1
kind: Secret
metadata:
  name: ${ID}
  labels:
    repo_name: "${REPO_NAME}"
    repo_owner: "${REPO_OWNER}"
    job_name: "${JOB_NAME}"
    build_id: "${BUILD_ID}"
    prow_job_id: "${PROW_JOB_ID}"
    pull_base_ref: "${PULL_BASE_REF}"
    pull_base_sha: "${PULL_BASE_SHA}"
data:
  hydraSecret: ${HYDRA_SECRET}
  postgresUsername: ${POSTGRES_USERNAME}
  postgresPassword: ${POSTGRES_PASSWORD}
  mongoUsername: ${MONGO_USERNAME}
  mongoPassword: ${MONGO_PASSWORD}
  postgresDSN: ${POSTGRES_DSN}
  tlsKey: ${TLS_KEY}
  tlsCert: ${TLS_CERT}
---
apiVersion: v1
kind: Service
metadata:
  name: ${ID}-mongo
spec:
  type: ClusterIP
  ports:
    - port: 27017
      protocol: TCP
      name: tcp-mongo
  selector:
    app: mongo
    repo_name: "${REPO_NAME}"
    repo_owner: "${REPO_OWNER}"
    job_name: "${JOB_NAME}"
    build_id: "${BUILD_ID}"
    prow_job_id: "${PROW_JOB_ID}"
    pull_base_ref: "${PULL_BASE_REF}"
    pull_base_sha: "${PULL_BASE_SHA}"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${ID}-mongo
spec:
  selector:
    matchLabels:
      app: mongo
      repo_name: "${REPO_NAME}"
      repo_owner: "${REPO_OWNER}"
      job_name: "${JOB_NAME}"
      build_id: "${BUILD_ID}"
      prow_job_id: "${PROW_JOB_ID}"
      pull_base_ref: "${PULL_BASE_REF}"
      pull_base_sha: "${PULL_BASE_SHA}"
  template:
    metadata:
      labels:
        app: mongo
        repo_name: "${REPO_NAME}"
        repo_owner: "${REPO_OWNER}"
        job_name: "${JOB_NAME}"
        build_id: "${BUILD_ID}"
        prow_job_id: "${PROW_JOB_ID}"
        pull_base_ref: "${PULL_BASE_REF}"
        pull_base_sha: "${PULL_BASE_SHA}"
    spec:
      containers:
        - name: mongo
          image: mongo
          ports:
            - containerPort: 27017
          env:
            - name: MONGO_INITDB_ROOT_PASSWORD
              valueFrom:
                secretKeyRef:
                  key: mongoPassword
                  name: ${ID}
            - name: MONGO_INITDB_ROOT_USERNAME
              valueFrom:
                secretKeyRef:
                  key: mongoUsername
                  name: ${ID}
---
apiVersion: v1
kind: Service
metadata:
  name: ${ID}-redis
spec:
  type: ClusterIP
  ports:
    - port: 6379
      protocol: TCP
      name: tcp-redis
  selector:
    app: redis
    repo_name: "${REPO_NAME}"
    repo_owner: "${REPO_OWNER}"
    job_name: "${JOB_NAME}"
    build_id: "${BUILD_ID}"
    prow_job_id: "${PROW_JOB_ID}"
    pull_base_ref: "${PULL_BASE_REF}"
    pull_base_sha: "${PULL_BASE_SHA}"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${ID}-redis
spec:
  selector:
    matchLabels:
      app: redis
      repo_name: "${REPO_NAME}"
      repo_owner: "${REPO_OWNER}"
      job_name: "${JOB_NAME}"
      build_id: "${BUILD_ID}"
      prow_job_id: "${PROW_JOB_ID}"
      pull_base_ref: "${PULL_BASE_REF}"
      pull_base_sha: "${PULL_BASE_SHA}"
  template:
    metadata:
      labels:
        app: redis
        repo_name: "${REPO_NAME}"
        repo_owner: "${REPO_OWNER}"
        job_name: "${JOB_NAME}"
        build_id: "${BUILD_ID}"
        prow_job_id: "${PROW_JOB_ID}"
        pull_base_ref: "${PULL_BASE_REF}"
        pull_base_sha: "${PULL_BASE_SHA}"
    spec:
      containers:
        - name: redis
          image: bitnami/redis
          ports:
            - containerPort: 6379
          env:
            - name: ALLOW_EMPTY_PASSWORD
              value: "yes"
---
apiVersion: v1
kind: Service
metadata:
  name: ${ID}-hydra-db
spec:
  type: ClusterIP
  ports:
    - port: 5432
      protocol: TCP
      name: tcp-hydra-db
  selector:
    app: hydra-db
    repo_name: "${REPO_NAME}"
    repo_owner: "${REPO_OWNER}"
    job_name: "${JOB_NAME}"
    build_id: "${BUILD_ID}"
    prow_job_id: "${PROW_JOB_ID}"
    pull_base_ref: "${PULL_BASE_REF}"
    pull_base_sha: "${PULL_BASE_SHA}"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${ID}-hydra-db
spec:
  selector:
    matchLabels:
      app: hydra-db
      repo_name: "${REPO_NAME}"
      repo_owner: "${REPO_OWNER}"
      job_name: "${JOB_NAME}"
      build_id: "${BUILD_ID}"
      prow_job_id: "${PROW_JOB_ID}"
      pull_base_ref: "${PULL_BASE_REF}"
      pull_base_sha: "${PULL_BASE_SHA}"
  template:
    metadata:
      labels:
        app: hydra-db
        repo_name: "${REPO_NAME}"
        repo_owner: "${REPO_OWNER}"
        job_name: "${JOB_NAME}"
        build_id: "${BUILD_ID}"
        prow_job_id: "${PROW_JOB_ID}"
        pull_base_ref: "${PULL_BASE_REF}"
        pull_base_sha: "${PULL_BASE_SHA}"
    spec:
      containers:
        - name: postgres
          image: postgres
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  key: postgresUsername
                  name: ${ID}
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  key: postgresPassword
                  name: ${ID}
            - name: POSTGRES_DB
              value: hydra
---
apiVersion: batch/v1
kind: Job
metadata:
  name: ${ID}-hydra-migrate
  labels:
    repo_name: "${REPO_NAME}"
    repo_owner: "${REPO_OWNER}"
    job_name: "${JOB_NAME}"
    build_id: "${BUILD_ID}"
    prow_job_id: "${PROW_JOB_ID}"
    pull_base_ref: "${PULL_BASE_REF}"
    pull_base_sha: "${PULL_BASE_SHA}"
spec:
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: hydra-migrate
          image: oryd/hydra
          args:
            - migrate
            - sql
            - --yes
            - --read-from-env
          env:
            - name: DSN
              valueFrom:
                secretKeyRef:
                  name: ${ID}
                  key: postgresDSN
---
apiVersion: v1
kind: Service
metadata:
  name: ${ID}-hydra
spec:
  type: ClusterIP
  ports:
    - port: 4444
      protocol: TCP
      name: tcp-hydra-public
    - port: 4445
      protocol: TCP
      name: tcp-hydra-admin
  selector:
    app: hydra
    repo_name: "${REPO_NAME}"
    repo_owner: "${REPO_OWNER}"
    job_name: "${JOB_NAME}"
    build_id: "${BUILD_ID}"
    prow_job_id: "${PROW_JOB_ID}"
    pull_base_ref: "${PULL_BASE_REF}"
    pull_base_sha: "${PULL_BASE_SHA}"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${ID}-hydra
spec:
  selector:
    matchLabels:
      app: hydra
      repo_name: "${REPO_NAME}"
      repo_owner: "${REPO_OWNER}"
      job_name: "${JOB_NAME}"
      build_id: "${BUILD_ID}"
      prow_job_id: "${PROW_JOB_ID}"
      pull_base_ref: "${PULL_BASE_REF}"
      pull_base_sha: "${PULL_BASE_SHA}"
  template:
    metadata:
      labels:
        app: hydra
        repo_name: "${REPO_NAME}"
        repo_owner: "${REPO_OWNER}"
        job_name: "${JOB_NAME}"
        build_id: "${BUILD_ID}"
        prow_job_id: "${PROW_JOB_ID}"
        pull_base_ref: "${PULL_BASE_REF}"
        pull_base_sha: "${PULL_BASE_SHA}"
    spec:
      volumes:
        - name: config
          secret:
            secretName: ${ID}
      containers:
        - name: hydra
          image: oryd/hydra
          ports:
            - containerPort: 4444
            - containerPort: 4445
          args:
            - serve
            - all
          volumeMounts:
            - mountPath: /certs/cert.pem
              name: config
              subPath: tlsCert
            - mountPath: /certs/key.pem
              name: config
              subPath: tlsKey
          env:
            - name: DSN
              valueFrom:
                secretKeyRef:
                  name: ${ID}
                  key: postgresDSN
            - name: SECRETS_SYSTEM
              valueFrom:
                secretKeyRef:
                  key: hydraSecret
                  name: ${ID}
            - name: URLS_SELF_ISSUER
              value: https://localhost:4444
            - name: URLS_CONSENT
              value: https://localhost:9000/auth/consent
            - name: URLS_LOGIN
              value: https://localhost:9000/auth/login
            - name: URLS_LOGOUT
              value: https://localhost:9000/auth/logout
            - name: LOG_LEAK_SENSITIVE_VALUES
              value: "true"
            - name: SERVE_TLS_KEY_PATH
              value: /certs/key.pem
            - name: SERVE_TLS_CERT_PATH
              value: /certs/cert.pem
            - name: HTTPS_TLS_KEY_PATH
              value: /certs/key.pem
            - name: HTTPS_TLS_CERT_PATH
              value: /certs/cert.pem
