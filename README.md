# UI Bakery Projects

UI Bakery code documentation

## UI Bakery

UI Bakery is a low-code platform to easily build internal tools and apps. It can
produce responsive UI, connect to various data sources, and can be self-hosted.

## Projects

Currently there are two projects developed with UI Bakery tool:

- Lightburn Portal
- Lightburn Vendor Compatibility Matrix

## Project maintenance

There are problems with maintaining UI Bakery projects.

- The current plan does not provide Git integration
- It is possible to export the project, but the export only covers UI. Backend
  code and data sources are not subjects to export.
- It is only possible to interact with the artifacts through visual editor which
  is not very flexible and only allows one particular workflow.
  - Say that we want to check that the first step to every action conforms to
    internal rules. The only way to do so is to open every action and there
    select the desired step.
- Changes can be applied on the fly unintentionally (say, dragging the mouse
  just a bit too much).
- None of the normal software development tools, techniques and workflows apply.
  - No source control
  - Questionable versioning
  - No refactoring

## Backend code

UI Bakery code is split into two categories:

- [Automations](./automations.md)
- [Actions Library](./actionslib.md)

The [call graph](./callgraph.md) (without calls to external services) shows some
automations and actions being unused (by the backend):

- Automations
  - [actionShopifyGetCompaniesList](./automations.md#actionshopifygetcompanieslist)
  - [utils_getStringChecksum](./automations.md#utils_getstringchecksum)
- Actions
  - [prod_utils_compareTwoObj](./actionslib.md#prod_utils_comparetwoobj)
  - [prod_access_createApiKey](./actionslib.md#prod_access_createapikey)
  - [prod_access_addApiKeyAccess](./actionslib.md#prod_access_addapikeyaccess)
  - [prod_pandaDocs_createResellerDoc](./actionslib.md#prod_pandadocs_createresellerdoc)
  - [prod_pandaDocs_getDocs](./actionslib.md#prod_pandadocs_getdocs)

## Data sources

UI Bakery projects are connected to several [data sources](./datasources.md).
