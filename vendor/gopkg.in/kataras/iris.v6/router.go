package iris

import (
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/kataras/go-errors"
	"github.com/kataras/go-fs"
)

const (
	// MethodGet "GET"
	MethodGet = "GET"
	// MethodPost "POST"
	MethodPost = "POST"
	// MethodPut "PUT"
	MethodPut = "PUT"
	// MethodDelete "DELETE"
	MethodDelete = "DELETE"
	// MethodConnect "CONNECT"
	MethodConnect = "CONNECT"
	// MethodHead "HEAD"
	MethodHead = "HEAD"
	// MethodPatch "PATCH"
	MethodPatch = "PATCH"
	// MethodOptions "OPTIONS"
	MethodOptions = "OPTIONS"
	// MethodTrace "TRACE"
	MethodTrace = "TRACE"
	// MethodNone is a Virtual method
	// to store the "offline" routes
	MethodNone = "NONE"
)

var (
	// AllMethods contains all the http valid methods:
	// "GET", "POST", "PUT", "DELETE", "CONNECT", "HEAD", "PATCH", "OPTIONS", "TRACE"
	AllMethods = [...]string{
		MethodGet,
		MethodPost,
		MethodPut,
		MethodDelete,
		MethodConnect,
		MethodHead,
		MethodPatch,
		MethodOptions,
		MethodTrace,
	}
)

const (
	// subdomainIndicator where './' exists in a registered path then it contains subdomain
	subdomainIndicator = "./"
	// DynamicSubdomainIndicator where a registered path starts with '*.' then it contains a dynamic subdomain, if subdomain == "*." then its dynamic
	DynamicSubdomainIndicator = "*."
	// slashByte is just a byte of '/' rune/char
	slashByte = byte('/')
	// slash is just a string of "/"
	slash = "/"
)

var errRouterIsMissing = errors.New(
	`
fatal error, router is missing!
Please .Adapt one of the available routers inside 'kataras/iris/adaptors'.
By-default Iris supports two routers, httprouter and gorillamux.
Edit your main .go source file to adapt one of these routers and restart your app.
	i.e: lines (<---) were missing.
	----------------------------HTTPROUTER----------------------------------
	import (
		"gopkg.in/kataras/iris.v6"
		"gopkg.in/kataras/iris.v6/adaptors/httprouter" // <--- this line
	)

	func main(){
		app := iris.New()
		// right below the iris.New()
		app.Adapt(httprouter.New()) // <--- and this line were missing.

		// the rest of your source code...
		// ...

		app.Listen("%s")
	}


	----------------------------OR GORILLA MUX-------------------------------

	import (
	    "gopkg.in/kataras/iris.v6"
	    "gopkg.in/kataras/iris.v6/adaptors/gorillamux" // <--- or this line
	)

	func main(){
		app := iris.New()
		// right below the iris.New()
		app.Adapt(gorillamux.New()) // <--- and this line were missing.

		app.Listen("%s")
	}
 `)

// Router the visible api for RESTFUL
type Router struct {

	// Ok I thought it very well
	// these changes are breaking for sure
	// but for the best design I have to risk stability.
	// so the router api it's the router
	// and new feature aka policies will be responsible
	// to build the handler and reverse routing
	// from this repo and errors
	// the global routes registry
	repository *routeRepository
	// the global errors registry
	Errors  *ErrorHandlers
	Context ContextPool
	handler http.Handler

	// per-party middleware
	middleware Middleware
	// per-party routes (useful only for done middleware)
	apiRoutes []*route
	// per-party done middleware
	doneMiddleware Middleware
	// per-party
	relativePath string
}

// Regex takes pairs with the named path (without symbols) following by its expression
// and returns a middleware which will do a pure but effective validation using the regexp package.
//
// Note: '/adaptors/gorillamux' already supports regex path validation.
// It's useful while the developer uses the '/adaptors/httprouter' instead.
func (s *Framework) Regex(pairParamExpr ...string) HandlerFunc {
	srvErr := func(ctx *Context) {
		ctx.EmitError(StatusInternalServerError)
	}

	wp := s.policies.RouterReversionPolicy.WildcardPath
	if wp == nil {
		s.Log(ProdMode, "expr cannot be used when a router policy is missing\n"+errRouterIsMissing.Format(s.Config.VHost).Error())
		return srvErr
	}

	if len(pairParamExpr)%2 != 0 {
		s.Log(ProdMode,
			"regexp expr pre-compile error: the correct format is paramName, expression"+
				"paramName2, expression2. The len should be %2==0")
		return srvErr
	}
	pairs := make(map[string]*regexp.Regexp, len(pairParamExpr)/2)

	for i := 0; i < len(pairParamExpr)-1; i++ {
		expr := pairParamExpr[i+1]
		r, err := regexp.Compile(expr)
		if err != nil {
			s.Log(ProdMode, "expr: regexp failed on: "+expr+". Trace:"+err.Error())
			return srvErr
		}

		pairs[pairParamExpr[i]] = r
		i++
	}

	// return the middleware
	return func(ctx *Context) {
		for k, v := range pairs {
			pathPart := ctx.Param(k)
			if pathPart == "" {
				// take care, the router already
				// does the param validations
				// so if it's empty here it means that
				// the router has label it as optional.
				// so we skip it, and continue to the next.
				continue
			}
			// the improtant thing:
			// if the path part didn't match with the relative exp, then fire status not found.
			if !v.MatchString(pathPart) {
				ctx.EmitError(StatusNotFound)
				return
			}
		}
		// otherwise continue to the next handler...
		ctx.Next()
	}
}

var (
	// errDirectoryFileNotFound returns an error with message: 'Directory or file %s couldn't found. Trace: +error trace'
	errDirectoryFileNotFound = errors.New("Directory or file %s couldn't found. Trace: %s")
)

func (router *Router) build(builder RouterBuilderPolicy) {
	router.handler = builder(router.repository, router.Context)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.handler.ServeHTTP(w, r)
}

// Routes returns the routes information,
// some of them can be changed at runtime some others not
// the result of this RoutesInfo is safe to use at RUNTIME.
func (router *Router) Routes() RoutesInfo {
	return router.repository
}

// UseGlobal registers Handler middleware  to the beginning, prepends them instead of append
//
// Use it when you want to add a global middleware to all parties, to all routes in  all subdomains
// It should be called right before Listen functions
func (router *Router) UseGlobal(handlers ...Handler) {
	router.repository.Visit(func(routeInfo RouteInfo) {
		router.repository.ChangeMiddleware(routeInfo, append(handlers, routeInfo.Middleware()...))
	})

	router.Use(handlers...)
}

// UseGlobalFunc registers HandlerFunc middleware to the beginning, prepends them instead of append
//
// Use it when you want to add a global middleware to all parties, to all routes in  all subdomains
// It should be called right before Listen functions
func (router *Router) UseGlobalFunc(handlersFn ...HandlerFunc) {
	router.UseGlobal(convertToHandlers(handlersFn)...)
}

// Party is just a group joiner of routes which have the same prefix and share same middleware(s) also.
// Party can also be named as 'Join' or 'Node' or 'Group' , Party chosen because it has more fun
func (router *Router) Party(relativePath string, handlersFn ...HandlerFunc) *Router {
	parentPath := router.relativePath
	dot := string(subdomainIndicator[0])
	if len(parentPath) > 0 && parentPath[0] == slashByte && strings.HasSuffix(relativePath, dot) { // if ends with . , example: admin., it's subdomain->
		parentPath = parentPath[1:] // remove first slash
	}

	fullpath := parentPath + relativePath
	middleware := convertToHandlers(handlersFn)
	// append the parent's +child's handlers
	middleware = joinMiddleware(router.middleware, middleware)

	return &Router{
		repository:     router.repository,
		Errors:         router.Errors,
		Context:        router.Context,
		handler:        router.handler, // not-needed
		doneMiddleware: router.doneMiddleware,
		apiRoutes:      make([]*route, 0),
		middleware:     middleware,
		relativePath:   fullpath,
	}
}

// Use registers Handler middleware
// returns itself
func (router *Router) Use(handlers ...Handler) *Router {
	router.middleware = append(router.middleware, handlers...)
	return router
}

// UseFunc registers HandlerFunc middleware
// returns itself
func (router *Router) UseFunc(handlersFn ...HandlerFunc) *Router {
	return router.Use(convertToHandlers(handlersFn)...)
}

// Done registers Handler 'middleware' the only difference from .Use is that it
// should be used BEFORE any party route registered or AFTER ALL party's routes have been registered.
//
// returns itself
func (router *Router) Done(handlers ...Handler) *Router {
	if len(router.apiRoutes) > 0 { // register these middleware on previous-party-defined routes, it called after the party's route methods (Handle/HandleFunc/Get/Post/Put/Delete/...)
		for i, n := 0, len(router.apiRoutes); i < n; i++ {
			router.apiRoutes[i].middleware = append(router.apiRoutes[i].middleware, handlers...)
		}
	} else {
		// register them on the doneMiddleware, which will be used on Handle to append these middlweare as the last handler(s)
		router.doneMiddleware = append(router.doneMiddleware, handlers...)
	}

	return router
}

// DoneFunc registers HandlerFunc 'middleware' the only difference from .Use is that it
// should be used BEFORE any party route registered or AFTER ALL party's routes have been registered.
//
// returns itself
func (router *Router) DoneFunc(handlersFn ...HandlerFunc) *Router {
	return router.Done(convertToHandlers(handlersFn)...)
}

// Handle registers a route to the server's router
// if empty method is passed then registers handler(s) for all methods, same as .Any, but returns nil as result
func (router *Router) Handle(method string, registeredPath string, handlers ...Handler) RouteInfo {
	if method == "" { // then use like it was .Any
		for _, k := range AllMethods {
			router.Handle(k, registeredPath, handlers...)
		}
		return nil
	}

	fullpath := router.relativePath + registeredPath // for now, keep the last "/" if any,  "/xyz/"

	middleware := joinMiddleware(router.middleware, handlers)

	// here we separate the subdomain and relative path
	subdomain := ""
	path := fullpath

	if dotWSlashIdx := strings.Index(path, subdomainIndicator); dotWSlashIdx > 0 {
		subdomain = fullpath[0 : dotWSlashIdx+1] // admin.
		path = fullpath[dotWSlashIdx+1:]         // /
	}

	// we splitted the path and subdomain parts so we're ready to check only the path,
	// otherwise we will had problems with subdomains
	// if the user wants beta:= iris.Default.Party("/beta"); beta.Get("/") to be registered as
	//: /beta/ then should disable the path correction OR register it like: beta.Get("//")
	// this is only for the party's roots in order to have expected paths,
	// as we do with iris.Default.Get("/") which is localhost:8080 as RFC points, not localhost:8080/
	///TODO: 31 Jan 2017 -> It does nothing I don't know why I code it but any way' I think it later...
	// if router.mux.correctPath && registeredPath == slash { // check the given relative path
	// 	// remove last "/" if any, "/xyz/"
	// 	if len(path) > 1 { // if it's the root, then keep it*
	// 		if path[len(path)-1] == slashByte {
	// 			// ok we are inside /xyz/
	// 		}
	// 	}
	// }

	path = strings.Replace(path, "//", "/", -1) // fix the path if double //

	if len(router.doneMiddleware) > 0 {
		middleware = append(middleware, router.doneMiddleware...) // register the done middleware, if any
	}
	r := router.repository.register(method, subdomain, path, middleware)

	router.apiRoutes = append(router.apiRoutes, r)
	// should we remove the router.apiRoutes on the .Party (new children party) ?, No, because the user maybe use this party later
	// should we add to the 'inheritance tree' the router.apiRoutes, No, these are for this specific party only, because the user propably, will have unexpected behavior when using Use/UseFunc, Done/DoneFunc
	return r
}

// HandleFunc registers and returns a route with a method string, path string and a handler
// registeredPath is the relative url path
func (router *Router) HandleFunc(method string, registeredPath string, handlersFn ...HandlerFunc) RouteInfo {
	return router.Handle(method, registeredPath, convertToHandlers(handlersFn)...)
}

// None registers an "offline" route
// see context.ExecRoute(routeName),
// iris.Default.None(...) and iris.Default.SetRouteOnline/SetRouteOffline
// For more details look: https://github.com/kataras/iris/issues/585
//
// Example: https://github.com/iris-contrib/examples/tree/master/route_state
func (router *Router) None(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodNone, path, handlersFn...)
}

// Get registers a route for the Get http method
func (router *Router) Get(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodGet, path, handlersFn...)
}

// Post registers a route for the Post http method
func (router *Router) Post(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodPost, path, handlersFn...)
}

// Put registers a route for the Put http method
func (router *Router) Put(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodPut, path, handlersFn...)
}

// Delete registers a route for the Delete http method
func (router *Router) Delete(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodDelete, path, handlersFn...)
}

// Connect registers a route for the Connect http method
func (router *Router) Connect(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodConnect, path, handlersFn...)
}

// Head registers a route for the Head http method
func (router *Router) Head(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodHead, path, handlersFn...)
}

// Options registers a route for the Options http method
func (router *Router) Options(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodOptions, path, handlersFn...)
}

// Patch registers a route for the Patch http method
func (router *Router) Patch(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodPatch, path, handlersFn...)
}

// Trace registers a route for the Trace http method
func (router *Router) Trace(path string, handlersFn ...HandlerFunc) RouteInfo {
	return router.HandleFunc(MethodTrace, path, handlersFn...)
}

// Any registers a route for ALL of the http methods (Get,Post,Put,Head,Patch,Options,Connect,Delete)
func (router *Router) Any(registeredPath string, handlersFn ...HandlerFunc) {
	for _, k := range AllMethods {
		router.HandleFunc(k, registeredPath, handlersFn...)
	}
}

// if / then returns /*wildcard or /something then /something/*wildcard
// if empty then returns /*wildcard too
func validateWildcard(reqPath string, paramName string) string {
	if reqPath[len(reqPath)-1] != slashByte {
		reqPath += slash
	}
	reqPath += "*" + paramName
	return reqPath
}

func (router *Router) registerResourceRoute(reqPath string, h HandlerFunc) RouteInfo {
	router.Head(reqPath, h)
	return router.Get(reqPath, h)
}

// StaticServe serves a directory as web resource
// it's the simpliest form of the Static* functions
// Almost same usage as StaticWeb
// accepts only one required parameter which is the systemPath ( the same path will be used to register the GET&HEAD routes)
// if second parameter is empty, otherwise the requestPath is the second parameter
// it uses gzip compression (compression on each request, no file cache)
func (router *Router) StaticServe(systemPath string, requestPath ...string) RouteInfo {
	var reqPath string

	if len(requestPath) == 0 {
		reqPath = strings.Replace(systemPath, fs.PathSeparator, slash, -1) // replaces any \ to /
		reqPath = strings.Replace(reqPath, "//", slash, -1)                // for any case, replaces // to /
		reqPath = strings.Replace(reqPath, ".", "", -1)                    // replace any dots (./mypath -> /mypath)
	} else {
		reqPath = requestPath[0]
	}

	return router.Get(reqPath+"/*file", func(ctx *Context) {
		filepath := ctx.Param("file")

		spath := strings.Replace(filepath, "/", fs.PathSeparator, -1)
		spath = path.Join(systemPath, spath)

		if !fs.DirectoryExists(spath) {
			ctx.NotFound()
			return
		}

		if err := ctx.ServeFile(spath, true); err != nil {
			ctx.EmitError(StatusInternalServerError)
		}
	})
}

// StaticContent serves bytes, memory cached, on the reqPath
// a good example of this is how the websocket server uses that to auto-register the /iris-ws.js
func (router *Router) StaticContent(reqPath string, cType string, content []byte) RouteInfo {
	modtime := time.Now()
	h := func(ctx *Context) {
		if err := ctx.SetClientCachedBody(StatusOK, content, cType, modtime); err != nil {
			ctx.Log(DevMode, "error while serving []byte via StaticContent: ", err.Error())
		}
	}

	return router.registerResourceRoute(reqPath, h)
}

// StaticEmbedded  used when files are distributed inside the app executable, using go-bindata mostly
// First parameter is the request path, the path which the files in the vdir will be served to, for example "/static"
// Second parameter is the (virtual) directory path, for example "./assets"
// Third parameter is the Asset function
// Forth parameter is the AssetNames function
//
// For more take a look at the
// example: https://github.com/iris-contrib/examples/tree/master/static_files_embedded
func (router *Router) StaticEmbedded(requestPath string, vdir string, assetFn func(name string) ([]byte, error), namesFn func() []string) RouteInfo {
	paramName := "path"
	s := router.Context.Framework()

	requestPath = s.policies.RouterReversionPolicy.WildcardPath(requestPath, paramName)

	if len(vdir) > 0 {
		if vdir[0] == '.' { // first check for .wrong
			vdir = vdir[1:]
		}
		if vdir[0] == '/' || vdir[0] == os.PathSeparator { // second check for /something, (or ./something if we had dot on 0 it will be removed
			vdir = vdir[1:]
		}
	}

	// collect the names we are care for, because not all Asset used here, we need the vdir's assets.
	allNames := namesFn()

	var names []string
	for _, path := range allNames {
		// check if path is the path name we care for
		if !strings.HasPrefix(path, vdir) {
			continue
		}

		path = strings.Replace(path, "\\", "/", -1) // replace system paths with double slashes
		path = strings.Replace(path, "./", "/", -1) // replace ./assets/favicon.ico to /assets/favicon.ico in order to be ready for compare with the reqPath later
		path = path[len(vdir):]                     // set it as the its 'relative' ( we should re-setted it when assetFn will be used)
		names = append(names, path)
	}

	if len(names) == 0 {
		// we don't start the server yet, so:
		s.Log(ProdMode, "error on StaticEmbedded: unable to locate any embedded files located to the (virtual) directory: "+vdir)
	}

	modtime := time.Now()
	h := func(ctx *Context) {
		reqPath := ctx.Param(paramName)
		for _, path := range names {
			// in order to map "/" as "/index.html"
			// as requested here: https://github.com/kataras/iris/issues/633#issuecomment-281691851
			if path == "/index.html" {
				if reqPath[len(reqPath)-1] == slashByte {
					reqPath = "/index.html"
				}
			}

			if path != reqPath {
				continue
			}

			cType := fs.TypeByExtension(path)
			fullpath := vdir + path

			buf, err := assetFn(fullpath)

			if err != nil {
				continue
			}

			if err := ctx.SetClientCachedBody(StatusOK, buf, cType, modtime); err != nil {
				ctx.EmitError(StatusInternalServerError)
				ctx.Log(DevMode, "error while serving via StaticEmbedded: ", err.Error())
			}
			return
		}

		// not found or error
		ctx.EmitError(StatusNotFound)

	}

	return router.registerResourceRoute(requestPath, h)
}

// Favicon serves static favicon
// accepts 2 parameters, second is optional
// favPath (string), declare the system directory path of the __.ico
// requestPath (string), it's the route's path, by default this is the "/favicon.ico" because some browsers tries to get this by default first,
// you can declare your own path if you have more than one favicon (desktop, mobile and so on)
//
// this func will add a route for you which will static serve the /yuorpath/yourfile.ico to the /yourfile.ico (nothing special that you can't handle by yourself)
// Note that you have to call it on every favicon you have to serve automatically (desktop, mobile and so on)
//
// panics on error
func (router *Router) Favicon(favPath string, requestPath ...string) RouteInfo {
	favPath = abs(favPath)

	f, err := os.Open(favPath)
	if err != nil {
		panic(errDirectoryFileNotFound.Format(favPath, err.Error()))
	}

	// ignore error f.Close()
	defer f.Close()
	fi, _ := f.Stat()
	if fi.IsDir() { // if it's dir the try to get the favicon.ico
		fav := path.Join(favPath, "favicon.ico")
		f, err = os.Open(fav)
		if err != nil {
			//we try again with .png
			return router.Favicon(path.Join(favPath, "favicon.png"))
		}
		favPath = fav
		fi, _ = f.Stat()
	}

	cType := fs.TypeByExtension(favPath)
	// copy the bytes here in order to cache and not read the ico on each request.
	cacheFav := make([]byte, fi.Size())
	if _, err = f.Read(cacheFav); err != nil {
		// Here we are before actually run the server.
		// So we could panic but we don't,
		// we just interrupt with a message
		// to the (user-defined) logger.
		router.Context.Framework().Log(DevMode,
			errDirectoryFileNotFound.
				Format(favPath, "favicon: couldn't read the data bytes for file: "+err.Error()).
				Error())
		return nil
	}
	modtime := ""
	h := func(ctx *Context) {
		if modtime == "" {
			modtime = fi.ModTime().UTC().Format(ctx.framework.Config.TimeFormat)
		}
		if t, err := time.Parse(ctx.framework.Config.TimeFormat, ctx.RequestHeader(ifModifiedSince)); err == nil && fi.ModTime().Before(t.Add(StaticCacheDuration)) {

			ctx.ResponseWriter.Header().Del(contentType)
			ctx.ResponseWriter.Header().Del(contentLength)
			ctx.SetStatusCode(StatusNotModified)
			return
		}

		ctx.ResponseWriter.Header().Set(contentType, cType)
		ctx.ResponseWriter.Header().Set(lastModified, modtime)
		ctx.SetStatusCode(StatusOK)
		if _, err := ctx.Write(cacheFav); err != nil {
			ctx.Log(DevMode, "error while trying to serve the favicon: %s", err.Error())
		}
	}

	reqPath := "/favicon" + path.Ext(fi.Name()) //we could use the filename, but because standards is /favicon.ico/.png.
	if len(requestPath) > 0 {
		reqPath = requestPath[0]
	}

	return router.registerResourceRoute(reqPath, h)
}

// StaticHandler returns a new Handler which serves static files
func (router *Router) StaticHandler(reqPath string, systemPath string, showList bool, enableGzip bool, exceptRoutes ...RouteInfo) HandlerFunc {
	// here we separate the path from the subdomain (if any), we care only for the path
	// fixes a bug when serving static files via a subdomain
	fullpath := router.relativePath + reqPath
	path := fullpath
	if dotWSlashIdx := strings.Index(path, subdomainIndicator); dotWSlashIdx > 0 {
		path = fullpath[dotWSlashIdx+1:]
	}

	h := NewStaticHandlerBuilder(systemPath).
		Path(path).
		Listing(showList).
		Gzip(enableGzip).
		Except(exceptRoutes...).
		Build()

	managedStaticHandler := func(ctx *Context) {
		h(ctx)
		prevStatusCode := ctx.ResponseWriter.StatusCode()
		if prevStatusCode >= 400 { // we have an error
			// fire the custom error handler
			router.Errors.Fire(prevStatusCode, ctx)
		}
		// go to the next middleware
		if ctx.Pos < len(ctx.Middleware)-1 {
			ctx.Next()
		}
	}
	return managedStaticHandler
}

// StaticWeb returns a handler that serves HTTP requests
// with the contents of the file system rooted at directory.
//
// first parameter: the route path
// second parameter: the system directory
// third OPTIONAL parameter: the exception routes
//      (= give priority to these routes instead of the static handler)
// for more options look router.StaticHandler.
//
//     router.StaticWeb("/static", "./static")
//
// As a special case, the returned file server redirects any request
// ending in "/index.html" to the same path, without the final
// "index.html".
//
// StaticWeb calls the StaticHandler(reqPath, systemPath, listingDirectories: false, gzip: false ).
func (router *Router) StaticWeb(reqPath string, systemPath string, exceptRoutes ...RouteInfo) RouteInfo {
	h := router.StaticHandler(reqPath, systemPath, false, false, exceptRoutes...)
	paramName := "file"
	routePath := validateWildcard(reqPath, paramName)
	handler := func(ctx *Context) {
		h(ctx)
		if fname := ctx.Param(paramName); fname != "" {
			cType := fs.TypeByExtension(fname)
			if cType != contentBinary && !strings.Contains(cType, "charset") {
				cType += "; charset=" + ctx.framework.Config.Charset
			}

			ctx.SetContentType(cType)
		}

	}

	return router.registerResourceRoute(routePath, handler)
}

// Layout oerrides the parent template layout with a more specific layout for this Party
// returns this Party, to continue as normal
// example:
// my := iris.Default.Party("/my").Layout("layouts/mylayout.html")
// 	{
// 		my.Get("/", func(ctx *iris.Context) {
// 			ctx.MustRender("page1.html", nil)
// 		})
// 	}
//
func (router *Router) Layout(tmplLayoutFile string) *Router {
	router.UseFunc(func(ctx *Context) {
		ctx.Set(TemplateLayoutContextKey, tmplLayoutFile)
		ctx.Next()
	})

	return router
}

// OnError registers a custom http error handler
func (router *Router) OnError(statusCode int, handlerFn HandlerFunc) {
	staticPath := router.Context.Framework().policies.RouterReversionPolicy.StaticPath(router.relativePath)

	if staticPath == "/" {
		router.Errors.Register(statusCode, handlerFn) // register the user-specific error message, as the global error handler, for now.
		return
	}

	// after this, we have more than one error handler for one status code, and that's dangerous some times, but use it for non-globals error catching by your own risk
	// NOTES:
	// subdomains error will not work if same path of a non-subdomain (maybe a TODO for later)
	// errors for parties should be registered from the biggest path length to the smaller.

	// get the previous
	prevErrHandler := router.Errors.GetOrRegister(statusCode)

	func(statusCode int, staticPath string, prevErrHandler Handler, newHandler Handler) { // to separate the logic
		errHandler := HandlerFunc(func(ctx *Context) {
			if strings.HasPrefix(ctx.Path(), staticPath) { // yes the user should use OnError from longest to lower static path's length in order this to work, so we can find another way, like a builder on the end.
				newHandler.Serve(ctx)
				return
			}
			// serve with the user-specific global ("/") pure iris.OnError receiver Handler or the standar handler if OnError called only from inside a no-relative Party.
			prevErrHandler.Serve(ctx)
		})

		router.Errors.Register(statusCode, errHandler)
	}(statusCode, staticPath, prevErrHandler, handlerFn)

}

// EmitError fires a custom http error handler to the client
//
// if no custom error defined with this statuscode, then iris creates one, and once at runtime
func (router *Router) EmitError(statusCode int, ctx *Context) {
	router.Errors.Fire(statusCode, ctx)
}
