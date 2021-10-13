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
		if _component, ok := component.(ImplementsInit); ok {
			st := time.Now()
			_component.Init(dp)
			et := time.Since(st)
			if BENCH_HANDLERS {
				log.Println("Init time", reflect.TypeOf(component), et)
			}
		} else if _component, ok := component.(ImplementsInitWithoutPage); ok {
			st := time.Now()
			_component.Init()
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
		if _component, ok := component.(ImplementsActions); ok {
			_component.Actions(dp)[aname](args...)
		} else if _component, ok := component.(ImplementsActionsWithoutPage); ok {
			_component.Actions()[aname](args...)
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
				if _component, ok := component.(ImplementsAsync); ok {
					wg.Add(1)
					go func(wg *sync.WaitGroup, err chan error, c ImplementsAsync, dp Page) {
						defer wg.Done()
						st := time.Now()
						_err := c.Async(dp)
						et := time.Since(st)
						if BENCH_LOWLEVEL {
							log.Println("Async time", reflect.TypeOf(component), et)
						}
						if _err != nil {
							err <- _err
						}
					}(&wg, err, _component, dp)
				} else if _component, ok := component.(ImplementsAsyncWithoutPage); ok {
					wg.Add(1)
					go func(wg *sync.WaitGroup, err chan error, c ImplementsAsyncWithoutPage) {
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
					}(&wg, err, _component)
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
	"(()=>{var z=11;function ne(e,n){var t=n.attributes,r,a,f,s,b;if(!(n.nodeType===z||e.nodeType===z)){for(var S=t.length-1;S>=0;S--)r=t[S],a=r.name,f=r.namespaceURI,s=r.value,f?(a=r.localName||a,b=e.getAttributeNS(f,a),b!==s&&(r.prefix==='xmlns'&&(a=r.name),e.setAttributeNS(f,a,s))):(b=e.getAttribute(a),b!==s&&e.setAttribute(a,s));for(var y=e.attributes,w=y.length-1;w>=0;w--)r=y[w],a=r.name,f=r.namespaceURI,f?(a=r.localName||a,n.hasAttributeNS(f,a)||e.removeAttributeNS(f,a)):n.hasAttribute(a)||e.removeAttribute(a)}}var D,ae='http://www.w3.org/1999/xhtml',o=typeof document=='undefined'?void 0:document,ie=!!o&&'content'in o.createElement('template'),le=!!o&&o.createRange&&'createContextualFragment'in o.createRange();function de(e){var n=o.createElement('template');return n.innerHTML=e,n.content.childNodes[0]}function fe(e){D||(D=o.createRange(),D.selectNode(o.body));var n=D.createContextualFragment(e);return n.childNodes[0]}function ue(e){var n=o.createElement('body');return n.innerHTML=e,n.childNodes[0]}function se(e){return e=e.trim(),ie?de(e):le?fe(e):ue(e)}function P(e,n){var t=e.nodeName,r=n.nodeName,a,f;return t===r?!0:(a=t.charCodeAt(0),f=r.charCodeAt(0),a<=90&&f>=97?t===r.toUpperCase():f<=90&&a>=97?r===t.toUpperCase():!1)}function ce(e,n){return!n||n===ae?o.createElement(e):o.createElementNS(n,e)}function oe(e,n){for(var t=e.firstChild;t;){var r=t.nextSibling;n.appendChild(t),t=r}return n}function $(e,n,t){e[t]!==n[t]&&(e[t]=n[t],e[t]?e.setAttribute(t,''):e.removeAttribute(t))}var j={OPTION:function(e,n){var t=e.parentNode;if(t){var r=t.nodeName.toUpperCase();r==='OPTGROUP'&&(t=t.parentNode,r=t&&t.nodeName.toUpperCase()),r==='SELECT'&&!t.hasAttribute('multiple')&&(e.hasAttribute('selected')&&!n.selected&&(e.setAttribute('selected','selected'),e.removeAttribute('selected')),t.selectedIndex=-1)}$(e,n,'selected')},INPUT:function(e,n){$(e,n,'checked'),$(e,n,'disabled'),e.value!==n.value&&(e.value=n.value),n.hasAttribute('value')||e.removeAttribute('value')},TEXTAREA:function(e,n){var t=n.value;e.value!==t&&(e.value=t);var r=e.firstChild;if(r){var a=r.nodeValue;if(a==t||!t&&a==e.placeholder)return;r.nodeValue=t}},SELECT:function(e,n){if(!n.hasAttribute('multiple')){for(var t=-1,r=0,a=e.firstChild,f,s;a;)if(s=a.nodeName&&a.nodeName.toUpperCase(),s==='OPTGROUP')f=a,a=f.firstChild;else{if(s==='OPTION'){if(a.hasAttribute('selected')){t=r;break}r++}a=a.nextSibling,!a&&f&&(a=f.nextSibling,f=null)}e.selectedIndex=t}}},N=1,ve=11,k=3,W=8;function A(){}function he(e){if(e)return e.getAttribute&&e.getAttribute('id')||e.id}function pe(e){return function(t,r,a){if(a||(a={}),typeof r=='string')if(t.nodeName==='#document'||t.nodeName==='HTML'||t.nodeName==='BODY'){var f=r;r=o.createElement('html'),r.innerHTML=f}else r=se(r);var s=a.getNodeKey||he,b=a.onBeforeNodeAdded||A,S=a.onNodeAdded||A,y=a.onBeforeElUpdated||A,w=a.onElUpdated||A,Q=a.onBeforeNodeDiscarded||A,O=a.onNodeDiscarded||A,Z=a.onBeforeElChildrenUpdated||A,B=a.childrenOnly===!0,T=Object.create(null),L=[];function M(d){L.push(d)}function X(d,l){if(d.nodeType===N)for(var i=d.firstChild;i;){var u=void 0;l&&(u=s(i))?M(u):(O(i),i.firstChild&&X(i,l)),i=i.nextSibling}}function x(d,l,i){Q(d)!==!1&&(l&&l.removeChild(d),O(d),X(d,i))}function G(d){if(d.nodeType===N||d.nodeType===ve)for(var l=d.firstChild;l;){var i=s(l);i&&(T[i]=l),G(l),l=l.nextSibling}}G(t);function C(d){S(d);for(var l=d.firstChild;l;){var i=l.nextSibling,u=s(l);if(u){var v=T[u];v&&P(l,v)?(l.parentNode.replaceChild(v,l),U(v,l)):C(l)}else C(l);l=i}}function ee(d,l,i){for(;l;){var u=l.nextSibling;(i=s(l))?M(i):x(l,d,!0),l=u}}function U(d,l,i){var u=s(l);u&&delete T[u],!(!i&&(y(d,l)===!1||(e(d,l),w(d),Z(d,l)===!1)))&&(d.nodeName!=='TEXTAREA'?te(d,l):j.TEXTAREA(d,l))}function te(d,l){var i=l.firstChild,u=d.firstChild,v,h,m,R,p;e:for(;i;){for(R=i.nextSibling,v=s(i);u;){if(m=u.nextSibling,i.isSameNode&&i.isSameNode(u)){i=R,u=m;continue e}h=s(u);var H=u.nodeType,g=void 0;if(H===i.nodeType&&(H===N?(v?v!==h&&((p=T[v])?m===p?g=!1:(d.insertBefore(p,u),h?M(h):x(u,d,!0),u=p):g=!1):h&&(g=!1),g=g!==!1&&P(u,i),g&&U(u,i)):(H===k||H==W)&&(g=!0,u.nodeValue!==i.nodeValue&&(u.nodeValue=i.nodeValue))),g){i=R,u=m;continue e}h?M(h):x(u,d,!0),u=m}if(v&&(p=T[v])&&P(p,i))d.appendChild(p),U(p,i);else{var I=b(i);I!==!1&&(I&&(i=I),i.actualize&&(i=i.actualize(d.ownerDocument||o)),d.appendChild(i),C(i))}i=R,u=m}ee(d,u,h);var K=j[d.nodeName];K&&K(d,l)}var c=t,E=c.nodeType,J=r.nodeType;if(!B){if(E===N)J===N?P(t,r)||(O(t),c=oe(t,ce(r.nodeName,r.namespaceURI))):c=r;else if(E===k||E===W){if(J===E)return c.nodeValue!==r.nodeValue&&(c.nodeValue=r.nodeValue),c;c=r}}if(c===r)O(t);else{if(r.isSameNode&&r.isSameNode(c))return;if(U(c,r,B),L)for(var V=0,re=L.length;V<re;V++){var F=T[L[V]];F&&x(F,F.parentNode,!1)}}return!B&&c!==t&&t.parentNode&&(c.actualize&&(c=c.actualize(t.ownerDocument||o)),t.parentNode.replaceChild(c,t)),c}}var ge=pe(ne),Y=ge;function _(e){let n=e.starter;if(e.id){let t=document.getElementById(e.id);if(!t)throw new Error(`Error while locating root: can't find direct with ${e}`);n=t}else{let t=0;for(;;){if(!n.parentElement)throw new Error(`Error while locating root: can't find parent with ${e}`);if(!n.getAttribute('state'))n=n.parentElement;else if(e.depth&&t!=e.depth)n=n.parentElement,t++;else break}}return n}function Ae(e){return e.includes(':')&&(e=e.split(':')[1]),e.includes('$')&&(e=e.replaceAll('$','')),e}function q(e,n,...t){let r=_({starter:e,depth:n.split('').filter(f=>f==='$').length,id:n.includes(':')?n.split(':')[0]:void 0}),a=new FormData;a.set('State',r.getAttribute('state')||'{}'),a.set('Args',JSON.stringify(t)),fetch(`/SSA/${r.getAttribute('name')}/${Ae(n)}`,{method:'POST',body:a}).then(f=>f.headers.get('X-Redirect')?(window.location.href=f.headers.get('X-Redirect'),''):f.text()).then(f=>{if(!!f){if(r.hasAttribute('ssa-replace')){r.outerHTML=f;return}try{Y(r,f)}catch(s){console.log('Fallback from morphdom to root.outerHTML due to error',s),r.outerHTML=f}}}).catch(f=>{console.log(f)})}function be(e,n){let t=_({starter:e,depth:n.split('').filter(a=>a==='$').length,id:n.includes(':')?n.split(':')[0]:void 0});if(!t.getAttribute('state'))throw new Error('Bind call error: component state is underfined');let r=JSON.parse(decodeURIComponent(t.getAttribute('state')));r[n]=e.value,t.setAttribute('state',JSON.stringify(r))}function Te(e,n){n.preventDefault();let t=_({starter:e});if(!t.getAttribute('state'))throw new Error('Bind call error: component state is underfined');let r=JSON.parse(decodeURIComponent(t.getAttribute('state'))),a=new FormData(n.target),f=Object.fromEntries(a.entries());return Object.entries(f).forEach(s=>{r[s[0]]=s[1]}),t.setAttribute('state',JSON.stringify(r)),q(t,'Submit'),!1}window._LocaleRoot=_;window.Action=q;window.Bind=be;window.FormSubmit=Te;})();" +
	"</script>"
