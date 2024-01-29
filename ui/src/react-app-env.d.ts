declare module '*.less' {
  const resource: { [key: string]: string }
  export = resource
}
declare module '*.module.less' {
  const classes: { readonly [key: string]: string }
  export default classes
}
declare module '*.png'
declare module '*.svg'
declare module '*.jpeg'
declare module '*.jpg'
