export function logError(...args: any[]) {
  // central place to manage error logging; adapt for remote logging later
  if (process.env.NODE_ENV === 'production') {
    // TODO: send to remote logging (Sentry, etc.)
    // no-op for now to avoid leaking local details
    return;
  }
  // keep a console fallback for development
   
  console.error(...args);
}

export function logInfo(...args: any[]) {
  if (process.env.NODE_ENV !== 'production') {
     
    console.log(...args);
  }
}
