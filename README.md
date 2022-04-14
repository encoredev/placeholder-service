<div align="center">
  <a href="https://encore.dev" alt="encore"><img width="189px" src="https://encore.dev/assets/img/logo.svg"></a>
  <h3><a href="https://encore.dev">Encore – The Backend Development Engine</a></h3>
</div>

# Placeholder Service

This is a placeholder service which Encore will use as an initial deployment image when new compute infrastructure is
provisioned on the cloud. It serves a dual purpose of:

1. Responding with a 200 OK to health checks made by the cloud provider & Encores' platform to verify the compute
   infrastructure is ready and routable from the public internet.
2. Responding with a 404 to all other requests with a human-readable error that the service has not been deployed yet.
   This is to prevent user confusion when they attempt to access the service before it has been deployed.

### Configuration

Configuration can be done from either environmental variables, a `.env` file within the working directory or a `config.yaml`
file in the working directory, `/etc/encore/placeholder-service` or `/$HOME/.encore/placeholder-service`.

See [example.env](./example.env) for a list of the configuration options and the documentation for each.
