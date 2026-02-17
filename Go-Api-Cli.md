# GO API Testing TUI/CLI App

Test API routes in the command line tui. The user can create a project. The project will have a base URL, and then the user can add routes to the project. The user can run all routes or run a specific route. They can also output results to a JSON file or to the terminal.

## name ideas

- goapi
- gorest

## tech stack

- Go
- Cobra
- BubbleTea
- SQLite?

## User Stories

- I want to be able to create a new project
- I want to be able to add routes to a project.
- I want to be able to set up authorization for the project: JWT, specifically.
- I want to be able to test all routes, seeing the responses in the terminal or in a json file.
- I want to be able to test a specific routes, seeing the response in the terminal and/or save to a json file.
- I want to be able to operate commands from the command line or the TUI.
- I want to be able to update routes
- I want to be able to variables in my routes. Specifically, if there are multiple resources I can enter the id for that one specifically without having to type in the whole route.
