This package implements the logic to enable offline usage.

While making API calls, the client thinks he will be interacting with a regular, remote server. This proxy will make
that process transparent.

Since the API requests are predictable, and that we can use discovery, we can infer the types of resources that an API
request refers to. 

We have to handle a few usecases.

1. The client is online

The requests going through this proxy should still be stored. For example, if a user creates a new resource, the proxy
should store that resource in the local store. In this case, the proxy will only store the resources if the API call was
successful.

2. The client is offline

The requests are ending at the proxy. The proxy must fake the responses that would otherwise be returned by the server.
The resources will be stored in the local store, so that the user can still access them.

3. The client was offline, and comes online

The user might have made some changes to the resources. For example, he might have added, deleted or modified resources
while offline. These changes are stored in the local storage, waiting to be reconciled with the remote server.

Simple case:

```
---------+------- time ----------------------------->
         | 
Server   |   V1     V2                               V3
         |   ^      ^                                ^
  User   |   V1     V2    (offline)   V3   (online)  V3     
```

In this scenario, the user while online creates a resource, and updates it too (v1, and v2)
Then, the client goes offline. While offline, he updates the resource again (v3). When online again, he uploads that new
version to the server (v3)

