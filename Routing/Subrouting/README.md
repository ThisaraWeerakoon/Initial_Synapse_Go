Sample demonstration of ContextVersionStrategy and URLBasedVersionStrategy using the subrouting utlizing net/http's StripPrefix() function.

Note:-
The current VersionStrategy using in MI.
ContextVersionStrategy
    Version should be specified in the Context
    Ex : api context="/bla/{version}/nn" name="api1" version="2.0" version-type="context"
    invoke : http://localhost:8290/bla/2.0/nn/asdf

URLBasedVersionStrategy
    Version should come immediately after the context
    Ex: api context="/bla2/" name="api2" version="1.1" version-type="url"
    Invoke : http://localhost:8290/bla2/1.1/asdf