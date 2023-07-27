import { action } from "./action";
import { root } from "./root";

// @ts-ignore
// Globals configuration
declare global {
    const actionpath: string
    const actionterminator: string
    interface Window {
        action: any
        root: any
    }
}

// Global scope
window.action = action
window.root = root
