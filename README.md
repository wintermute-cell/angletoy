# TODO
- fix scene template (remove draw funcs, add fixed-/update, impl scene root)
- !!! document input code, document hierachical transforms

- make a sample implementation of the ai module, remove the ai enemy grunt file afterwards.
- clean up utils
- figure out how to do depth sorting / draw sorting
    - figure out how to determine what objects are rendered for what render stage
    - how do render stages interfere with depth sorted rendering?
        - each stage could have it's own depth buffer, and at the end they are composited using the depth buffer

# To document
- Settings -> how do they work, how do you use them, how to set defaults, how to use scripts/sync_settings.py
- The Gem -> how/why are entities stored? how to interface with the Gem? what is fast, what is slow? why are scenes an optional feature?
