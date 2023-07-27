

export interface Locate {
    start: HTMLElement
    depth?: number // parent call with $.<depth>.<action>
    id?: string // direct call with #.<id>.<action>
}

// Root allows to find the root element of the current component.
// Takes start node and searches for a root within parent nodes.
export function root(l: Locate): HTMLElement {
    // First, let's define root.
    let root: HTMLElement | null = null
    // We have different behavior depending on the type of the call.
    // First, let's handle direct call.
    if (l.id) {
        // Get root by id.
        root = document.getElementById(l.id)
    }
    // Traverse up to the root by default.
    if (!root) {
        // Cursor with a start point.
        let cursor = l.start
        // For locating deeper parents (like parent->parent->...).
        let depth = 0
        // Traverse mainloop.
        while (true) {
            // Handle no parent case (means we haven't found anything and reached the top).
            if (!cursor.parentElement) {
                break
            }
            // Handle cursor is not a component case.
            if (!cursor.getAttribute('state')) {
                cursor = cursor.parentElement
                continue
            }
            // If we are here, cursor is a component.
            // If current depth is exactly what we need, we found the root.
            if (l.depth === depth) {
                root = cursor
                break
            }
            // Otherwise, we need to go deeper.
            cursor = cursor.parentElement
            depth++
        }
    }
    // Handle empty root.
    if (!root) {
        throw new Error(`Root element not found with parameters ${l}.`)
    }
    // Return root
    return root
}