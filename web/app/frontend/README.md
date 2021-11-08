# How to run this (for now)

from the repo root directory:

`make build`

this will bootstrap everything.

Then, with the server running in another terminal (`make serve`) run:

`make frontend`


You can also start the frontend directly from its directory with yarn:

`yarn start` or `yarn web`/`yarn android` to launch the app directly.

In the metro bundler page that automatically opens after doing `make frontend` (or the alternative commands), 
chose the `Local` option in the bottom left corner above the QR code. 
This ensures that the frontend will be able to communicate with the login server. 
