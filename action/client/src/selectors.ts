
/**
 * Display behavior attributes.
 * Applies specific behavior on action.
 *
 * display - set given display parameter on action.
 */
export enum Display {
  display = 'action\\:display'
}

/**
 * Render behavior attributes.
 *
 * render - use a given render strategy.
 * renderMorph - use morph render strategy.
 * renderReplace - use replace render strategy.
 */
export enum Render {
  render = 'action\\:render',
  renderMorph = 'action\\:render.morph',
  renderReplace = 'action\\:render.replace'
}

/**
 * Redirect behavior attributes.
 *
 * redirect - redirect to given location.
 */
export enum Redirect {
  redirect = 'action\\:redirect'
}

/**
 * Morph behavior attributes.
 *
 * ignore - ignore current node and all children.
 * ignoreThis - ignore only current node (children will be morphed).
 * ignoreAttr - ignore particular attribute morphing.
 */
export enum Morph {
  ignore = 'action\\:morph.ignore',
  ignoreThis = 'action\\:morph.ignore.this',
  ignoreAttr = 'action\\:morph.ignore.attr'
}
