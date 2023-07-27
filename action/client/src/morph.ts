import morphdom from 'morphdom'

// MorphOptions provides a way to extend morph with own methods and configurations.
export interface MorphOptions {
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

// morph provides a way to morph one DOM node into another,
// without touching similar existing nodes.
// Provides a way to ignore elements with certain attributes and elements.
//
// action:morph.ignore - ignore current node and all children.
// action:morph.ignore.this - ignore only current node (children will be morphed).
// action:morph.ignore.attr - ignore particular attribute morphing.
export function morph(fromNode: Node, toNode: Node | string, options?: MorphOptions) {
    let afterMorphIgnore = new Array()
    let newOptions: MorphOptions = {}
    if (options) {
        newOptions = options
    }
    newOptions!.onBeforeElUpdated = function(fromEl, toEl) {
        if (fromEl.getAttribute('action:morph.ignore.attr') != null) {
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
        if (fromEl.getAttribute('action:morph.ignore') != null) {
            return false
        }
        if (fromEl.getAttribute('action:morph.ignore.this') != null && fromEl != fromNode) {
            afterMorphIgnore.push({ fromEl, toEl })
            return false
        }
        return true
    }
    morphdom(fromNode, toNode, newOptions)
    if (afterMorphIgnore.length > 0) {
        afterMorphIgnore.forEach(el => {
            morph(el.fromEl, el.toEl, {
                childrenOnly: true
            })
        })
    }
}