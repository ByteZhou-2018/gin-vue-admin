import service from '@/utils/request'

export const zip = (data) => {
  return service({
    url: '/server/zip',
    method: 'POST',
    data
  })
}
export const deploy = (data) => {
  return service({
    url: '/server/deploy',
    method: 'POST',
    data
  })
}
