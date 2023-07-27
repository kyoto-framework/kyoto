import { root, Locate } from "./root";
import { morph, MorphOptions } from "./morph";


interface Action {
    name: string
    depth: number
    id?: string
}


// action allows to call server side function.
// After function execution,
// component will be rendered and morphed into current DOM.
export async function action(self: HTMLElement, action: string, ...args: any[]) {
    // Parse action
    let _action = actionParse(action)
    // Determine component root
    let _root = root({
        start: self,
        depth: _action.depth,
        id: _action.id
    })
    // Trigger action display
    actionDisplay(_root)
    // Build action url
    let url = actionpath.endsWith('/') ? actionpath : actionpath + '/'
    url += `${_root.getAttribute('name')}`
    url += `/${_action.name}`
    // Build action payload
    let payload = new FormData()
    payload.set('State', _root.getAttribute('state') as string)
    payload.set('Args', JSON.stringify(args))
    // Use XHR to load chunks.
    // Each chunk is a new component layout update.
    // Using XHR due to lack of support TextDecoderStream by Firefox.
    //
    // Also, we are using buffer & terminator sequence to ensure integrity.
    // Somehow chunk becomes split sometimes which leads to broken render.
    // This is a workaround solution, actual reason of such behavior wasn't found.
    let xhr = new XMLHttpRequest()
    xhr.open('POST', url, true)
    let buffer = ''
    let cursor = 0
    xhr.onprogress = () => {
        // Determine current cursor
        let cursorNow = xhr.responseText.length
        // If we are in the end, exit
        if (cursor == cursorNow) return
        // Get chunk
        const chunk = xhr.responseText.substring(cursor, cursorNow)
        // Add to buffer
        buffer += chunk
        console.log(buffer, buffer.endsWith(actionterminator))
        // If buffer ends with terminator sequence, remove it and render
        if (buffer.endsWith(actionterminator)) {
            // Remove terminator
            buffer = buffer.slice(0, -(actionterminator.length))
            // Handle redirect
            if (buffer.startsWith('action:redirect=')) {
                window.location.href = buffer.replace('action:redirect=', '')
                return
            }
            // Determine render mode (morph by default)
            const mode = _root.getAttribute('action:render.mode') || 'morph'
            // Render
            switch (mode) {
                case 'replace':
                    _root.outerHTML = buffer
                    break;
                case 'morph':
                    try {
                        morph(_root, buffer)
                    } catch (e: any) {
                        console.log('Fallback from "morphdom" to "replace" due to an error:', e)
                        _root.outerHTML = buffer
                    }
                    break
                default:
                    console.log('Render mode is not supported, fallback to "replace"')
                    _root.outerHTML = buffer
                    break;
            }
            // Cleanup buffer
            buffer = ''
        }
        // Increment cursor
        cursor = cursorNow
    }
    // Send request
    xhr.send(payload)
}

// actionParse is a part of the action,
// that parses action call according to specific semantic.
// It allows to extract action name, call depth, direct id, etc.
function actionParse(action: string): Action {
    // Split action into tokens
    let tokens = action.split('.')
    // Define action properties
    let name = tokens[tokens.length - 1]
    let depth = 0
    let id = undefined
    if (action.startsWith('#')) {
        id = tokens[1]
    }
    if (action.startsWith('$')) {
        depth = parseInt(tokens[1])
    }
    // Return action
    return {
        name: name,
        depth: depth,
        id: id
    }
}

// actionDisplay is a part of the action,
// which responsible for handling display state during action.
// Default display will be passed as a ready for use markup, so we don't need to provide a reset function.
function actionDisplay(root: HTMLElement) {
    // Find all elements with action display attribute
    let attr = `action\\\\:display`
    let elements = root.querySelectorAll(`[${attr}]`)
    elements.forEach(element => {
        // Get display attribute value
        let display = element.getAttribute(attr)
        // Set display value, if not empty
        if (display) {
            (element as HTMLElement).style.display = display
        }
    })
}