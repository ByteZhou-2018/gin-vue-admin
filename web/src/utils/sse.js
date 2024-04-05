import { fetchEventSource } from '@microsoft/fetch-event-source'

const basePath = import.meta.env.VITE_BASE_API
export const useSSE = (url, {
  onmessage,
  onopen,
  onclose,
  onerror,
  data,
  headers = {
    'Content-Type': 'application/json',
  },
  openWhenHidden = true,
}) => {
  const ctrl = new AbortController()

  return [ctrl, () => {
    fetchEventSource(` ${basePath}${url}`, {
      method: 'POST',
      openWhenHidden,
      headers,
      body: JSON.stringify(data),
      signal: ctrl.signal,
      onmessage: onmessage,
      onopen: onopen,
      onclose: onclose,
      onerror: onerror
    })
  }]
}
