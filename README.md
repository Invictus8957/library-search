thinking just aboutl libby interactions, ideally, would be able to:
1. ask what libraries you have a card for 
1. store those libraries somewhere (config file? sqlite?)
1. query for an entry by title/author across those libraries and aggregate results



things that would be interesting to do:
1. search libraries via libby
1. search libraries for "real" books
1. search kindle unlimited inclusions
1. do all of the above from the same interface to make it seamless

Probably can't "check out" books without getting into questionable security territory.


### Libby Interactions
Libby has a locate api which allows looking up a library. 
It looks like it is probably backed by elasticsearch or similar.
It takes search strings directly on the path.
https://locate.libbyapp.com/autocomplete/40509

After finding the library, there is a magic id string like "lexpublib"
which has to be looked up. There is a separate api for that.
Like this https://thunder.api.overdrive.com/v2/libraries/?websiteIds=571
Note the websiteId is returned as part of the response from lookup of the
library in the first place.

After you have id string, you can query availability on a per-library basis
https://thunder.api.overdrive.com/v2/libraries/lexpublib/media?query=foo

There does not appear to be a multi-query endpoint to search multiple libraries
at the same time.
