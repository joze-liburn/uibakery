# uibakery

UI Bakery code documentation

## UI Bakery 

UI Bakery is a low-code platform to easily build internal tools and apps. It can
produce responsice UI, connect to various data sources, and can be self-hosted.

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

