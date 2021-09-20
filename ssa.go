package ssc

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Component Store SSA is a storage for component types.
// When SSA is called, page's general lifecycle components store is not available (we have dummy page instead).
var csSSA = map[string]reflect.Type{}
var csSSALock = &sync.Mutex{}

// SSAHandlerFactory is a factory for building Server Side Action handler.
// Check documentation for lifecycle details (different comparing to page's).
// Example of usage:
// func ssatemplate(p ssc.Page) *template.Template {
// 	return template.Must(template.New("SSA").Funcs(ssc.Funcs()).ParseGlob("*.html"))
// }
// func ssahandler() http.HandlerFunc {
//     return func(rw http.ResponseWriter, r *http.Request) {
// 	       ssc.SSAHandlerFactory(ssatemplate, map[string]interface{}{
//	           "internal:rw": rw,
//             "internal:r": r,
//         })(rw, r)
//     }
// }
func SSAHandlerFactory(tb TemplateBuilder, context map[string]interface{}) http.HandlerFunc {
	// Init dummy page
	dp := &dummypage{
		TemplateBuilder: tb,
	}
	// Set context
	for k, v := range context {
		SetContext(dp, k, v)
	}
	// Return handler
	return func(rw http.ResponseWriter, r *http.Request) {
		// Async specific state
		var wg sync.WaitGroup
		var err = make(chan error, 1000)
		// Extract component action and name from route
		tokens := strings.Split(r.URL.Path, "/")
		cname := tokens[2]
		aname := tokens[3]
		// Find component type in store
		ctype, found := csSSA[cname]
		// Panic, if not found
		if !found {
			panic("Can't find component. Seems like it's not registered")
		}
		// Create component
		component := reflect.New(ctype).Interface().(Component)
		// Init
		if component, ok := component.(ImplementsNestedInit); ok {
			st := time.Now()
			component.Init(dp)
			et := time.Since(st)
			if BENCH_HANDLERS {
				log.Println("Init time", reflect.TypeOf(component), et)
			}
		}
		// Populate component state
		st := time.Now()
		state, _ := url.QueryUnescape(r.PostFormValue("State"))
		if err := json.Unmarshal([]byte(state), &component); err != nil {
			panic(err)
		}
		et := time.Since(st)
		if BENCH_HANDLERS {
			log.Println("Populate time", reflect.TypeOf(component), et)
		}
		// Extract arguments
		st = time.Now()
		var args []interface{}
		json.Unmarshal([]byte(r.PostFormValue("Args")), &args)
		et = time.Since(st)
		if BENCH_HANDLERS {
			log.Println("Extract args time", reflect.TypeOf(component), et)
		}
		// Call action
		st = time.Now()
		if component, ok := component.(ImplementsActions); ok {
			component.Actions()[aname](args...)
		} else {
			panic("Component not implements Actions, unexpected behavior")
		}
		et = time.Since(st)
		if BENCH_HANDLERS {
			log.Println("Action time", reflect.TypeOf(component), et)
		}
		// If new components registered, trigger async
		st = time.Now()
		subset := 0
		for {
			cslLock.RLock()
			regc := csl[dp][subset:]
			cslLock.RUnlock()
			subset += len(regc)
			if len(regc) == 0 {
				break
			}
			for _, component := range regc {
				if component, ok := component.(ImplementsAsync); ok {
					wg.Add(1)
					go func(wg *sync.WaitGroup, err chan error, c ImplementsAsync) {
						defer wg.Done()
						st := time.Now()
						_err := c.Async()
						et := time.Since(st)
						if BENCH_LOWLEVEL {
							log.Println("Async time", reflect.TypeOf(component), et)
						}
						if _err != nil {
							err <- _err
						}
					}(&wg, err, component)
				}
			}
			wg.Wait()
		}
		et = time.Since(st)
		if BENCH_HANDLERS {
			log.Println("Nested async time", et)
		}
		// Extact flags
		redirected := GetContext(dp, "internal:redirected")
		// Render page
		if redirected == nil {
			// Prepare template
			st = time.Now()
			t := dp.Template()
			t = template.Must(t.Parse(fmt.Sprintf(`{{ template "%s" . }}`, cname)))
			et = time.Since(st)
			if BENCH_HANDLERS {
				log.Println("Template prepare time", reflect.TypeOf(component), et)
			}
			// Render component
			st = time.Now()
			terr := t.Execute(rw, component)
			if terr != nil {
				panic(terr)
			}
			et = time.Since(st)
			if BENCH_HANDLERS {
				log.Println("Executiton time", reflect.TypeOf(component), et)
			}
		}
		// Clear context
		DelContext(dp, "")
	}
}

// SSA client side code
var ssaclient = "<script>" +
	"(()=>{var z=11;function re(e,n){var t=n.attributes,r,a,f,u,b;if(!(n.nodeType===z||e.nodeType===z)){for(var S=t.length-1;S>=0;S--)r=t[S],a=r.name,f=r.namespaceURI,u=r.value,f?(a=r.localName||a,b=e.getAttributeNS(f,a),b!==u&&(r.prefix===\"xmlns\"&&(a=r.name),e.setAttributeNS(f,a,u))):(b=e.getAttribute(a),b!==u&&e.setAttribute(a,u));for(var O=e.attributes,w=O.length-1;w>=0;w--)r=O[w],a=r.name,f=r.namespaceURI,f?(a=r.localName||a,n.hasAttributeNS(f,a)||e.removeAttributeNS(f,a)):n.hasAttribute(a)||e.removeAttribute(a)}}var H,ae=\"http://www.w3.org/1999/xhtml\",v=typeof document==\"undefined\"?void 0:document,ie=!!v&&\"content\"in v.createElement(\"template\"),le=!!v&&v.createRange&&\"createContextualFragment\"in v.createRange();function de(e){var n=v.createElement(\"template\");return n.innerHTML=e,n.content.childNodes[0]}function fe(e){H||(H=v.createRange(),H.selectNode(v.body));var n=H.createContextualFragment(e);return n.childNodes[0]}function se(e){var n=v.createElement(\"body\");return n.innerHTML=e,n.childNodes[0]}function ue(e){return e=e.trim(),ie?de(e):le?fe(e):se(e)}function P(e,n){var t=e.nodeName,r=n.nodeName,a,f;return t===r?!0:(a=t.charCodeAt(0),f=r.charCodeAt(0),a<=90&&f>=97?t===r.toUpperCase():f<=90&&a>=97?r===t.toUpperCase():!1)}function ce(e,n){return!n||n===ae?v.createElement(e):v.createElementNS(n,e)}function ve(e,n){for(var t=e.firstChild;t;){var r=t.nextSibling;n.appendChild(t),t=r}return n}function X(e,n,t){e[t]!==n[t]&&(e[t]=n[t],e[t]?e.setAttribute(t,\"\"):e.removeAttribute(t))}var j={OPTION:function(e,n){var t=e.parentNode;if(t){var r=t.nodeName.toUpperCase();r===\"OPTGROUP\"&&(t=t.parentNode,r=t&&t.nodeName.toUpperCase()),r===\"SELECT\"&&!t.hasAttribute(\"multiple\")&&(e.hasAttribute(\"selected\")&&!n.selected&&(e.setAttribute(\"selected\",\"selected\"),e.removeAttribute(\"selected\")),t.selectedIndex=-1)}X(e,n,\"selected\")},INPUT:function(e,n){X(e,n,\"checked\"),X(e,n,\"disabled\"),e.value!==n.value&&(e.value=n.value),n.hasAttribute(\"value\")||e.removeAttribute(\"value\")},TEXTAREA:function(e,n){var t=n.value;e.value!==t&&(e.value=t);var r=e.firstChild;if(r){var a=r.nodeValue;if(a==t||!t&&a==e.placeholder)return;r.nodeValue=t}},SELECT:function(e,n){if(!n.hasAttribute(\"multiple\")){for(var t=-1,r=0,a=e.firstChild,f,u;a;)if(u=a.nodeName&&a.nodeName.toUpperCase(),u===\"OPTGROUP\")f=a,a=f.firstChild;else{if(u===\"OPTION\"){if(a.hasAttribute(\"selected\")){t=r;break}r++}a=a.nextSibling,!a&&f&&(a=f.nextSibling,f=null)}e.selectedIndex=t}}},N=1,oe=11,k=3,W=8;function A(){}function he(e){if(e)return e.getAttribute&&e.getAttribute(\"id\")||e.id}function ge(e){return function(t,r,a){if(a||(a={}),typeof r==\"string\")if(t.nodeName===\"#document\"||t.nodeName===\"HTML\"||t.nodeName===\"BODY\"){var f=r;r=v.createElement(\"html\"),r.innerHTML=f}else r=ue(r);var u=a.getNodeKey||he,b=a.onBeforeNodeAdded||A,S=a.onNodeAdded||A,O=a.onBeforeElUpdated||A,w=a.onElUpdated||A,Q=a.onBeforeNodeDiscarded||A,y=a.onNodeDiscarded||A,Z=a.onBeforeElChildrenUpdated||A,_=a.childrenOnly===!0,T=Object.create(null),L=[];function x(d){L.push(d)}function $(d,l){if(d.nodeType===N)for(var i=d.firstChild;i;){var s=void 0;l&&(s=u(i))?x(s):(y(i),i.firstChild&&$(i,l)),i=i.nextSibling}}function U(d,l,i){Q(d)!==!1&&(l&&l.removeChild(d),y(d),$(d,i))}function G(d){if(d.nodeType===N||d.nodeType===oe)for(var l=d.firstChild;l;){var i=u(l);i&&(T[i]=l),G(l),l=l.nextSibling}}G(t);function C(d){S(d);for(var l=d.firstChild;l;){var i=l.nextSibling,s=u(l);if(s){var o=T[s];o&&P(l,o)?(l.parentNode.replaceChild(o,l),E(o,l)):C(l)}else C(l);l=i}}function ee(d,l,i){for(;l;){var s=l.nextSibling;(i=u(l))?x(i):U(l,d,!0),l=s}}function E(d,l,i){var s=u(l);s&&delete T[s],!(!i&&(O(d,l)===!1||(e(d,l),w(d),Z(d,l)===!1)))&&(d.nodeName!==\"TEXTAREA\"?te(d,l):j.TEXTAREA(d,l))}function te(d,l){var i=l.firstChild,s=d.firstChild,o,h,m,R,g;e:for(;i;){for(R=i.nextSibling,o=u(i);s;){if(m=s.nextSibling,i.isSameNode&&i.isSameNode(s)){i=R,s=m;continue e}h=u(s);var D=s.nodeType,p=void 0;if(D===i.nodeType&&(D===N?(o?o!==h&&((g=T[o])?m===g?p=!1:(d.insertBefore(g,s),h?x(h):U(s,d,!0),s=g):p=!1):h&&(p=!1),p=p!==!1&&P(s,i),p&&E(s,i)):(D===k||D==W)&&(p=!0,s.nodeValue!==i.nodeValue&&(s.nodeValue=i.nodeValue))),p){i=R,s=m;continue e}h?x(h):U(s,d,!0),s=m}if(o&&(g=T[o])&&P(g,i))d.appendChild(g),E(g,i);else{var I=b(i);I!==!1&&(I&&(i=I),i.actualize&&(i=i.actualize(d.ownerDocument||v)),d.appendChild(i),C(i))}i=R,s=m}ee(d,s,h);var K=j[d.nodeName];K&&K(d,l)}var c=t,M=c.nodeType,J=r.nodeType;if(!_){if(M===N)J===N?P(t,r)||(y(t),c=ve(t,ce(r.nodeName,r.namespaceURI))):c=r;else if(M===k||M===W){if(J===M)return c.nodeValue!==r.nodeValue&&(c.nodeValue=r.nodeValue),c;c=r}}if(c===r)y(t);else{if(r.isSameNode&&r.isSameNode(c))return;if(E(c,r,_),L)for(var V=0,ne=L.length;V<ne;V++){var F=T[L[V]];F&&U(F,F.parentNode,!1)}}return!_&&c!==t&&t.parentNode&&(c.actualize&&(c=c.actualize(t.ownerDocument||v)),t.parentNode.replaceChild(c,t)),c}}var pe=ge(re),Y=pe;function B(e){let n=e.starter;if(e.id){let t=document.getElementById(e.id);if(!t)throw new Error(`Error while locating root: can't find direct with ${e}`);n=t}else{let t=0;for(;;){if(!n.parentElement)throw new Error(`Error while locating root: can't find parent with ${e}`);if(!n.getAttribute(\"state\"))n=n.parentElement;else if(e.depth&&t!=e.depth)n=n.parentElement,t++;else break}}return n}function q(e,n,...t){let r=B({starter:e,depth:n.split(\"\").filter(f=>f===\"$\").length,id:n.includes(\":\")?n.split(\":\")[0]:void 0}),a=new FormData;a.set(\"State\",r.getAttribute(\"state\")||\"{}\"),a.set(\"Args\",JSON.stringify(t)),fetch(`/SSA/${r.getAttribute(\"name\")}/${n}`,{method:\"POST\",body:a}).then(f=>f.headers.get(\"X-Redirect\")?(window.location.href=f.headers.get(\"X-Redirect\"),\"\"):f.text()).then(f=>{f&&Y(r,f)}).catch(f=>{console.log(f)})}function Ae(e,n){let t=B({starter:e,depth:n.split(\"\").filter(a=>a===\"$\").length,id:n.includes(\":\")?n.split(\":\")[0]:void 0});if(!t.getAttribute(\"state\"))throw new Error(\"Bind call error: component state is underfined\");let r=JSON.parse(decodeURIComponent(t.getAttribute(\"state\")));r[n]=e.value,t.setAttribute(\"state\",JSON.stringify(r))}function be(e,n){n.preventDefault();let t=B({starter:e});if(!t.getAttribute(\"state\"))throw new Error(\"Bind call error: component state is underfined\");let r=JSON.parse(decodeURIComponent(t.getAttribute(\"state\"))),a=new FormData(n.target),f=Object.fromEntries(a.entries());return Object.entries(f).forEach(u=>{r[u[0]]=u[1]}),t.setAttribute(\"state\",JSON.stringify(r)),q(t,\"Submit\"),!1}window._LocaleRoot=B;window.Action=q;window.Bind=Ae;window.FormSubmit=be;})();" +
	"</script>"
