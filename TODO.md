
# Dredger ToDo's:

## Fixes
- Symlink error in generation process
- Async functionality💁
- localize generation

### Fixes for generated code
- verify tracing works, fix if not
- check and test monitoring code
- update OPA code
- organized shutdown
- reflect required fields of objects from spec in validation

## Features
- Auto-genrate frontend: When pages exist in spec, set AddFrontend to true
- Convert .templ files to go files after generation
- Init git repo (and create initial commit)
- Post-generation suggestions: "We are done, you can cd ... and start the service just run ..."
- Generate usecase handler functions in usecases/ to fully "hide" application logic (like dredger-rs)
- add authenticating code (OpenID, TOTP)


### Features for generated code
- allow form data
- valitate body against entitiy struct, f.e. with "github.com/go-playground/validator/v10" and custom validator if yaml HAS TO BE supported
