import morphdom from 'morphdom';


// Internals

interface LocateParameters {
    starter: HTMLElement
    depth?: number  // parent call with $
    id?: string  // direct call with <id>:<action>
}

function _LocateRoot(parameters: LocateParameters): HTMLElement {
    let root: HTMLElement = parameters.starter
    // Direct call case
    if (parameters.id) {
        let element = document.getElementById(parameters.id)
        if (!element) {
            throw new Error(`Error while locating root: can't find direct with ${parameters}`)
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


// Public

export function Action(self: HTMLElement, action: string, ...args: Array<string>) {
    // Determine component root
    let root = _LocateRoot({
        starter: self,
        depth: action.split('').filter(x => x === '$').length,
        id: action.includes(':') ? action.split(':')[0] : undefined,
    })
    // Set loading state
    TriggerLoaders(self)
    // Prepare form data
    let formdata = new FormData()
    formdata.set('State', root.getAttribute('state') || '{}')
    formdata.set('Args', JSON.stringify(args))
    // Make request
    fetch(`/SSA/${root.getAttribute('name')}/${_NameCleanup(action)}`, {
        method: 'POST',
        body: formdata
    }).then(resp => {
        // Handle redirect header
        if (resp.headers.get('X-Redirect')) {
            window.location.href = resp.headers.get('X-Redirect') as string
            return ''
        }
        return resp.text()
    }).then(data => {
        // Handle no data case
        if (!data) {
            return
        }
        // Handle replace case
        if (root.hasAttribute('ssa:replace')) {
            console.log('Replace detected')
            root.outerHTML = data
            return
        }
        // Morph
        try {
            morphdom(root, data)
        } 
        catch (e: any) {
            console.log('Fallback from morphdom to root.outerHTML due to error', e)
            root.outerHTML = data
        }
    }).catch(err => {
        console.log(err)
    })
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
    root.setAttribute('state', btoa(JSON.stringify(state)))
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

function TriggerLoaders(self: HTMLElement) {
    // Find component root
    let root = _LocateRoot({
        starter: self,
    })
    // Find loader elements
    // Need escape for escaping (for payload wrapper)
    let loader = root.querySelectorAll('[ssa\\\\:oncall\\\\.display]')
    loader.forEach(element => {
        // Check attribute value
        let loadertype = element.getAttribute('ssa:oncall.display')
        // Show element, if value is empty
        if (loadertype == "") {
            element.setAttribute("style", "display: block");
        } else {
            // Set display value if exist
            element.setAttribute("style", "display: " + loadertype);
        }
    });
}


// Export to global

declare global {
    interface Window {
        _LocaleRoot: any;
        Action: any,
        Bind: any,
        FormSubmit: any
    }
}

window._LocaleRoot = _LocateRoot
window.Action = Action
window.Bind = Bind
window.FormSubmit = FormSubmit
