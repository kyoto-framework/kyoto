import { action } from './action'
import { bind } from './bind'
import { root } from './root'

// @ts-ignore
// Globals configuration
declare global {
    const actionpath: string
    const actionterminator: string
    interface Window {
        action: any
        bind: any
        root: any
    }
}

// Global scope
window.action = action
window.bind = bind
window.root = root
