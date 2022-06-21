import morphdom from 'morphdom';

// Internals

interface LocateParameters {
    starter: HTMLElement
    depth?: number  // parent call with $
    id?: string  // direct call with <id>:<action>
}

interface MorphDomOptions {
    getNodeKey?: (node: Node) => any;
    onBeforeNodeAdded?: (node: Node) => Node;
    onNodeAdded?: (node: Node) => Node;
    onBeforeElUpdated?: (fromEl: HTMLElement, toEl: HTMLElement) => boolean;
    onElUpdated?: (el: HTMLElement) => void;
    onBeforeNodeDiscarded?: (node: Node) => boolean;
    onNodeDiscarded?: (node: Node) => void;
    onBeforeElChildrenUpdated?: (fromEl: HTMLElement, toEl: HTMLElement) => boolean;
    childrenOnly?: boolean;
}

function _LocateRoot(parameters: LocateParameters): HTMLElement {
    let root: HTMLElement = parameters.starter
    // Direct call case
    if (parameters.id) {
        let element = document.getElementById(parameters.id)
        if (!element) {
            throw new Error(`Error while locating root with id: can't find direct with ${parameters}`)
        }
        root = element
    } else {
        // Depth counter
        let dcount = 0
        // Search loop
        while (true) {
            // Handle error case
            if (!root.parentElement) {
                throw new Error(`Error while locating root: can't find parent with ${parameters}`)
            }
            if (!root.getAttribute('state')) { // Not a component, get parent
                root = root.parentElement
            } else { // Found a component
                if (parameters.depth && dcount != parameters.depth) { // Parent call case
                    root = root.parentElement
                    dcount++
                } else { // Exit clause
                    break
                }
            }
        }
    }
    return root
}

function _NameCleanup(action: string): string {
    if (action.includes(':')) {
        action = action.split(':')[1]
    }
    if (action.includes('$')) {
        action = action.replaceAll('$', '')
    }
    return action
}

// Updates display parameter for component DOM elements according to provided config
function _TriggerLoaders(root: HTMLElement) {
    // Find loader elements
    // Need escape for escaping (for payload wrapper)
    let loader = root.querySelectorAll('[ssa\\\\:oncall\\\\.display]')
    loader.forEach(element => {
        // Check attribute value
        let loadertype = element.getAttribute('ssa:oncall.display')
        // Set display value if exist
        if (loadertype != "") {
            element.setAttribute("style", "display: " + loadertype);
        }
    });
}

// Finds elements with onload attribute and executes provided action (without arguments)
export function _OnLoad() {
    document.querySelectorAll('[ssa\\\\:onload]').forEach(element => {
        let action = element.getAttribute('ssa:onload')
        if (action && action != "") {
            Action(element as HTMLElement, action)
        }
    });
}

// Find elements with poll attribute and executes provided action with interval
export function _Poll() {
    document.querySelectorAll('[ssa\\\\:poll]').forEach(element => {
        let action = element.getAttribute('ssa:poll') || ''
        let interval = element.getAttribute('ssa:poll.interval')

        if (action && action != "" && interval && interval != "") {
            setInterval(() => {
                Action(element as HTMLElement, action)
            }, parseInt(interval))
        }
    });
}

// Find elements with onintersect attribute and executes provided action on intersection
export function _OnIntersect() {
    document.querySelectorAll('[ssa\\\\:onintersect]').forEach(element => {
        let action = element.getAttribute('ssa:onintersect') || ''
        let threshold = element.getAttribute('ssa:onintersect.threshold') || '1.0'
        if (action != '') {
            let observer = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.intersectionRatio >= parseFloat(threshold)) {
                        Action(element as HTMLElement, action, parseFloat(threshold))
                    }
                })
            }, { threshold: parseFloat(threshold) })
            observer.observe(element)
        }
    });
}

function _Morph(fromNode: Node, toNode: Node | string, options?: MorphDomOptions) {
    let afterMorphIgnore = new Array()
    let newOptions: MorphDomOptions = {}
    if (options) {
        newOptions = options
    }
    newOptions!.onBeforeElUpdated = function(fromEl, toEl) {
        if (fromEl.getAttribute('ssa:morph.ignore.attr') != null) {
            let attr = fromEl.getAttribute('ssa:morph.ignore.attr')
            if (attr) {
                if (attr == 'innerHTML') {
                    toEl.innerHTML = fromEl.innerHTML
                } else {
                    let attrValue = fromEl.getAttribute(attr)
                    if (attrValue) toEl.setAttribute(attr, attrValue)
                }
            }
        }
        if (fromEl.getAttribute('ssa:morph.ignore') != null) {
            return false
        }
        if (fromEl.getAttribute('ssa:morph.ignore.this') != null && fromEl != fromNode) {
            afterMorphIgnore.push({ fromEl, toEl })
            return false
        }
        return true
    }
    morphdom(fromNode, toNode, newOptions)
    if (afterMorphIgnore.length > 0) {
        afterMorphIgnore.forEach(el => {
            _Morph(el.fromEl, el.toEl, {
                childrenOnly: true
            })
        })
    }
}


// Public

export async function Action(self: HTMLElement, action: string, ...args: Array<any>): Promise<void> {
    // Determine component root
    let root = _LocateRoot({
        starter: self,
        depth: action.split('').filter(x => x === '$').length,
        id: action.includes(':') ? action.split(':')[0] : undefined,
    })
    // Set loading state
    _TriggerLoaders(root)
    // Build URL
    let url = ssapath // Base
    if (!url.endsWith('/')) {
        url += '/'
    }
    url += `${root.getAttribute('name')}`  // Component name
    url += `/${_NameCleanup(action)}` // Action name
    // Make request
    const payload = new FormData()
    payload.set('State', root.getAttribute('state') as string)
    payload.set('Args', JSON.stringify(args))
    const response = await fetch(url, {
        method: 'POST',
        body: payload,
    })
    // Extract stream reader
    const reader = response.body?.pipeThrough(new TextDecoderStream()).getReader()
    if (!reader) {
        return
    }
    // Read responses
    while (true) {
        // Extract value and status
        const { value, done } = await reader.read()
        if (done) {
            break
        }
        // Handle redirect
        if (value.startsWith('ssa:redirect=')) {
            window.location.href = value.replace('ssa:redirect=', '')
            continue
        }
        // Handle replace
        if (root.getAttribute('ssa:render.mode') == 'replace') {
            root.outerHTML = value
            continue
        }
        // Morph by default
        try {
            _Morph(root, value)
        } catch (e: any) {
            console.log('Fallback from morphdom to root.outerHTML due to error', e)
            root.outerHTML = value
        }
    }
}

export function Bind(self: HTMLElement, field: string) {
    // Find component root
    let root = _LocateRoot({
        starter: self,
        depth: field.split('').filter(x => x === '$').length,
        id: field.includes(':') ? field.split(':')[0] : undefined,
    })
    // Check state
    if (!root.getAttribute('state')) {
        throw new Error('Bind call error: component state is underfined')
    }
    // Load state
    let state = JSON.parse(atob(root.getAttribute('state') as string))
    // Set value
    state[field] = (self as HTMLInputElement).value
    // Set state
    root.setAttribute('state', btoa(state))
}

export function FormSubmit(self: HTMLElement, e: Event) {
    // Prevent default submit
    e.preventDefault()
    // Find component root
    let root = _LocateRoot({
        starter: self
    })
    // Check state
    if (!root.getAttribute('state')) {
        throw new Error('Bind call error: component state is underfined')
    }
    // Load state
    let state = JSON.parse(atob(root.getAttribute('state') as string))
    // Update state with form data
    let form = new FormData((e.target as HTMLFormElement))
    let formdata = Object.fromEntries(form.entries())
    Object.entries(formdata).forEach(pair => {
        state[pair[0]] = pair[1]
    })
    // Set state
    root.setAttribute('state', btoa(JSON.stringify(state)))
    // Trigger "Submit" action
    Action(root, 'Submit')
    // Fix for ...?
    // Can't remember the issue
    return false
}


// Export to global

declare global {
    const ssapath: string
    interface Window {
        
        _LocaleRoot: any
        Action: any
        Bind: any
        FormSubmit: any
    }
}

window._LocaleRoot = _LocateRoot
window.Action = Action
window.Bind = Bind
window.FormSubmit = FormSubmit

document.addEventListener('DOMContentLoaded', _OnLoad)
document.addEventListener('DOMContentLoaded', _OnIntersect)
document.addEventListener('DOMContentLoaded', _Poll)