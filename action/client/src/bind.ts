import { root } from './root'

/**
 * bind allows to bind input field to the component field.
 * Works only with universal components because of the direct access to the component state.
*/
export async function bind(self: HTMLElement, field: string) {
  // Determine component root
  let _root = root({
      start: self,
  })
  // Check if we have a state
  if (!_root.getAttribute('state')) {
    throw new Error(`bind call error, component state is undefined`)
  }
  // Try to load state
  let state: any
  try {
    state = JSON.parse(decodeURIComponent(atob(_root.getAttribute('state') as string)))
  } catch (err) {
    throw new Error(`bind call error, can't decode state (probably it's not universal')`)
  }
  // Set value
  state[field] = (self as HTMLInputElement).value
  // Set state
  _root.setAttribute('state', btoa(encodeURIComponent(JSON.stringify(state))))
}
