# snippetbox

## Third-party routers

The Go (1.23) standard library routing doesn't support the following:
 - Sending custom 404 Not Found  and 405 Method Not Allowed  responses to the user.
 - Using regular expressions in your route patterns or wildcards.
 - Matching multiple HTTP methods in a single route declaration.
 - Automatic support for OPTIONS requests.
 - Routing requests to handlers based on unusual things, like HTTP request headers.

If you need these features in your application, you’ll need to use a third-party router package.
The ones that I recommend are httprouter, chi, flow and gorilla/mux, and you can find a comparison
of them and guidance about which one to use in this  [blog post](https://www.alexedwards.net/blog/which-go-router-should-i-use).

## Tips
 - [Patterns for processing and validating different types if input](https://www.alexedwards.net/blog/validation-snippets-for-go)
